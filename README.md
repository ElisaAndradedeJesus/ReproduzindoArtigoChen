# Reprodução do artigo de Chen et al. com eBPF/XDP

Este repositório acompanha a reprodução do artigo *Efficient DDoS Detection
and Mitigation in Cloud Data Centers Using eBPF and XDP*.

O protótipo ativo está em `t3/`. Atualmente, ele identifica fluxos IPv4
TCP/UDP por cinco tuplas, mantém métricas em um mapa BPF e envia eventos do
kernel para o espaço de usuário por um ring buffer.

## Estrutura

```text
.
├── t1/       # Primeiras tentativas e materiais exploratórios
├── t2/       # Protótipo intermediário
├── t3/       # Implementação ativa em C, libbpf e XDP
│   ├── common.h
│   ├── flow_monitor.bpf.c
│   ├── main.c
│   └── Makefile
├── lab/      # Ambiente isolado com Docker e Containerlab
├── 02_Chen et al. - 2024 - Efficient DDoS Detection and Mitigation in Cloud Data Centers Using eBPF and XDP.pdf
└── README.md
```

## Requisitos

O projeto foi preparado para Linux x86-64. Em Ubuntu/Debian, instale:

```bash
sudo apt update
sudo apt install clang llvm bpftool libbpf-dev libelf-dev zlib1g-dev \
    build-essential make iproute2
```

Também é necessário um kernel Linux com suporte a eBPF e XDP. O carregamento
do programa requer privilégios administrativos.

## Compilação

Na raiz do repositório:

```bash
cd t3
make
```

A compilação produz três arquivos gerados:

- `flow_monitor.bpf.o`: bytecode eBPF;
- `flow_monitor.skel.h`: skeleton criado pelo `bpftool`;
- `flow_monitor`: executável do espaço de usuário.

Esses arquivos não devem ser versionados. Eles são ignorados pelo
`.gitignore` e podem ser recriados com `make`.

## Execução

Liste as interfaces disponíveis:

```bash
ip -br link
```

Para um primeiro teste local na interface de loopback:

```bash
cd t3
make run IFACE=lo
```

Em um segundo terminal, inicie um servidor TCP local:

```bash
python3 -m http.server 8000
```

Em um terceiro terminal, faça uma requisição:

```bash
curl http://127.0.0.1:8000
```

Use `Ctrl+C` para encerrar o monitor. O programa remove o link XDP durante a
finalização normal. Se uma execução terminar de forma inesperada, remova-o
manualmente:

```bash
make detach IFACE=lo
```

Para monitorar outra interface, substitua `lo`, por exemplo:

```bash
make run IFACE=eth0
```

## Comandos do Makefile

| Comando | Função |
|---|---|
| `make` | Compila o eBPF, gera o skeleton e compila o userspace |
| `make run IFACE=lo` | Compila e executa na interface indicada |
| `make inspect IFACE=lo` | Mostra informações da interface |
| `make detach IFACE=lo` | Remove um programa XDP da interface |
| `make clean` | Remove todos os arquivos gerados |
| `make help` | Lista os comandos disponíveis |

## Ambiente com Containerlab

O diretório `lab/` contém uma topologia isolada com um cliente normal, um
gerador de tráfego, um switch Linux e uma vítima que executa o `t3`.

```bash
cd lab
make deploy
make monitor
```

As instruções completas estão em [`lab/README.md`](lab/README.md).

## Estado atual e limitações

O `t3` ainda é um protótipo em desenvolvimento:

- processa apenas IPv4 TCP e UDP;
- envia um evento ao userspace para cada pacote;
- ainda não implementa janelas temporais e expiração de fluxos;
- ainda não possui classificador de aprendizado de máquina;
- ainda não possui blacklist nem mitigação com `XDP_DROP`.

Não execute testes de ataque em redes ou sistemas sem autorização. Use um
ambiente local ou isolado controlado pela equipe.

## Colaboração

Antes de enviar alterações:

```bash
cd t3
make clean
make
git status
```

Confirme que somente código-fonte e documentação aparecem no commit. O
skeleton, o objeto BPF e o executável sempre devem ser gerados localmente.
