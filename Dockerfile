# Debian base
FROM debian:bullseye

# Go compiler setup
ARG file=go1.22.0.linux-amd64.tar.gz
ARG url=https://go.dev/dl/${file}
ENV PATH="${PATH}:/home/go-compiler/go/bin"

# Update software list, install Go compiler, install wget, and Tex Live
RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && \
    apt-get -y install wget && \
    apt-get -y install texlive texinfo && \
    apt-get -y install texlive-fonts-recommended && \
    apt-get -y install texlive-fonts-extra && \
    apt-get -y install texlive-latex-extra && \
    apt-get -y install texlive-luatex && \
    wget ${url} && \
    rm -rf /home/go-compiler  && mkdir -p /home/go-compiler && tar -C /home/go-compiler -xzf ${file}

# Copy source
COPY . /home/user1/
WORKDIR /home/user1/

# Add non-root user
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod \
    go mod vendor && \
    useradd -u 8877 user1 && \
    mkdir -p /home/user1/files && \
    chown -R user1:user1 /home/user1

USER user1

# Compile and run the app
RUN go build -o app

# Expose port
EXPOSE 8080

# Use /data as base where files directory will be
WORKDIR /home/user1/

# Run
CMD ["/home/user1/app", "&"]
