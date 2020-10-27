FROM golang:alpine as builder
RUN mkdir -p /go/src/rssAgregator
RUN apk update && apk upgrade && apk --no-cache add ca-certificates && update-ca-certificates && apk add git && export GOOS=linux && export GOARCH=amd64
WORKDIR /go/src/rssAgregator
COPY . .
RUN ls
RUN go get . 
RUN go build -o /compiled/rssAgregator 
COPY config.json /compiled/

FROM alpine:latest
RUN mkdir -p /go/src/rssAgregator && chmod -R 0777 /go/* && apk add --no-cache bash && apk add --no-cache tzdata
RUN apk update && apk upgrade && apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=builder /compiled/ /go/src/rssAgregator
WORKDIR /go/src/rssAgregator/
RUN chmod +x rssAgregator 
ENV PATH="/go/src/rssAgregator"
EXPOSE 8585
CMD ["rssAgregator"]
