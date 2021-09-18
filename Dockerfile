# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.14 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN mkdir pkg
COPY pkg/reddit ./pkg/reddit
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server ./cmd/reddit-rss

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:edge

ARG S6_OVERLAY_RELEASE=https://github.com/just-containers/s6-overlay/releases/latest/download/s6-overlay-amd64.tar.gz
ENV S6_OVERLAY_RELEASE=${S6_OVERLAY_RELEASE}
ADD ${S6_OVERLAY_RELEASE} /tmp/s6overlay.tar.gz
RUN apk upgrade --update --no-cache \
    && rm -rf /var/cache/apk/* \
    && tar xzf /tmp/s6overlay.tar.gz -C / \
    && rm /tmp/s6overlay.tar.gz
COPY s6 /etc/services.d
ENV S6_BEHAVIOUR_IF_STAGE2_FAILS=2

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server
EXPOSE 8080

CMD ["/init"]
