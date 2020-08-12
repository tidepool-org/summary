# Development
FROM golang:1.13.10-alpine AS development
WORKDIR /go/src/github.com/tidepool-org/summary
RUN adduser -D tidepool && \
    chown -R tidepool /go/src/github.com/tidepool-org/summary
USER tidepool
COPY --chown=tidepool . .
RUN ./build.sh
CMD ["./dist/main"]
