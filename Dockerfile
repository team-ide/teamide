FROM alpine:3.17.0

EXPOSE 21080/tcp

COPY release/teamide-server/teamide /opt/teamide/teamide
COPY release/teamide-server/conf /opt/teamide/conf
COPY release/teamide-server/libaci.so /opt/teamide/lib/libaci.so
COPY docker-entrypoint.sh /opt/teamide/docker-entrypoint.sh
RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN chmod +x /opt/teamide/docker-entrypoint.sh

ENV LD_LIBRARY_PATH=/opt/teamide/lib/

WORKDIR /opt/teamide

CMD ["/opt/teamide/docker-entrypoint.sh"]