FROM golang:1.22 AS builder

ARG TARGETPLATFORM
ARG VERSION=v0.7.5

WORKDIR /workspace
RUN curl -fsSL https://github.com/jmoney/tokenizers/releases/download/${VERSION}/libtokenizers.$(echo ${TARGETPLATFORM} | tr / -).tar.gz | tar xvz
COPY go.mod go.sum cmd/server/main.go ./
COPY internal ./internal
RUN go mod download
RUN mkdir -p ./lib
RUN mv ./libtokenizers.a ./lib
RUN go build -ldflags '-extldflags "-L./lib -static"' -tags lambda.norpc -o lambda .
RUN go build -ldflags '-extldflags "-L./lib -static"' -o http .

FROM alpine:3.14 AS tokenizer-server
COPY --from=builder /workspace/lambda ./lambda
ENTRYPOINT [ "./lambda" ]

FROM public.ecr.aws/lambda/provided:al2023 AS tokenizer-lambda
COPY --from=builder /workspace/http ./http
ENTRYPOINT [ "./http" ]