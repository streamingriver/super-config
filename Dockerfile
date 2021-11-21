FROM golang:1.17-alpine as compiler

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /super-config


FROM scratch
WORKDIR /
COPY --from=compiler /super-config /

EXPOSE 80

ENTRYPOINT ["/super-config"]

