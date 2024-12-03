  //go:build ignore

#include "../vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>

char LICENSE[] SEC("license") = "Dual BSD/GPL";
struct event {
    __u32 pid;
    __u64 bytes; // 记录文件读写字节数
    char comm[16];
    char operation[16];
    char filename[256];
};
struct event *unused __attribute__((unused));

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 256 * 1024);
} es SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries,1024);
    __type(key,char[16]);
    __type(value,__u64);
} imageMap SEC(".maps");

SEC("tracepoint/syscalls/sys_enter_openat")
int trace_openat(struct trace_event_raw_sys_enter *ctx) {
    struct event *e;
    const char *filename = (const char *)ctx->args[1];
    __u32 pid = bpf_get_current_pid_tgid() >> 32;
    e = bpf_ringbuf_reserve(&es, sizeof(*e), 0);
    if (!e) return 0;

    bpf_get_current_comm(&e->comm, sizeof(e->comm));
    bpf_probe_read_str(&e->filename, sizeof(e->filename), filename);
    bpf_probe_read_str(&e->operation, sizeof(e->operation), "openat");
    e->pid = pid;
    e->bytes = 0; // 初始为 0
    bpf_ringbuf_submit(e, 0);

    return 0;
}

SEC("tracepoint/syscalls/sys_enter_read")
int trace_read(struct trace_event_raw_sys_enter *ctx) {
    __u64 fd = ctx->args[0];
    __u64 bytes = ctx->args[2];
    __u32 pid = bpf_get_current_pid_tgid() >> 32;
    struct event *e;
    e = bpf_ringbuf_reserve(&es, sizeof(*e), 0);
    if (!e) return 0;

    bpf_get_current_comm(&e->comm, sizeof(e->comm));
    bpf_probe_read_str(&e->operation, sizeof(e->operation), "read");
    e->pid = pid;
    e->bytes = bytes;

    bpf_ringbuf_submit(e, 0);
    return 0;
}