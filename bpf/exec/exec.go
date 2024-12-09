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
	ContainerPid int32
	rb           *ringbuf.Reader
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
				if err = binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &e); err != nil {
					log.Printf("reading record: %s\n", err)
				}
				var pid int32
				err = objs.CgPidMap.Lookup(CgMapKey, &pid)
				if err != nil {
					log.Printf("unable to lookup cgmap, cg_pid:%d,%s\n", req.ContainerPid, err.Error())
				}
				log.Printf("catch pid is %d\n", e.Pid)
				out <- e
			}
		}
	}()
	return out
}

// add more action function here
