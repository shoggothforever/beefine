
//go:build ignore
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

// 数据结构：用于存储网络包的信息
struct pkt_info {
    __u32 src_ip;    // 源 IP 地址
    __u32 dst_ip;    // 目标 IP 地址
    __u16 src_port;  // 源端口号
    __u16 dst_port;  // 目标端口号
    __u8  protocol;  // 协议类型 (TCP/UDP)
};
struct pkt_info *unused __attribute__((unused));
// 包计数 Map
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 1);
} pkt_count SEC(".maps");

// 用于传递数据包信息到用户空间的 perf event map
struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
} events SEC(".maps");

// XDP 程序：统计包数量并提取网络包信息
SEC("xdp")
int count_and_info(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;

    // 1. 包计数
    __u32 key = 0;
    __u64 *count = bpf_map_lookup_elem(&pkt_count, &key);
    if (count) {
        __sync_fetch_and_add(count, 1);
    }

    // 2. 解析以太网头部
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end) {
        return XDP_PASS; // 检查包大小合法性
    }

    // 只处理 IPv4 数据包
    if (eth->h_proto != __constant_htons(ETH_P_IP)) {
        return XDP_PASS;
    }

    // 3. 解析 IP 头部
    struct iphdr *ip = data + ETH_HDR_LEN;
    if ((void *)(ip + 1) > data_end) {
        return XDP_PASS;
    }

    struct pkt_info pkt = {};
    pkt.src_ip = ip->saddr;      // 源 IP
    pkt.dst_ip = ip->daddr;      // 目标 IP
    pkt.protocol = ip->protocol; // 协议类型

    // 4. 解析 TCP/UDP 头部
    if (ip->protocol == IPPROTO_TCP) {
        struct tcphdr *tcp = (void *)ip + ip->ihl * 4;
        if ((void *)(tcp + 1) > data_end) {
            return XDP_PASS;
        }
        pkt.src_port = __constant_ntohs(tcp->source);
        pkt.dst_port = __constant_ntohs(tcp->dest);
    } else if (ip->protocol == IPPROTO_UDP) {
        struct udphdr *udp = (void *)ip + ip->ihl * 4;
        if ((void *)(udp + 1) > data_end) {
            return XDP_PASS;
        }
        pkt.src_port = __constant_ntohs(udp->source);
        pkt.dst_port = __constant_ntohs(udp->dest);
    }

    // 5. 传递网络包信息到用户空间
    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &pkt, sizeof(pkt));

    return XDP_PASS;
}

char __license[] SEC("license") = "Dual MIT/GPL";