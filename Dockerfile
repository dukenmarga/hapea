# Use the official Go image as the base image
FROM golang:1.22.8-bullseye

# Update software list, install Go compiler, install wget, and Tex Live
RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && \
    apt-get -y install wget && \
    apt-get -y install texlive texinfo && \
    apt-get -y install texlive-fonts-recommended && \
    apt-get -y install texlive-fonts-extra && \
    apt-get -y install texlive-latex-extra && \
    apt-get -y install texlive-luatex

# Expose port
EXPOSE 8080

# Set the working directory inside the container
WORKDIR /goapp

# Copy the application code
COPY . .

# Download the dependencies
RUN go mod download

# Build the Go application
RUN go build -o main .

# Run
ENV GIN_MODE=release
CMD ["./main"]