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

Os enlaces usam MTU 1500. O Containerlab normalmente cria links `veth` com MTU
9500, valor que impedia o programa XDP atual de ser anexado.

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
dependências da imagem Ubuntu. O comando primeiro usa `t3/Makefile` para
compilar o programa na máquina host e depois copia o executável para a imagem da
vítima.

O fluxo interno do comando é:

```text
t3/Makefile → flow_monitor → imagem Docker → topologia Containerlab
```

## Executar o monitor

No primeiro terminal:

```bash
cd lab
make monitor
```

Mantenha esse terminal aberto. O comando executa `flow_monitor eth1` dentro da
vítima. Use `Ctrl+C` para encerrar somente o monitor.

## Gerar tráfego

Em outro terminal, gere primeiro um fluxo TCP pequeno e controlado:

```bash
cd lab
sudo docker exec clab-chen-ddos-normal \
  iperf3 -c 10.0.0.2 -n 128K
```

No terminal do monitor devem aparecer fluxos com origem `10.0.0.1`, destino
`10.0.0.2`, porta de destino `5201` e protocolo `6` (TCP).

Os alvos abaixo geram mais tráfego e devem ser usados somente depois do teste
pequeno:

```bash
make test-normal
make test-udp
```

`test-normal` gera TCP durante cinco segundos. `test-udp` gera UDP a 10 Mbit/s
durante cinco segundos. Como a versão atual emite um evento por pacote, esses
testes podem produzir muitas linhas no terminal.

Esses testes servem apenas para validar a coleta. O `t3` ainda não classifica
nem bloqueia ataques.

## Entender a saída

Uma linha como:

```text
Flow: 10.0.0.1:45010 -> 10.0.0.2:5201 (Proto: 6)
```

indica um fluxo TCP do cliente normal para o servidor `iperf3` na vítima. Os
totais e flags são acumulados no `flow_map`; pacotes por segundo, bytes por
segundo e tamanho médio são calculados no userspace.

As primeiras taxas podem ser muito altas porque o intervalo entre os primeiros
pacotes é extremamente curto. A implementação de janelas temporais corrigirá
essa instabilidade.

## Entrar na vítima

```bash
make shell-victim
```

Alguns comandos úteis dentro do container:

```bash
ip -br addr
ip -details link show eth1
```

Após encerrar o monitor com `Ctrl+C`, confirme a desanexação:

```bash
sudo docker exec clab-chen-ddos-victim \
  sh -c "ip -details link show dev eth1 | grep -i 'prog/xdp' || echo 'XDP desanexado'"
```

## Solução de problemas

### `Numerical result out of range` ao anexar XDP

Confirme que todos os links em `chen-ddos.clab.yml` possuem:

```yaml
mtu: 1500
```

Depois recrie a topologia com `make deploy`.

### Muitos eventos no terminal

Isso é esperado na versão atual, pois o ring buffer recebe uma fotografia das
métricas depois de cada pacote. Encerre com `Ctrl+C` e prefira o teste TCP
pequeno enquanto as janelas temporais não forem implementadas.

## Encerrar o laboratório

```bash
make destroy
```

Isso remove os containers, os links virtuais e o diretório de estado criado
pelo Containerlab.
