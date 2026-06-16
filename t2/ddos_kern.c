// +build ignore
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include <bpf/bpf_helpers.h>

char LICENSE[] SEC("license") = "GPL";

// Estrutura enviada para o User Space (Tabela I do artigo)
struct flow_metrics {
    __u32 src_ip;
    __u32 dst_ip;
    __u16 src_port;
    __u16 dst_port;
    __u8  protocol;
    __u64 packets_per_sec;
    __u64 bytes_per_sec;
    __u32 syn_cnt;
    __u32 ack_cnt;
};

// 1. Mapa LPM Trie para a Blacklist de mitigação rápida (XDP)
struct {
    __uint(type, BPF_MAP_TYPE_LPM_TRIE);
    __type(key, struct bpf_lpm_trie_key_u4 { __u32 prefixlen; __u32 ipv4; });
    __type(value, __u32); // 1 = Drop
    __uint(max_entries, 65536);
    __uint(map_flags, BPF_F_NO_PREALLOC);
} blacklist_map SEC(".maps");

// 2. Ring Buffer para enviar estatísticas de fluxo para o espaço de usuário
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24); // 16MB buffer
} ringbuf_map SEC(".maps");

// --- MÓDULO DE MITIGAÇÃO (XDP) ---
SEC("xdp")
int xdp_mitigate(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data     = (void *)(long)ctx->data;

    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end) return XDP_PASS;

    if (eth->h_proto != __constant_htons(ETH_P_IP)) return XDP_PASS;

    struct iphdr *iph = (void *)(eth + 1);
    if ((void *)(iph + 1) > data_end) return XDP_PASS;

    // Chave para buscar no mapa LPM Trie
    struct bpf_lpm_trie_key_u4 key = {
        .prefixlen = 32, // Match exato no IP
        .ipv4 = iph->saddr
    };

    // Verifica se o IP está na Blacklist gerada pelo XGBoost
    __u32 *action = bpf_map_lookup_elem(&blacklist_map, &key);
    if (action && *action == 1) {
        return XDP_DROP; // Mitigação imediata na placa de rede
    }

    return XDP_PASS;
}

// --- MÓDULO DE COLETA (eBPF Kprobe na pilha de rede) ---
SEC("kprobe/ip_rcv")
int kprobe_ip_rcv(struct pt_regs *ctx) {
    // Nota: Em cenários reais de produção, extrai-se o sk_buff. 
    // Para fins de replicação didática da coleta de telemetria do artigo:
    struct flow_metrics *metrics;
    
    metrics = bpf_ringbuf_reserve(&ringbuf_map, sizeof(*metrics), 0);
    if (!metrics) return 0;

    // Exemplo de preenchimento fictício simulando a extração dos metadados
    metrics->src_ip = 0x0100007F; // 127.0.0.1 simulação
    metrics->protocol = 6;        // TCP
    metrics->syn_cnt = 1;
    metrics->packets_per_sec = 180163; // Exemplo de taxa simulada

    bpf_ringbuf_submit(metrics, 0);
    return 0;
}

