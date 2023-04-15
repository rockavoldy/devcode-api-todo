FROM --platform=amd64 golang:1.20-alpine3.17 as build

WORKDIR /app
COPY . ./
ENV CGO_ENABLED 0
ENV MYSQL_HOST "$MYSQL_HOST"
ENV MYSQL_USER "$MYSQL_USER"
ENV MYSQL_PASSWORD "$MYSQL_PASSWORD"
ENV MYSQL_DBNAME "$MYSQL_DBNAME"
RUN go mod tidy && go vet . && go build -ldflags="-s -w" -o devcode .

FROM --platform=amd64 alpine:3.17.0
COPY --from=build /app/devcode /devcode
EXPOSE 3030

CMD ["./devcode"]