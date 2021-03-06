# --- local ---
# go
FROM golang:1.12-alpine as build
WORKDIR /go/load-test-sample
COPY . .
ENV GO111MODULE on
RUN apk add --no-cache git make curl && \
  curl -fLo /go/bin/air https://git.io/linux_air && \
  chmod +x /go/bin/air && \
  make build

# --- production ---
# go
FROM alpine as production
WORKDIR /app
COPY --from=build /go/load-test-sample/bin/api .
COPY --from=build /go/load-test-sample/_tools ./_tools
RUN addgroup go \
  && adduser -D -G go go \
  && chown -R go:go /app/api
CMD ["./api"]
