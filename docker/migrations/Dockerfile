FROM golang:1.23-alpine

# Install golang-migrate
RUN apk add --no-cache curl bash && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Set the working directory for migrations
WORKDIR /migrations

# Default entrypoint shows migrate version and keeps the container alive
ENTRYPOINT ["/bin/sh", "-c", "migrate -version && migrate -path=./ -database=$DB_DSN up"]