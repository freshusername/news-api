FROM golang:1.20.1 as builder

ARG DB_HOST
ARG DB_NAME
ARG DB_PASSWORD
ARG DB_PORT
ARG DB_USER

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

#retrieve application dependencies, copy go.mod & go.sum
COPY go.mod ./
COPY go.sum ./
RUN go mod download

#copy local code to the container image
COPY . ./ 

# Set the working directory to /app/api where the main.go file is
WORKDIR /app/api

#print packages when building the binary
# Build the Go application, and place the output binary at /app/server
RUN CGO_ENABLED=0 go build -v -o /app/server . #

FROM alpine:latest as deploy
COPY --from=builder /app/server /app/server
#COPY .env.local /app - works only locally

ENV DB_HOST=${DB_HOST}
ENV DB_NAME=${DB_NAME}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_PORT=${DB_PORT}
ENV DB_USER=${DB_USER}

EXPOSE 3000
CMD ["/app/server"]