FROM golang:1.19-bullseye

# Install necessary packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        bind9-utils \
        bind9-host \
        bash \
        jq && \
    rm -rf /var/lib/apt/lists/*
        
# copy code into docker
COPY ./PaxiBFT /PaxiBFT
COPY .env /.env

COPY entry_point.sh /entry_point.sh
COPY setup_ip.sh /PaxiBFT/bin/setup_ip.sh
COPY setup_config.sh /PaxiBFT/bin/setup_config.sh
COPY start.sh /PaxiBFT/bin/start.sh

RUN chmod +x /entry_point.sh
RUN chmod +x /PaxiBFT/bin/setup_ip.sh
RUN chmod +x /PaxiBFT/bin/setup_config.sh
RUN chmod +x /PaxiBFT/bin/start.sh



WORKDIR /PaxiBFT
RUN go mod tidy
RUN go mod download
RUN go install github.com/salemmohammed/PaxiBFT
WORKDIR /PaxiBFT/bin
RUN chmod +x build.sh
RUN ./build.sh




# add arg server/client to entrypoint
ENTRYPOINT ["/bin/bash", "/entry_point.sh"] 