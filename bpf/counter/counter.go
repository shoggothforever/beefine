package counter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/rlimit"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type pkt_info bpf counter.c -- -I /sys/kernel/btf
type CounterReq struct {
	IfName     string
	perfReader *perf.Reader
}
type CounterRes struct {
	Count    uint64
	SrcIp    string
	DstIp    string
	Protocol string
	SrcPort  uint16
	DstPort  uint16
}

func Start(req *CounterReq) (<-chan CounterRes, func()) {
	stopper := make(chan struct{})
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Print(err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Printf("loading objects: %v", err)
	}

	ifname := req.IfName // Change this to an interface on your machine.
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		log.Fatalf("Getting interface %s: %s", ifname, err)
	}
	// Attach count_packets to the network interface.
	attachXDP, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.CountAndInfo,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatal("Attaching XDP:", err)
	}
	req.perfReader, err = perf.NewReader(objs.Events, os.Getpagesize())
	if err != nil {
		log.Fatal("open perf reader:", err)
	}
	out := Action(objs, req, stopper)
	buildClose := func() func() {
		once := sync.Once{}
		return func() {
			once.Do(func() {
				objs.Close()
				req.perfReader.Close()
				attachXDP.Close()
				close(stopper)
				time.Sleep(1 * time.Second)
				close(out)

			})
		}
	}
	log.Printf("Counting incoming packets on %s..", ifname)
	return out, buildClose()
}

func Action(objs bpfObjects, req *CounterReq, stopper chan struct{}) chan CounterRes {
	// add your link logic here
	out := make(chan CounterRes)
	tick := time.Tick(time.Second)
	go func() {
		var count uint64
		var res CounterRes
		for {
			select {
			case <-tick:
				record, err := req.perfReader.Read()
				if err != nil {
					fmt.Printf("Error reading perf event: %v\n", err)
					continue
				}
				err = objs.PktCount.Lookup(uint32(0), &count)
				if err != nil {
					fmt.Println("Map lookup:", err)
				}
				var info bpfPktInfo
				if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &info); err != nil {
					fmt.Printf("Failed to parse packet info: %v\n", err)
					continue
				}
				fmt.Printf("Src IP: %s, Dst IP: %s, Src Port: %d, Dst Port: %d, Protocol: %s\n",
					ipToString(info.SrcIp), ipToString(info.DstIp), info.SrcPort, info.DstPort, protoToString(info.Protocol))
				res.build(&info, count)
				out <- res
			case <-stopper:
				log.Print("Received signal, exiting..")
				return
			}
		}
	}()
	return out
}

// add more action function here
var protoMap = map[uint8]string{
	0:   "IP",
	1:   "ICMP",
	2:   "IGMP",
	4:   "IPIP",
	6:   "TCP",
	8:   "EGP",
	12:  "PUP",
	17:  "UDP",
	22:  "IDP",
	29:  "TP",
	33:  "DCCP",
	41:  "IPV6",
	46:  "RSVP",
	47:  "GRE",
	50:  "ESP",
	51:  "AH",
	92:  "MTP",
	94:  "BEETPH",
	98:  "ENCAP",
	103: "PIM",
	108: "COMP",
	115: "L2TP",
	132: "SCTP",
	136: "UDPLITE",
	137: "MPLS",
	143: "ETHERNET",
	255: "RAW",
}

func protoToString(proto uint8) string {
	return protoMap[proto]
}

// ipToString 将 IP 地址转换为字符串
func ipToString(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip), byte(ip>>8), byte(ip>>16), byte(ip>>24))
}

func (i *CounterRes) build(info *bpfPktInfo, count uint64) {
	i.Count = count
	i.SrcIp = ipToString(info.SrcIp)
	i.DstIp = ipToString(info.DstIp)
	i.Protocol = protoToString(info.Protocol)
	i.SrcPort = info.SrcPort
	i.DstPort = info.DstPort
}
