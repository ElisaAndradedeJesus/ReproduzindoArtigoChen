FROM ubuntu:24.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        bpftool \
        build-essential \
        clang \
        libbpf-dev \
        libelf-dev \
        make \
        zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src/t3

COPY t3/Makefile t3/common.h t3/flow_monitor.bpf.c t3/main.c ./

RUN make

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        bpftool \
        ca-certificates \
        iperf3 \
        iproute2 \
        iputils-ping \
        libbpf1 \
        libelf1t64 \
        procps \
        tcpdump \
        zlib1g \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/t3/flow_monitor /usr/local/bin/flow_monitor

CMD ["sleep", "infinity"]

