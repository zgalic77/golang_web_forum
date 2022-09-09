FROM golang:1.18.6-alpine3.15 as usersbuilder
ENV GOPATH /go
WORKDIR /go/src
COPY . /go/src/users_app
RUN cd /go/src/users_app && go build .

FROM alpine
WORKDIR /app
COPY --from=usersbuilder /go/src/users_app/users_app /app
COPY .env /app
EXPOSE 81
ENTRYPOINT ["./users_app"]