FROM teamide/build:0.6

EXPOSE 21080/tcp

COPY release/teamide-server/teamide /opt/teamide/teamide
COPY release/teamide-server/conf /opt/teamide/conf
COPY release/teamide-server/libaci.so /opt/teamide/lib/libaci.so
COPY docker-entrypoint.sh /opt/teamide/docker-entrypoint.sh

RUN chmod +x /opt/teamide/docker-entrypoint.sh
# RUN yum install -y unixODBC libtool unixODBC-devel

WORKDIR /opt/teamide

CMD ["/opt/teamide/docker-entrypoint.sh"]


# cd html
# install
# npm run build
# go test -v -timeout 3600s -run ^TestStatic$ teamide/internal/static
# go build -ldflags "-w -s -X main.buildFlags=--isServer -X teamide/pkg/util.version=1.8.9" .