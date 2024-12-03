  //go:build ignore

#include "../vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>

  char LICENSE[] SEC("license") = "Dual BSD/GPL";
  struct event {
  	__u8 comm[16];
  	__u16 val;
  };
  struct event *unused __attribute__((unused));

  SEC("XXX")
  int handle_XXX(){
      // write your code here
  	return 0;
  }
