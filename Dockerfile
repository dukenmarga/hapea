# Debian base
FROM debian:bullseye

# Update software list and install wget
RUN apt-get update
RUN apt-get -y install wget

# Download and install Go compiler
ARG file=go1.22.0.linux-amd64.tar.gz
ARG url=https://go.dev/dl/${file}
RUN wget ${url}
RUN rm -rf /home/go-compiler  && mkdir -p /home/go-compiler && tar -C /home/go-compiler -xzf ${file}
ENV PATH="${PATH}:/home/go-compiler/go/bin"

# Compile and run the app
COPY . /go/src
WORKDIR /go/src
RUN go mod vendor
RUN go build -o app

# Expose port
EXPOSE 8080

# Run
CMD ["/go/src/app", "&"]
