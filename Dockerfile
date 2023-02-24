FROM teamide/build:0.8

EXPOSE 21080/tcp

COPY release/teamide-server/teamide /opt/teamide/teamide
COPY release/teamide-server/conf /opt/teamide/conf
COPY release/teamide-server/lib /opt/teamide/lib
COPY docker/docker-entrypoint.sh /opt/teamide/docker-entrypoint.sh
COPY docker/server.sh /opt/teamide/server.sh

ENV LD_LIBRARY_PATH=/opt/teamide/lib:$LD_LIBRARY_PATH
RUN chmod +x /opt/teamide/server.sh
RUN chmod +x /opt/teamide/docker-entrypoint.sh
# RUN yum install -y unixODBC libtool unixODBC-devel

WORKDIR /opt/teamide

CMD ["/opt/teamide/docker-entrypoint.sh"]


# cd html
# install
# npm run build
# go test -v -timeout 3600s -run ^TestStatic$ teamide/internal/static
# go build -ldflags "-w -s -X main.buildFlags=--isServer -X teamide/pkg/util.version=1.8.9" .