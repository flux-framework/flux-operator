FROM golang:latest
RUN git clone https://github.com/researchapps/dragonboat-example /opt/dragon
WORKDIR /opt/dragon
RUN  make 
COPY ./entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
