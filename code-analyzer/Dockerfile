FROM golang:1.21

WORKDIR /app

# Copy the application code
COPY . .

# Handle dependencies
RUN if [ ! -f go.mod ]; then go mod init github.com/terraform-provider-kasm/code-analyzer; fi
RUN go mod tidy
RUN go get golang.org/x/tools@v0.16.0

# Build the application
RUN go build -o /app/analyzer .

# Verify the executable exists and is executable
RUN ls -la /app/analyzer && chmod +x /app/analyzer

CMD ["/app/analyzer"]
