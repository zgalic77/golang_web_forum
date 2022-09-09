FROM golang:1.18.6-alpine3.15 as mainbuilder
ENV GOPATH /go
WORKDIR /go/src
COPY . /go/src/main_app
RUN cd /go/src/main_app && go build .

FROM alpine
WORKDIR /app
COPY --from=mainbuilder /go/src/main_app/main_app /app
COPY .env /app
COPY --from=mainbuilder /go/src/main_app/templates ./templates
COPY --from=mainbuilder /go/src/main_app/static ./static
EXPOSE 80
ENTRYPOINT ["./main_app"]