FROM golang:1.18.6-alpine3.15 as postsbuilder
ENV GOPATH /go
WORKDIR /go/src
COPY . /go/src/posts_app
RUN cd /go/src/posts_app && go build .

FROM alpine
WORKDIR /app
COPY --from=postsbuilder /go/src/posts_app/posts_app /app
COPY .env /app
EXPOSE 8002
ENTRYPOINT ["./posts_app"]