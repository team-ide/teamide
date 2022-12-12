FROM alpine:3.17.0

EXPOSE 21080/tcp

COPY teamide /opt/teamide/teamide
COPY conf /opt/teamide/conf
RUN chmod +x /opt/teamide/server

WORKDIR /opt/teamide

ENTRYPOINT ["./teamide"]
# docker build -t teamide/toolbox .


# docker run --name toolbox-18000 -m 256m -p 18000:18000 teamide/toolbox
# docker run -itd --name toolbox-18000 -m 256m -p 18000:18000 --restart=always teamide/toolbox


# docker stop toolbox-18000
# docker rm toolbox-18000