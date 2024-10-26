# Use the official Go image
FROM golang:1.22.3

# Set the working directory
WORKDIR /app

# Copy and download Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application with the specified output name
# Specify the main file location if it's in a subdirectory, e.g., cmd/main.go
RUN go build -o server ./cmd/main.go

# Specify the command to run the compiled binary
CMD ["./server"]
