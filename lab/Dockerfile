FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        iperf3 \
        iproute2 \
        iputils-ping \
        libbpf1 \
        libelf1t64 \
        libzstd1 \
        procps \
        tcpdump \
        zlib1g \
    && rm -rf /var/lib/apt/lists/*

COPY t3/flow_monitor /usr/local/bin/flow_monitor

RUN chmod 0755 /usr/local/bin/flow_monitor

CMD ["sleep", "infinity"]
