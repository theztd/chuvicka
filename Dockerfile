## Build
FROM golang:1.19-buster AS build


COPY . /usr/src/chuvicka/

WORKDIR /usr/src/chuvicka

RUN go build -o /usr/local/bin/chuvicka


## Deploy
FROM debian:stable-slim

COPY --from=build /usr/local/bin/chuvicka /usr/local/bin/chuvicka

# Pouzivam embedFS
# COPY templates /opt/chuvicka/templates
# COPY static /opt/chuvicka/static
# COPY public /opt/chuvicka/public
# COPY v2_api.yaml /opt/chuvicka/v2_api.yaml


WORKDIR /opt/chuvicka

ENV ADDRESS=":8080"

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/chuvicka"]