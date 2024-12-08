package exec

import (
	"bytes"
	"encoding/binary"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	"log"
	"sync"
	"time"
)

// remove -type event if you won't use diy struct in kernel

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -type event -target bpfel  bpf exec.c -- -I  /sys/kernel/btf
type ExecReq struct {
	rb *ringbuf.Reader
}
type ExecRes struct {
	Pid       int32
	Prio      int32
	Ts        uint64
	Comm      [16]byte
	ExitEvent bool
	_         [7]byte
}

func Start(req *ExecReq) (chan ExecRes, func()) {
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
	execStart, err := link.Tracepoint("sched", "sched_process_exec", objs.HandleExec, nil)
	if err != nil {
		log.Fatalf("opening tracepoint: %s", err)
	}
	execExit, err := link.Tracepoint("sched", "sched_process_exit", objs.HandleExit, nil)
	if err != nil {
		log.Fatalf("opening tracepoint: %s", err)
	}
	req.rb, err = ringbuf.NewReader(objs.bpfMaps.Rb)
	if err != nil {
		log.Fatalf("opening ringbuf reader: %s", err)
	}
	out := Action(objs, req, stopper)
	buildClose := func() func() {
		once := sync.Once{}
		return func() {
			once.Do(func() {
				objs.Close()
				execStart.Close()
				execExit.Close()
				req.rb.Close()
				// close attach
				close(out)
				close(stopper)
			})
		}
	}

	return out, buildClose()

}
func Action(objs bpfObjects, req *ExecReq, stopper chan struct{}) chan ExecRes {
	out := make(chan ExecRes)
	go func() {
		var e ExecRes
		for {
			select {
			case <-stopper:
				return
			default:
				record, err := req.rb.Read()
				if err != nil {
					log.Printf("reading ringbuf: %s", err)
					return
				}
				if err = binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &e); err != nil {
					log.Printf("reading record: %s", err)
				}
				out <- e
				time.Sleep(1 * time.Second)
			}
		}
	}()
	return out
}

// add more action function here
