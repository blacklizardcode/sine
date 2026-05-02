FROM golang:1.25.0 AS build
WORKDIR /app
COPY src/ .
RUN go build .

FROM debian:latest
COPY --from=build /app/sine /
EXPOSE 8080
ENTRYPOINT [ "/sine" ]