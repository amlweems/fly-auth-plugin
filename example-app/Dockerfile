ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk --update add git \
 && go install github.com/amlweems/fly-auth-plugin@latest

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .

FROM alpine:latest
ENV GOOGLE_APPLICATION_CREDENTIALS=/adc.json
ENV GOOGLE_EXTERNAL_ACCOUNT_ALLOW_EXECUTABLES=1
COPY entrypoint.sh /
COPY --from=builder /go/bin/fly-auth-plugin /usr/local/bin/
COPY --from=builder /run-app /usr/local/bin/
CMD ["/entrypoint.sh"]
