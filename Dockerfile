FROM golang:1.10-alpine AS build-stage

RUN apk add --update make git
RUN mkdir -p /go/src/github.com/gemsorg/boilerplate
WORKDIR /go/src/github.com/gemsorg/boilerplate

COPY . /go/src/github.com/gemsorg/boilerplate

ARG GIT_COMMIT
ARG VERSION
ARG BUILD_DATE

RUN make build-service

# Final Stage
FROM alpine

RUN apk --update add ca-certificates
RUN mkdir /app
WORKDIR /app

COPY --from=build-stage  /go/src/github.com/gemsorg/boilerplate/bin/boilerplate .

EXPOSE 3000

CMD ["./boilerplate"]