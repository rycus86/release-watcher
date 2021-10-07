FROM golang:1.16 as builder

ARG CC=""
ARG CC_PKG=""
ARG CC_GOARCH=""

ADD . /go/src/github.com/rycus86/release-watcher

RUN if [ -n "$CC_PKG" ]; then \
      apt-get update && apt-get install -y $CC_PKG; \
    fi \
    && export CC=$CC \
    && export GOOS=linux \
    && export GOARCH=$CC_GOARCH \
    && export CGO_ENABLED=1 \
    && cd /go/src/github.com/rycus86/release-watcher \
    && go build -mod vendor -o /var/tmp/app .

FROM <target>

LABEL application="Release Watcher" \
      description="Release Watcher - Backend service to send slack notifactions after a new release of a lib" \
      version="0.0.2" \
      maintainer="Viktor Adam <rycus86@gmail.com>" \
      lastUpdatedBy="Pascal Zimmermann" \
      lastUpdatedOn="2021-03-21"

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /var/tmp/app /release-watcher

CMD [ "/release-watcher" ]
