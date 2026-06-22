# Laboratório Containerlab do `t3`

Este laboratório cria uma rede isolada com três hosts e um switch Linux:

```text
normal (10.0.0.1)   ─┐
                     ├─ switch ─ victim (10.0.0.2)
attacker (10.0.0.3) ─┘                 └─ t3/XDP em eth1
```

O Docker compartilha o kernel do host. Por isso, o programa eBPF executado na
vítima usa o kernel da máquina real, mas observa somente a interface `eth1` do
namespace de rede daquele container.

## Pré-requisitos

```bash
docker --version
containerlab version
```

Os comandos usam `sudo`, pois o Containerlab e o carregamento de programas
eBPF/XDP precisam de privilégios administrativos.

## Primeiro uso

Na raiz do repositório:

```bash
cd lab
make deploy
make inspect
```

O primeiro `deploy` pode demorar porque baixa as imagens e instala as
dependências usadas para compilar o `t3` dentro da imagem Ubuntu.

## Executar o monitor

No primeiro terminal:

```bash
cd lab
make monitor
```

Mantenha esse terminal aberto. O comando executa `flow_monitor eth1` dentro da
vítima. Use `Ctrl+C` para encerrar somente o monitor.

## Gerar tráfego

Em outro terminal, gere um fluxo TCP comum:

```bash
cd lab
make test-normal
```

Para gerar UDP controlado a 10 Mbit/s durante cinco segundos:

```bash
make test-udp
```

Esses testes servem apenas para validar a coleta. O `t3` ainda não classifica
nem bloqueia ataques.

## Entrar na vítima

```bash
make shell-victim
```

Alguns comandos úteis dentro do container:

```bash
ip -br addr
ip -details link show eth1
bpftool prog show
bpftool map show
```

## Encerrar o laboratório

```bash
make destroy
```

Isso remove os containers, os links virtuais e o diretório de estado criado
pelo Containerlab.

