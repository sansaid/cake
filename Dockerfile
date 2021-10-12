FROM golang:1.17-bullseye as base

COPY . .

RUN go get -v golang.org/x/tools/gopls

CMD ["bash"]