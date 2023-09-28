# Utilisez l'image de base Golang
FROM golang:1.21 AS builder

#Set the working directory inside the container
WORKDIR /go-cc

#Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

#Download dependencies
RUN go mod download

#Copy the entire project
COPY . .

#Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

#Start from a small Alpine image and copy the static binary
FROM alpine:latest

#Install git and zip
RUN apk --no-cache add git zip

#Set the working directory
WORKDIR /root/

#Copy the binary from the builder stage
COPY --from=builder /go-cc/main .

#Copy .env file
COPY --from=builder /go-cc/.env .env

#Expose the port the app runs on
EXPOSE 8080

#Command to run the application
CMD ["./main"]