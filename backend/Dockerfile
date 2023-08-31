FROM golang:alpine as build

# Install necessary packages and set the working directory
RUN apk add --no-cache gcc musl-dev
WORKDIR /build

# Copy the entire source code from the current directory to the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o ./app

# Compress the binary
RUN apk add upx
RUN upx ./app

# Start a new stage from scratch
FROM scratch

# Copy the compiled Go binary and other necessary files and directories
COPY --from=build /build/app /app
COPY --from=build /build/.env /.env
COPY --from=build /build/pkg/db/migrations /pkg/db/migrations
COPY --from=build /build/pkg/db/files/images/post pkg/db/files/images/post
COPY --from=build /build/pkg/db/files/images/profile pkg/db/files/images/profile

# Set the entry point for the container
EXPOSE 80
ENTRYPOINT ["/app"]
