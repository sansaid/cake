FROM golang:1.17 as builder

WORKDIR /app
COPY ./cake /app/.
COPY ../caked.proto /app/.
RUN cd /app && \
    protoc --go-grpc_out=./cake/pb --go-grpc_opt=paths=source_relative --go_out=./cake/pb --go_opt=paths=source_relative caked.proto \
    go get -d . \
    go build -o cakectl .

FROM golang:1.17 as main

WORKDIR /app
COPY --from=builder  /app/cakectl /app/.

ENTRYPOINT [ "cakectl" ]