FROM debian:bookworm-slim

EXPOSE 830

RUN apt-get update -yq && \
    apt-get install -yq \
        --no-install-recommends \
        ca-certificates \
        build-essential \
        binutils-x86-64-linux-gnu \
        pkg-config \
        cmake \
        libpcre2-dev \
        libssl-dev \
        libssh-dev \
        libcurl4-openssl-dev \
        git \
        vim \
        curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN echo "root:password" | chpasswd

RUN mkdir build

WORKDIR /build

RUN curl -L https://github.com/CESNET/libyang/archive/refs/tags/v3.7.8.tar.gz -o libyang.tar.gz && \
    curl -L https://github.com/CESNET/libnetconf2/archive/refs/tags/v3.5.5.tar.gz -o libnetconf2.tar.gz && \
    curl -L https://github.com/sysrepo/sysrepo/archive/refs/tags/v3.3.10.tar.gz -o sysrepo.tar.gz && \
    curl -L https://github.com/CESNET/netopeer2/archive/refs/tags/v2.2.35.tar.gz -o netopeer2.tar.gz

RUN tar -xzf libyang.tar.gz && \
    tar -xzf libnetconf2.tar.gz && \
    tar -xzf sysrepo.tar.gz && \
    tar -xzf netopeer2.tar.gz

RUN mkdir libyang-3.7.8/build && \
    mkdir libnetconf2-3.5.5/build && \
    mkdir sysrepo-3.3.10/build && \
    mkdir netopeer2-2.2.35/build

WORKDIR /build/libyang-3.7.8/build
RUN cmake .. && \
    make && \
    make install

WORKDIR /build/libnetconf2-3.5.5/build
RUN cmake .. && \
    make && \
    make install

WORKDIR /build/sysrepo-3.3.10/build
RUN cmake .. && \
    make && \
    make install

WORKDIR /build/netopeer2-2.2.35/build
RUN ldconfig && \
    cmake .. && \
    make && \
    make install

CMD ["netopeer2-server", "-d", "-v3"]
