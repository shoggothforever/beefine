  //go:build ignore


#include "../vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <stdbool.h>
//go:build ignore

#define TASK_COMM_LEN	 16
#define MAX_FILENAME_LEN 127
const volatile unsigned long long min_duration_ns = 0;

char LICENSE[] SEC("license") = "Dual BSD/GPL";
const int32 container_map_key=0;
struct {
    __uint(type,BPF_MAP_TYPE_ARRAY);
    __uint(max_entries,1024);
    __type(key,int32);
    __type(value,int32);
} cg_pid_map SEC(".maps");
struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 256 * 1024);
} rb SEC(".maps");

struct event {
	int32 pid;
	int32 prio;
	__u64 ts;
	char comm[TASK_COMM_LEN];
	bool exit_event;
};
struct event *unused __attribute__((unused));

SEC("tracepoint/sched/sched_process_exec")
int handle_exec(struct trace_event_raw_sched_process_exec *ctx)
{
	unsigned fname_off;
	struct event *e;
	int pid;
	int *cg_pid;
	__u64 ts;

	/* remember time exec() was executed for this PID */
	pid = bpf_get_current_pid_tgid() >> 32;
    /* don't emit exec events when minimum duration is specified */
    cg_pid=bpf_map_lookup_elem(&cg_pid_map,&container_map_key);
    if (!cg_pid)
        return 0;
    if(pid != *cg_pid)
        return 0;
	/* reserve sample from BPF ringbuf */
	e = bpf_ringbuf_reserve(&rb, sizeof(*e), 0);
	if (!e)
		return 0;

	e->exit_event = false;
	e->pid = pid;
	ts = bpf_ktime_get_ns();
	e->ts=ts;
	bpf_get_current_comm(&e->comm, sizeof(e->comm));
	bpf_ringbuf_submit(e, 0);
	/* successfully submit it to user-space for post-processing */

	return 0;
}


SEC("tracepoint/sched/sched_process_exit")
int handle_exit(struct trace_event_raw_sched_process_template *ctx)
{
	struct event *e;
	int pid, tid,prio;
    int *cg_pid;
	__u64 id, ts, *start_ts, duration_ns = 0;
	/* get PID and TID of exiting thread/process */
	id = bpf_get_current_pid_tgid();
	pid = id >> 32;
	tid = (__u32)id;
	/* ignore thread exits */
	if (pid != tid)
		return 0;
    cg_pid=bpf_map_lookup_elem(&cg_pid_map,&container_map_key);
    if (!cg_pid)
        return 0;
    struct task_struct *task = (struct task_struct *)bpf_get_current_task();
    __u32 host_ppid = BPF_CORE_READ(task,real_parent,tgid);
    if (pid != *cg_pid)
        return 0;
	/* reserve sample from BPF ringbuf */
	e = bpf_ringbuf_reserve(&rb, sizeof(struct event), 0);
	if (!e)
		return 0;

	/* fill out the sample with data */
	// task = (struct task_struct *)bpf_get_current_task();
	/* if we recorded start of the process, calculate lifetime duration */
//	start_ts = bpf_map_lookup_elem(&ct, &pid);

	bpf_map_delete_elem(&cg_pid_map, &pid);
    e->prio=ctx->prio;
	e->exit_event = true;
	e->ts = bpf_ktime_get_ns();
	e->pid = pid;
	bpf_get_current_comm(&e->comm, sizeof(e->comm));

	/* send data to user-space for post-processing */
	bpf_ringbuf_submit(e, 0);



	return 0;
}
