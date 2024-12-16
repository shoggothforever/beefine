package image_prep

import (
	"bytes"
	"encoding/binary"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	"log"
	"shoggothforever/beefine/internal/helper"
	"sync"
	"time"
)

// remove -type event if you won't use diy struct in kernel
//
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -target bpf -fbuiltin " -target bpfel -type event bpf image_prep.c -- -I /sys/kernel/btf
type ImagePrepReq struct {
	rb *ringbuf.Reader
}
type ImagePreRes struct {
	Pid       uint32
	_         [4]byte
	Bytes     uint64
	Comm      [16]byte
	Operation [16]byte
	Filename  [256]byte
}

func Start(req *ImagePrepReq) (chan ImagePreRes, func()) {
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
	oa, err := link.Tracepoint("syscalls", "sys_enter_openat", objs.TraceOpenat, nil)
	if err != nil {
		log.Fatalf("link tracepoint/syscalls/sys_enter_openat: %v", err)
	}
	//oe, err := link.Tracepoint("syscalls", "sys_enter_read", objs.TraceRead, nil)
	//if err != nil {
	//	log.Fatalf("link tracepoint/syscalls/sys_enter_read: %v", err)
	//}
	req.rb, err = ringbuf.NewReader(objs.bpfMaps.Es)
	if err != nil {
		log.Fatalf("ringbuf.NewReader: %v", err)
	}
	out := Action(objs, req, stopper)
	buildClose := func() func() {
		once := sync.Once{}
		return func() {
			once.Do(func() {
				// close attach
				objs.Close()
				oa.Close()
				//oe.Close()
				req.rb.Close()
				close(stopper)
				time.Sleep(time.Second)
				close(out)
			})
		}
	}
	return out, buildClose()

}
func Action(objs bpfObjects, req *ImagePrepReq, stopper chan struct{}) chan ImagePreRes {
	// add your link logic here
	out := make(chan ImagePreRes)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				return
			}
		}()

		var event ImagePreRes
		for {
			// write your logical code here
			select {
			case <-stopper:
				return
			default:
				record, err := req.rb.Read()
				if err != nil {
					log.Printf("reading ringbuf: %s", err)
					return
				}
				if err = binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &event); err != nil {
					log.Printf("reading record: %s", err)
					return
				}
				if helper.Bytes2String(event.Comm[:]) == "systemd-oomd" {
					continue
				}
				out <- event

			}

		}
	}()
	return out
}

// add more action function here
