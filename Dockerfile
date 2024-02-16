# Debian base
FROM debian:bullseye

# Go compiler setup
ARG file=go1.22.0.linux-amd64.tar.gz
ARG url=https://go.dev/dl/${file}
ENV PATH="${PATH}:/home/go-compiler/go/bin"

# Copy source
COPY . /go/src
WORKDIR /go/src

# Update software list, install Go compiler, install wget, and Tex Live
RUN apt-get update && \
    apt-get -y install wget && \
    apt-get -y install texlive texinfo && \
    apt-get -y install texlive-fonts-recommended && \
    apt-get -y install texlive-fonts-extra && \
    apt-get -y install texlive-latex-extra && \
    wget ${url} && \
    rm -rf /home/go-compiler  && mkdir -p /home/go-compiler && tar -C /home/go-compiler -xzf ${file} && \
    # Compile and run the app
    go mod vendor && \
    go build -o app

# Expose port
EXPOSE 8080

# Run
CMD ["/go/src/app", "&"]
