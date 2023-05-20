#build image
FROM docker.io/golang:1.17.4 AS builder 

RUN go version

# Set the current working directory inside the container
WORKDIR /go/src/app

COPY . .

RUN go build -mod vendor main.go

#product imager 
FROM docker.io/golang:1.17.4

WORKDIR /go/src/app

COPY --from=builder /go/src/app/main .

# This container exposes HTTP/HTTPS port to the outside world
EXPOSE 8080 10443 443

CMD ["/go/src/app/main"]
