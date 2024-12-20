package mount

import (
	"bytes"
	"encoding/binary"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	"log"
	"sync"
)

// remove -type event if you won't use diy struct in kernel
//
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpfel -type mount_event bpf mount.c -- -I /sys/kernel/btf
type MountReq struct {
	rd *ringbuf.Reader
}
type MountRes struct {
	Pid     uint32
	DevName [64]byte
	DirName [64]byte
	Type    [32]byte
}

func Start(req *MountReq) (chan MountRes, func()) {
	stopper := make(chan struct{})
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}
	// Load pre-compiled programs and maps into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	// write your link code here
	enterLink, err := link.Tracepoint("syscalls", "sys_enter_mount", objs.TraceEnterMount, nil)
	if err != nil {
		log.Fatalf("Failed to attach sys_enter_mount tracepoint: %v", err)
	}
	// 打开 perf buffer
	rd, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		log.Fatalf("Failed to open perf buffer: %v", err)
	}
	req.rd = rd
	out := Action(objs, req, stopper)
	// using closure to create only once close function
	buildClose := func() func() {
		once := sync.Once{}
		return func() {
			once.Do(func() {
				objs.Close()
				enterLink.Close()
				rd.Close()
				// close attach
				close(out)
				close(stopper)
			})
		}
	}

	return out, buildClose()

}
func Action(objs bpfObjects, req *MountReq, stopper chan struct{}) chan MountRes {
	// add your link logic here
	out := make(chan MountRes)
	go func() {
		var event MountRes
		for {
			// write your logical code here
			select {
			case <-stopper:
				return
			default:
				record, err := req.rd.Read()
				if err != nil {
					log.Printf("reading ringbuf: %s", err)
					return
				}
				if err = binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &event); err != nil {
					return
				}
				out <- event
			}
		}
	}()
	return out
}

// add more action function here
