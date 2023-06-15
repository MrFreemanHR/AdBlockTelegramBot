# FROM golang:1.20.4-bullseye AS golang
FROM ubuntu:22.04 AS bot

COPY --from=wcsiu/tdlib:1.7-alpine /usr/local/include/td /usr/local/include/td
COPY --from=wcsiu/tdlib:1.7-alpine /usr/local/lib/libtd* /usr/local/lib/
COPY --from=wcsiu/tdlib:1.7-alpine /usr/lib/libssl.a /usr/local/lib/libssl.a
COPY --from=wcsiu/tdlib:1.7-alpine /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
COPY --from=wcsiu/tdlib:1.7-alpine /lib/libz.a /usr/local/lib/libz.a
# RUN apk add build-base gcompat libstdc++
RUN apt update && apt install -y build-essential libc++-dev libc++abi-dev wget
RUN apt install -y libbsd-dev git

RUN wget https://go.dev/dl/go1.20.4.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz
ENV GOPATH /usr/local/go
ENV PATH $GOPATH/bin:$PATH

RUN git clone https://github.com/sqlcipher/sqlcipher.git
RUN DEBIAN_FRONTEND=noninteractive apt install -y libssl-dev tcl zlib1g-dev
RUN cd sqlcipher && ./configure --enable-tempstore=yes CFLAGS="-DSQLITE_HAS_CODEC" \
	LDFLAGS="-lcrypto" && make && make install

WORKDIR /myApp

COPY . .

# RUN go build --ldflags "-extldflags '-static -L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz'" -o /tmp/bot main.go
# go build --ldflags "-extldflags '-static -L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz -lbsd'" -o /tmp/bot main.go
# go build --ldflags "-extldflags '-static -L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz -lbsd -lmd'" -o /tmp/bot main.go
# go build --ldflags "-extldflags '-L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz -lbsd'" -o /tmp/bot main.go

RUN go get
RUN go build --ldflags "-extldflags '-L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz -lbsd'" -o /bot main.go

FROM alpine:latest AS runner
COPY --from=golang /tmp/bot /bot
ENTRYPOINT [ "/bot" ]