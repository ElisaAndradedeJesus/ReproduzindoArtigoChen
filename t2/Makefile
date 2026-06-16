IFACE ?= lo
PROG ?= xdp_ddos_protection
MAP ?= rate_limit_map

.PHONY: all compile clean attach detach dump iface-inspect

all: detach clean compile attach dump

compile: clean
	clang -O2 -g -target bpf -c ddos_kern.c -o ddos_kern.o

train-ia:
	pip install xgboost sklearn nyoka pandas numpy
	python3 train_model.py


clean:
# 	rm -f xdp_ddos_protection.o
	rm -f $(PROG).o

attach:
	sudo ip link set dev $(IFACE) xdp obj $(PROG).o sec xdp

detach:
	sudo ip link set dev $(IFACE) xdp off

iface-inspect:
	sudo ip link show $(IFACE)

dump:
	sudo bpftool map dump name $(MAP) &> /dev/null