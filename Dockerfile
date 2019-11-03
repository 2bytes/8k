FROM golang:1.13

ARG VERSION

WORKDIR /go/src/8192bytes
COPY . .

RUN CGO_ENABLED=0 go build -o 8192 -ldflags "-X 8192bytes/internal/flags.Version=${VERSION}" cmd/8192/main.go
RUN ./8192 -v

FROM scratch
WORKDIR /
COPY pkg/frontend/tmpl /frontend
COPY --from=0 /go/src/8192bytes/8192 .
ENTRYPOINT ["./8192"]