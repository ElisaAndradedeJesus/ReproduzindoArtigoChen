IMAGE := chen-t3-lab:latest
TOPOLOGY := chen-ddos.clab.yml
LAB := chen-ddos

.PHONY: build deploy inspect monitor test-normal test-udp shell-victim destroy help

build:
	@echo "Construindo a imagem da vítima com o t3..."
	sudo docker build -f Dockerfile -t $(IMAGE) ..

deploy: build
	@echo "Criando a topologia Containerlab..."
	sudo containerlab deploy -t $(TOPOLOGY) --reconfigure

inspect:
	sudo containerlab inspect -t $(TOPOLOGY)

monitor:
	@echo "Iniciando o t3 na interface eth1 da vítima..."
	sudo docker exec -it clab-$(LAB)-victim flow_monitor eth1

test-normal:
	@echo "Gerando um fluxo TCP normal durante 5 segundos..."
	sudo docker exec clab-$(LAB)-normal iperf3 -c 10.0.0.2 -t 5

test-udp:
	@echo "Gerando tráfego UDP controlado durante 5 segundos..."
	sudo docker exec clab-$(LAB)-attacker iperf3 -c 10.0.0.2 -u -b 10M -t 5

shell-victim:
	sudo docker exec -it clab-$(LAB)-victim bash

destroy:
	@echo "Destruindo a topologia e removendo os arquivos de estado..."
	sudo containerlab destroy -t $(TOPOLOGY) --cleanup

help:
	@echo "make deploy       Constrói a imagem e cria o laboratório"
	@echo "make inspect      Exibe o estado dos nós"
	@echo "make monitor      Executa o t3 na vítima (Ctrl+C encerra)"
	@echo "make test-normal  Gera tráfego TCP normal"
	@echo "make test-udp     Gera tráfego UDP controlado"
	@echo "make shell-victim Abre um shell na vítima"
	@echo "make destroy      Remove o laboratório"

