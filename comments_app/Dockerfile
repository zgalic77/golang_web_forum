FROM golang:1.18.6-alpine3.15 as commentsbuilder
ENV GOPATH /go
WORKDIR /go/src
COPY . /go/src/comments_app
RUN cd /go/src/comments_app && go build .

FROM alpine
WORKDIR /app
COPY --from=commentsbuilder /go/src/comments_app/comments_app /app
COPY .env /app
EXPOSE 8003
ENTRYPOINT ["./comments_app"]