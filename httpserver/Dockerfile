FROM golang:1.16-alpine AS build
RUN apk add --no-cache git

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

WORKDIR /go/src/httpserver/
COPY . .
RUN  go mod tidy && go build -o /bin/httpserver

FROM scratch
COPY --from=build /bin/httpserver /bin/httpserver
ENTRYPOINT ["/bin/httpserver"]
CMD ["help"]
