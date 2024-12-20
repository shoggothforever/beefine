  //go:build ignore

#include "../vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>

char LICENSE[] SEC("license") = "Dual BSD/GPL";

// 定义输出到 perf buffer 的数据结构
struct mount_event {
    __u32 pid;
    char dev_name[64];
    char dir_name[64];
    char type[32];
};
struct mount_event *unused __attribute__((unused));

// 定义 perf buffer
struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 64 * 1024);
} events SEC(".maps");

// 捕获 sys_enter_mount 事件
SEC("tracepoint/syscalls/sys_enter_mount")
int trace_enter_mount(struct trace_event_raw_sys_enter *ctx) {
    struct mount_event *evt;
    __u32 pid = bpf_get_current_pid_tgid() >> 32;
    evt = bpf_ringbuf_reserve(&events, sizeof(*evt), 0);
    if (!evt){
        return 0;
    }
    evt->pid = pid;
    // 读取挂载参数
    bpf_probe_read_str(&evt->dev_name, sizeof(evt->dev_name), (const char*)ctx->args[0]);
    bpf_probe_read_str(&evt->dir_name, sizeof(evt->dir_name),  (const char*)ctx->args[1]);
    bpf_probe_read_str(&evt->type, sizeof(evt->type),  (const char*)ctx->args[2]);
    // 输出事件到用户空间
    bpf_ringbuf_submit(evt, 0);
    return 0;
}