//go:build ignore

#include "../vmlinux.h"
#include "../headers/common.h"
#include "../headers/bpf_endian.h"
#include "../headers/bpf_tracing.h"

char LICENSE[] SEC("license") = "Dual BSD/GPL";
struct event {
	u8 comm[16];
	__u16 val;
};
struct event *unused __attribute__((unused));

SEC("XXX")
int handle_XXX(){
    // write your code here
	return 0;
}
