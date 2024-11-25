package counter

import (
	"fmt"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"log"
	"net"
	"sync"
	"time"
)

// remove -type event if you won't use diy struct in kernel
//
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go bpf counter.c -- -I /sys/kernel/btf
type CounterReq struct {
	IfName string
}
type CounterRes struct {
	Count uint64
}

func Start(req CounterReq) (<-chan CounterRes, func()) {
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
		Program:   objs.CountPackets,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatal("Attaching XDP:", err)
	}
	out := Action(objs, req, stopper)
	buildClose := func() func() {
		once := sync.Once{}
		return func() {
			once.Do(func() {
				objs.Close()
				attachXDP.Close()
				close(out)
				close(stopper)
			})
		}
	}
	log.Printf("Counting incoming packets on %s..", ifname)
	return out, buildClose()
}

func Action(objs bpfObjects, req CounterReq, stopper chan struct{}) chan CounterRes {
	// add your link logic here
	out := make(chan CounterRes)
	tick := time.Tick(time.Second)
	go func() {
		for {
			select {
			case <-tick:
				var count uint64
				err := objs.PktCount.Lookup(uint32(0), &count)
				if err != nil {
					fmt.Println("Map lookup:", err)
				}
				out <- CounterRes{Count: count}
			case <-stopper:
				log.Print("Received signal, exiting..")
				return
			}
		}
	}()
	return out
}

// add more action function here
