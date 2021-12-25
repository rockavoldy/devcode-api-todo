FROM golang:1.17.5-alpine as build

WORKDIR /app
COPY . ./
ENV CGO_ENABLED 0
ENV MYSQL_HOST "$MYSQL_HOST"
ENV MYSQL_USER "$MYSQL_USER"
ENV MYSQL_PASSWORD "$MYSQL_PASSWORD"
ENV MYSQL_DBNAME "$MYSQL_DBNAME"
RUN go mod download && go mod tidy && go vet . && go build -ldflags="-s -w" -o devcode .
RUN rm -rf Dockerfile go.mod go.sum

FROM gcr.io/distroless/base
COPY --from=build /app/devcode /devcode
EXPOSE 8090

CMD ["./devcode"]