FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    curl \
    wget \
    vim \
    git \
    build-essential \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/tokamak-network/tokamak-thanos.git
RUN cd tokamak-thanos && git fetch --all && git checkout OR-1948-make-L1-contracts-deploy-script
