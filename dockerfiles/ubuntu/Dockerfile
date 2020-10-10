FROM concourse/golang-builder as builder
COPY . /src
WORKDIR /src
ENV CGO_ENABLED 0
RUN go build -o /assets/in ./cmd/in
RUN go build -o /assets/check ./cmd/check
RUN set -e; for pkg in $(go list ./...); do \
    go test -o "/tests/$(basename $pkg).test" -c $pkg; \
  done

FROM ubuntu:bionic AS resource
RUN apt-get update \
      && DEBIAN_FRONTEND=noninteractive \
      apt-get install -y --no-install-recommends \
        ca-certificates \
        tzdata \
      && rm -rf /var/lib/apt/lists/*
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
ARG GITHUB_TOKEN
COPY --from=builder /tests /go-tests
WORKDIR /go-tests
RUN set -e; for test in /go-tests/*.test; do \
    $test; \
  done

FROM resource
