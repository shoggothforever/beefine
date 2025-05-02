  //go:build ignore

//#include "../vmlinux.h"
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#define IPPROTO_UDP 17
#define IPPROTO_TCP 6
#define ETH_HDR_LEN 14 // 以太网头部长度
char LICENSE[] SEC("license") = "Dual BSD/GPL";

SEC("xdp")
int pod_network_monitor(struct xdp_md *ctx) {
    void* data_end=(void*)(long)ctx->data_end;
    void* data=(void*)(long)ctx->data;
    // 捕获以太网帧
    struct ethhdr *eth =data;
    // 只关心 IP 包
    if ((void*)(eth+1)>data_end){
            return XDP_PASS;
    }
    // 获取 IP 层
    struct iphdr *ip = data+ETH_HDR_LEN;
    if ((void*)(ip+1)>data_end){
        return XDP_PASS;
    }
    // 进一步检查 TCP 或 UDP 层
    if (ip->protocol == IPPROTO_TCP) {
        struct tcphdr *tcp = (void*)ip+ip->ihl*4;
        if((void*)(tcp+1)>data_end){
            return XDP_PASS;
        }

    } else if (ip->protocol == IPPROTO_UDP) {
        struct udphdr *udp = (void*)ip+ip->ihl*4;
        if((void*)(udp+1)>data_end){
            return XDP_PASS;
        }
    }
    return XDP_PASS;  // 继续传递数据包
}