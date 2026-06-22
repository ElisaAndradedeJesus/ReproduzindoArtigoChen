CLANG ?= clang
CC ?= cc
BPFTOOL ?= bpftool
IFACE ?= lo

BPF_OBJECT := flow_monitor.bpf.o
SKELETON := flow_monitor.skel.h
PROGRAM := flow_monitor

BPF_CFLAGS := -O2 -g -Wall -target bpf \
	-D__TARGET_ARCH_x86 \
	-I/usr/include/x86_64-linux-gnu

USER_CFLAGS := -O2 -g -Wall
USER_LIBS := -lbpf -lelf -lz

.PHONY: all run clean detach inspect help

all: $(PROGRAM)
	@echo "Compilação concluída: $(PROGRAM)"

$(BPF_OBJECT): flow_monitor.bpf.c common.h
	@echo "Compilando programa eBPF..."
	$(CLANG) $(BPF_CFLAGS) -c flow_monitor.bpf.c -o $@

$(SKELETON): $(BPF_OBJECT)
	@echo "Gerando skeleton..."
	$(BPFTOOL) gen skeleton $< > $@

$(PROGRAM): main.c common.h $(SKELETON)
	@echo "Compilando programa userspace..."
	$(CC) $(USER_CFLAGS) main.c -o $@ $(USER_LIBS)

run: all
	@echo "Executando na interface $(IFACE)..."
	sudo ./$(PROGRAM) $(IFACE)

inspect:
	ip link show $(IFACE)

detach:
	sudo ip link set dev $(IFACE) xdp off

clean:
	rm -f $(BPF_OBJECT) $(SKELETON) $(PROGRAM)

help:
	@echo "make                 Compila todo o projeto"
	@echo "make run IFACE=lo    Compila e executa"
	@echo "make inspect         Exibe a interface"
	@echo "make detach          Remove XDP da interface"
	@echo "make clean           Remove arquivos gerados"