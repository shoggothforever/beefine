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

const CgMapKey int32 = 0

type ExecReq struct {
	ContainerPid uint32
	rb           *ringbuf.Reader
}
type ExecRes struct {
	Pid       uint32
	Prio      uint32
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
	// using closure to create only once close function
	buildClose := func() func() {
		once := sync.Once{}
		return func() {
			once.Do(func() {
				objs.Close()
				execStart.Close()
				execExit.Close()
				req.rb.Close()
				close(stopper)
				time.Sleep(1 * time.Second)
				close(out)
			})
		}
	}

	return out, buildClose()

}
func Action(objs bpfObjects, req *ExecReq, stopper chan struct{}) chan ExecRes {
	out := make(chan ExecRes)
	go func() {
		var e ExecRes
		var event bpfEvent
		err := objs.CgPidMap.Put(CgMapKey, req.ContainerPid)
		if err != nil {
			log.Printf("unable to set cgmap, cg_pid:%d\n ,%s", req.ContainerPid, err.Error())
			return
		}
		log.Printf("container's pid: %d\n", req.ContainerPid)
		for {
			select {
			case <-stopper:
				return
			default:
				record, err := req.rb.Read()
				if err != nil {
					log.Printf("reading ringbuf: %s\n", err)
					return
				}
				if err = binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &event); err != nil {
					log.Printf("reading record: %s\n", err)
				}
				if err != nil {
					log.Printf("unable to lookup cgmap, cg_pid:%d,%s\n", req.ContainerPid, err.Error())
				}
				e.ExitEvent = event.ExitEvent
				e.Pid = event.Pid
				e.Ts = event.Ts
				e.Prio = event.Prio
				for k, v := range event.Comm {
					e.Comm[k] = byte(v)
				}
				out <- e
			}
		}
	}()
	return out
}

// add more action function here
