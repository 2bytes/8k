FROM --platform=${BUILDPLATFORM} golang:1.17.0 as builder

WORKDIR /go/src/8k
COPY . .

ARG VERSION=dev
ENV VERSION="${VERSION}"

ARG TARGETOS
ARG TARGETARCH

RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o 8192 -ldflags "-X github.com/2bytes/8k/internal/config.Version=${VERSION}" cmd/8192/main.go

FROM scratch
WORKDIR /
COPY --from=builder /go/src/8k/8192 .
ENTRYPOINT ["./8192"]