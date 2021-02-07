FROM golang:1.15

WORKDIR /go/src/8k
COPY . .

ARG VERSION=dev
ENV VERSION="${VERSION}"

RUN go test -v ./...
RUN CGO_ENABLED=0 go build -o 8192 -ldflags "-X github.com/2bytes/8k/internal/config.Version=${VERSION}" cmd/8192/main.go
RUN ./8192 -v

FROM scratch
WORKDIR /
COPY pkg/frontend/tmpl /frontend
COPY --from=0 /go/src/8k/8192 .
ENTRYPOINT ["./8192"]
