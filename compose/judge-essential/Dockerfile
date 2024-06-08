FROM ghcr.io/minhnhatnoe/isolate:latest

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Etc/UTC

RUN echo 'deb-src http://archive.ubuntu.com/ubuntu/ jammy main' >> /etc/apt/sources.list
RUN apt-get update
RUN apt-get build-dep -y python3
RUN apt-get install --no-install-recommends -y \
    ca-certificates gcc g++ wget && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Install Go 1.22.4
RUN wget -q --no-check-certificate https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.22.4.linux-amd64.tar.gz
RUN rm go1.22.4.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# Install python 3.10.14
RUN wget -q https://www.python.org/ftp/python/3.10.14/Python-3.10.14.tgz
RUN tar -xzf Python-3.10.14.tgz
RUN rm Python-3.10.14.tgz
WORKDIR /app/Python-3.10.14
RUN ./configure --enable-optimizations
RUN make -s -j 4
RUN make -s install
WORKDIR /app
RUN rm -rf Python-3.10.14

# Install python 2.7.18
RUN wget -q https://www.python.org/ftp/python/2.7.18/Python-2.7.18.tgz
RUN tar -xzf Python-2.7.18.tgz
RUN rm Python-2.7.18.tgz
WORKDIR /app/Python-2.7.18
RUN ./configure --enable-optimizations
RUN make -s -j 4
RUN make -s install
WORKDIR /app
RUN rm -rf Python-2.7.18

# Install PyPy 2.7
RUN wget -q --no-check-certificate https://downloads.python.org/pypy/pypy2.7-v7.3.16-linux64.tar.bz2
RUN tar -C /usr/local -xjf pypy2.7-v7.3.16-linux64.tar.bz2
RUN rm pypy2.7-v7.3.16-linux64.tar.bz2
ENV PATH=$PATH:/usr/local/pypy2.7-v7.3.16-linux64/bin

# Install PyPy 3.10
RUN wget -q --no-check-certificate https://downloads.python.org/pypy/pypy3.10-v7.3.16-linux64.tar.bz2
RUN tar -C /usr/local -xjf pypy3.10-v7.3.16-linux64.tar.bz2
RUN rm pypy3.10-v7.3.16-linux64.tar.bz2
ENV PATH=$PATH:/usr/local/pypy3.10-v7.3.16-linux64/bin

COPY compose/judge/entrypoint /entrypoint
RUN sed -i 's/\r$//g' /entrypoint
RUN chmod +x /entrypoint

ENTRYPOINT ["/entrypoint"]
