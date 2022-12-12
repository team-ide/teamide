FROM alpine:3.17.0

EXPOSE 21080/tcp

COPY release/teamide-server/teamide /opt/teamide/teamide
COPY release/teamide-server/conf /opt/teamide/conf
COPY release/teamide-server/libaci.so /opt/teamide/lib/libaci.so
RUN chmod +x /opt/teamide/teamide

ENV LD_LIBRARY_PATH=/opt/teamide/lib/

WORKDIR /opt/teamide

ENTRYPOINT ["/opt/teamide/teamide"]
# docker build -t teamide/toolbox .


# docker run --name toolbox-18000 -m 256m -p 18000:18000 teamide/toolbox
# docker run -itd --name toolbox-18000 -m 256m -p 18000:18000 --restart=always teamide/toolbox


# docker stop toolbox-18000
# docker rm toolbox-18000