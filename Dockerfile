FROM golang:1.13

# -- Recompile openssl with weak ssl cipher support --
# Get openssl dependencies
RUN apt-get update
RUN apt-get install wget build-essential zlib1g-dev -y
## Get openssl source
WORKDIR /usr/local/src/
RUN wget https://www.openssl.org/source/openssl-1.1.1f.tar.gz
RUN tar -xf openssl*
WORKDIR /usr/local/src/openssl-1.1.1f
# Compile OpenSSL
RUN ./config enable-weak-ssl-ciphers enable-ssl3 enable-ssl3-method --prefix=/usr/local/ssl --openssldir=/usr/local/ssl shared zlib
RUN make
RUN make test
RUN make install

# Remove old openssl
RUN apt-get remove -y openssl
ENV PATH="${PATH}:/usr/local/ssl/bin"

# Add env var for library path
ENV LD_LIBRARY_PATH=/usr/local/ssl/lib