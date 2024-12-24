# Use the official Golang image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.0/migrate.linux-amd64.tar.gz | tar xvz && mv ./migrate /usr/local/bin

RUN git config --global --add safe.directory /app

# Download Go dependencies
RUN go mod download

# Copy the entrypoint script into the container
COPY docker-entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/entrypoint.sh"]

# Expose port 8080

EXPOSE 8080