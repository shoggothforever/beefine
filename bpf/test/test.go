package test

import (
	"github.com/cilium/ebpf/rlimit"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// remove -type event if you won't use diy struct in kernel
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpfel -type event bpf test.c -- -I../headers
type TestReq struct {
}
type TestRes struct {
}

func Start(req TestReq) (<-chan TestRes, func()) {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}
	// Load pre-compiled programs and maps into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()
	// write your link code here
	return Action(objs, req, stopper), func() { stopper <- os.Interrupt }

}
func Action(objs bpfObjects, req TestReq, stopper chan os.Signal) <-chan TestRes {
	// add your link logic here
	out := make(chan TestRes)
	go func() {
		for {
			// write your logical code here
			select {
			case <-stopper:
				return
			default:
				time.Sleep(1 * time.Second)
			}

		}
	}()
	return out
}

// add more action function here
