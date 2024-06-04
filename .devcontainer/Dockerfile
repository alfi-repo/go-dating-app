FROM golang:1.22-alpine

RUN apk add --no-cache libxext libxrender libxtst libxi freetype procps gcompat bash curl make

RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.59.0

RUN wget -O- -nv https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh -s v3.20.0
