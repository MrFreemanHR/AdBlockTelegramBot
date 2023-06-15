FROM ubuntu:22.04 AS bot

COPY --from=wcsiu/tdlib:1.7-alpine /usr/local/include/td /usr/local/include/td
COPY --from=wcsiu/tdlib:1.7-alpine /usr/local/lib/libtd* /usr/local/lib/
COPY --from=wcsiu/tdlib:1.7-alpine /usr/lib/libssl.a /usr/local/lib/libssl.a
COPY --from=wcsiu/tdlib:1.7-alpine /usr/lib/libcrypto.a /usr/local/lib/libcrypto.a
COPY --from=wcsiu/tdlib:1.7-alpine /lib/libz.a /usr/local/lib/libz.a

RUN apt update && apt install -y build-essential libc++-dev libc++abi-dev wget libbsd-dev git screen openssh-server

RUN wget https://go.dev/dl/go1.20.4.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz
ENV GOPATH /usr/local/go
ENV PATH $GOPATH/bin:$PATH

WORKDIR /myApp

COPY . .
RUN go get
RUN go build --ldflags "-extldflags '-L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz -lbsd'" -o /root/bot main.go

ENTRYPOINT ["/bin/bash", "-c", "while true; do sleep 10; done"]