FROM golang:1.12.3-stretch AS build
MAINTAINER "Diogo Xavier <diogoxavierpinto@gmail.com>"

WORKDIR /go/src/github.com/diogox/REST-JWT
COPY . .

RUN rm -rf ./vendor

# Install dep
RUN apt-get install curl git openssl
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Deploy prisma
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get -y install nodejs
RUN npm install -g prisma
RUN prisma generate
RUN sed -i 's/localhost:4467/prisma:4467/g' prisma.yml

RUN dep ensure -v

# Build smallest binary possible
RUN go build -o rest-server -ldflags="-s -w" server/cmd/main.go

# Compress binary even further
RUN apt-get install -y upx
RUN upx --brute rest-server

# Make smaller image with just the executable
#FROM alpine
#COPY --from=build rest-server /rest-server

EXPOSE 8090

CMD ["./rest-server"]
