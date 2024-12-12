# Base image with Go installed
FROM golang:1.20-alpine

# Install supervisor to manage multiple services
RUN apk add --no-cache supervisor

# Set the working directory
WORKDIR /app

# Copy the source code for all services
COPY ./ /app

# Copy the .env file into the container
COPY ./.env /app/.env

# Copy the supervisor config file
COPY supervisord.conf /etc/supervisord.conf

# Expose port for gateway-service
EXPOSE 8080

# Start all services using supervisord
CMD ["supervisord", "-c", "/etc/supervisord.conf"]
