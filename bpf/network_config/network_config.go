package network_config

import (
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"
  "github.com/cilium/ebpf/link"
  "github.com/cilium/ebpf/rlimit"
)
// remove -type event if you won't use diy struct in kernel
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpfel -type event bpf network_config.c -- -I /sys/kernel/btf
type Network_configReq struct {

}
type Network_configRes struct {

}
func Start(req *Network_configReq) ( chan Network_configRes,func()) {
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

	  out:=Action(objs, req, stopper)
  	buildClose := func() func() {
  		once := sync.Once{}
  		return func() {
  			once.Do(func() {
  				objs.Close()
  				// close attach
  				close(out)
  				close(stopper)
  			})
  		}
  	}

  return out,buildClose()

}
func Action(objs bpfObjects , req Network_configReq , stopper chan struct{}) chan Network_configRes{
  // add your link logic here
  out := make(chan Network_configRes)
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
