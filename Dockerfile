FROM golang:1.13-alpine AS build-stage

RUN apk add --update make git
RUN mkdir -p /go/src/github.com/expandorg/verification
WORKDIR /go/src/github.com/expandorg/verification

COPY . /go/src/github.com/expandorg/verification
COPY .env /go/src/github.com/expandorg/verification/.env

ARG GIT_COMMIT
ARG VERSION
ARG BUILD_DATE

RUN make build-service

# Final Stage
FROM alpine

RUN apk --update add ca-certificates
RUN mkdir /app
WORKDIR /app

COPY --from=build-stage  /go/src/github.com/expandorg/verification/bin/verification .

EXPOSE 8186

CMD ["./verification"]