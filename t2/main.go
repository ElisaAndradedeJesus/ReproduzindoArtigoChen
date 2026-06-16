package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
)

// Estrutura correspondente à do arquivo C (Tabela I)
type FlowMetrics struct {
	SrcIP         uint32
	DstIP         uint32
	SrcPort       uint16
	DstPort       uint16
	Protocol      uint8
	PacketsPerSec uint64
	BytesPerSec   uint64
	SynCnt        uint32
	AckCnt        uint32
}

// Chave para a LPM Trie
type LpmTrieKey struct {
	PrefixLen uint32
	IPv4      uint32
}

func main() {
	// 1. Permitir que o processo trave memória no kernel (requisito eBPF)
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("Falha ao remover memlock: %v", err)
	}

	// 2. Carregar o bytecode compilado (assumindo arquivo ddos_kern.o pré-compilado)
	spec, err := ebpf.LoadCollectionSpec("ddos_kern.o")
	if err != nil {
		log.Fatalf("Falha ao carregar especificação do eBPF: %v", err)
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		log.Fatalf("Falha ao criar coleção eBPF: %v", err)
	}
	defer coll.Close()

	// 3. Atachar o programa XDP à interface de rede local (ex: eth0)
	xdpProg := coll.Programs["xdp_mitigate"]
	ifname := "eth0" // Mude para sua interface de teste
	iface, _ := net.InterfaceByName(ifname)
	
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   xdpProg,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("Erro ao atachar XDP na interface %s: %v", ifname, err)
	}
	defer l.Close()
	log.Printf("Sistema XDP de mitigação ativo na interface %s!", ifname)

	// 4. Escutar o Ring Buffer de dados coletados do Kernel
	rd, err := ringbuf.NewReader(coll.Maps["ringbuf_map"])
	if err != nil {
		log.Fatalf("Erro ao abrir ringbuf reader: %v", err)
	}
	defer rd.Close()

	blacklistMap := coll.Maps["blacklist_map"]

	// Canal para finalizar o programa de forma graciosa
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		var record ringbuf.Record
		for {
			if err := rd.Read(&record); err != nil {
				return
			}

			var metric FlowMetrics
			err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &metric)
			if err != nil {
				continue
			}

			// --- INTERFACE COM O MODELO XGBOOST ---
			// Simulando a tomada de decisão do modelo preditivo exportado (Seção III-B)
			isAttack := predictDDoS(metric)

			if isAttack {
				log.Printf("[ALERTA] DDoS Detectado do IP de Origem: %d! Bloqueando via XDP...", metric.SrcIP)
				
				// Insere o IP do atacante dinamicamente no mapa LPM Trie do Kernel
				key := LpmTrieKey{PrefixLen: 32, IPv4: metric.SrcIP}
				var value uint32 = 1 // Ação: Drop

				blacklistMap.Put(&key, &value)
			}
		}
	}()

	<-stop
	log.Println("Encerrando aplicação e limpando ganchos de kernel...")
}

// Função simuladora da inferência matemática do modelo XGBoost
func predictDDoS(m FlowMetrics) bool {
	// De acordo com o artigo, o modelo avalia taxas e flags anômalas (Ex: inundação de SYN)
	if m.PacketsPerSec > 100000 && m.SynCnt > m.AckCnt {
		return true // Classificado como tráfego malicioso
	}
	return false
}
