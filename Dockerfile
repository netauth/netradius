FROM golang:1.19-alpine as build
WORKDIR /netauth/netradius
COPY . .
RUN go mod vendor && \
        CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /netradius . && \
        apk add upx binutils && \
        strip /netradius && \
        upx /netradius && \
        ls -alh /netradius

FROM scratch
LABEL org.opencontainers.image.source https://github.com/netauth/netradius
ENTRYPOINT ["/netradius"]
USER 1000
COPY --from=build /netradius /netradius
