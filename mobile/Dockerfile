FROM golang:1.22.0-alpine as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# WORKDIR $GOPATH/src/mobile
WORKDIR /app
COPY . .

RUN echo $PWD && ls -lah

# CMD go build -v
# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o ./mobile .
RUN go build -o mobile

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
# COPY --from=builder /go/src/mobile/pkg/mail/smtp/template /go/src/mobile/pkg/mail/smtp/template
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

EXPOSE 9008

CMD ["./mobile"]
