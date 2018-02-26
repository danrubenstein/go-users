FROM golang:1.9.2
WORKDIR /go/go-users/src
RUN	go get -v github.com/go-redis/redis\
	&& go get -v github.com/gorilla/mux
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o src .

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /go/go-users/src/src .
CMD ["./src"]  