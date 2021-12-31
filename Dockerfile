FROM golang:1.17-alpine as build

WORKDIR /app
COPY . ./
ENV CGO_ENABLED 0
ENV MYSQL_HOST "$MYSQL_HOST"
ENV MYSQL_USER "$MYSQL_USER"
ENV MYSQL_PASSWORD "$MYSQL_PASSWORD"
ENV MYSQL_DBNAME "$MYSQL_DBNAME"
RUN go mod download && go mod tidy && go build -ldflags="-s -w" -o devcode .

FROM alpine:3.12
RUN apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
COPY --from=build /app/devcode /devcode
EXPOSE 3030

CMD ["sh", "-c", "/devcode"]