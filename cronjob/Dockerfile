FROM golang:1.23.5-alpine as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# WORKDIR $GOPATH/src/cronjob
WORKDIR /app
COPY . .

RUN echo $PWD && ls -lah

RUN go build -o cronjob

#####################
# MAKE SMALL BINARY #
#####################
FROM alpine:3.12

RUN apk update && apk add --no-cache tzdata
ENV TZ=UTC

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the executable.
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

EXPOSE 9100

CMD ["./cronjob"]
