FROM ubuntu:22.10

RUN apt-get update \
		&& apt-get install -y ca-certificates \
		&& apt-get install -y libc6 \
		&& rm -rf /var/lib/apt/lists/* \
		&& apt-get clean \
		&& ls -a /opt

# 指定RUN工作目录
WORKDIR /home

ENV LD_LIBRARY_PATH /opt/lib

# 拷贝程序
COPY bin /opt/bin 

# 启动程序
ENTRYPOINT ["/opt/bin/multiverse-server"]