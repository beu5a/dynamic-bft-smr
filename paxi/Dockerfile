FROM golang:1.22.2-alpine
RUN apk add --no-cache \
                bind-tools \
	            bash \
                jq
        
# copy code into docker
COPY ./paxi /paxi
COPY .env /.env

COPY entry_point.sh /entry_point.sh
COPY setup_ip.sh /paxi/bin/setup_ip.sh
COPY setup_config.sh /paxi/bin/setup_config.sh
COPY start.sh /paxi/bin/start.sh

RUN chmod +x /entry_point.sh
RUN chmod +x /paxi/bin/setup_ip.sh
RUN chmod +x /paxi/bin/setup_config.sh
RUN chmod +x /paxi/bin/start.sh


WORKDIR /paxi
RUN go mod download
WORKDIR /paxi/bin
RUN chmod +x build.sh
RUN ./build.sh




# add arg server/client to entrypoint
ENTRYPOINT ["/bin/bash", "/entry_point.sh"] 