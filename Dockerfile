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

# Build website
WORKDIR ./web
RUN npm run build
WORKDIR ./..

# Build smallest binary possible
RUN CGO_ENABLED=0 go build -o rest-server -ldflags="-s -w" server/cmd/main.go

# Compress binary even further
#RUN apt-get install -y upx
#RUN upx --brute rest-server

# Make smaller image with just the executable
FROM alpine
COPY --from=build /go/src/github.com/diogox/REST-JWT/rest-server /server/rest-server
COPY --from=build /go/src/github.com/diogox/REST-JWT/server/email_body.html /server/email_body.html
COPY --from=build /go/src/github.com/diogox/REST-JWT/web/build /web/build/

EXPOSE 8090

CMD ["server/rest-server"]
