FROM golang:1.22-alpine as builder

WORKDIR /app

RUN apk add --no-cache make nodejs npm ca-certificates
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy the source code into the container 
COPY ./ .

RUN make install

# Build the go application
RUN make build

FROM alpine
COPY --from=builder /app/main /
# COPY .env .env

EXPOSE 8080
ENTRYPOINT [ "./main" ]
