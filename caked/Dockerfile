FROM golang:1.17 as builder

WORKDIR /app
COPY ./cake /app/.
COPY caked.proto /app/.
RUN cd /app && \
    protoc --go-grpc_out=./caked/pb --go-grpc_opt=paths=source_relative --go_out=./caked/pb --go_opt=paths=source_relative caked.proto \
    go get -d . \
    go build -o caked .

FROM golang:1.17 as main

WORKDIR /app
COPY --from=builder  /app/caked /app/.

ENTRYPOINT [ "caked" ]