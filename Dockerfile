FROM golang:alpine AS build
RUN apk add --update --no-cache ca-certificates git
ADD . /src
RUN cd /src && CGO_ENABLED=0 go build -o hikvision_exporter


FROM alpine:latest
RUN apk add --update --no-cache ca-certificates fping
COPY --from=build /src/hikvision_exporter /
COPY --from=build /src/cameras.yml /
ENV CAMERAS=/cameras.yml
ENV LISTEN=:19101
ENV PERIOD=60
ENTRYPOINT ["/hikvision_exporter"]
