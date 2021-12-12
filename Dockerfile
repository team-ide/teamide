FROM centos:7

EXPOSE 21080/tcp

COPY release/teamide-linux-x64 /data/teamide
RUN chmod +x /data/teamide/server
CMD cd /data/teamide && ./server


# docker build -t teamide/toolbox .


# docker run --name toolbox-18000 -m 256m -p 18000:18000 teamide/toolbox
# docker run -itd --name toolbox-18000 -m 256m -p 18000:18000 --restart=always teamide/toolbox


# docker stop toolbox-18000
# docker rm toolbox-18000