FROM golang:1.20-alpine

RUN apk add --no-cache make git

# Set the Current Working Directory inside the container
WORKDIR /app/rinha

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Build the Go app
RUN make build

# This container exposes port 80 to the outside world
EXPOSE 80

ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=postgres
ENV DB_SSLMODE=disable

# Run the binary program produced by `go install`
CMD ["./bin/app"]