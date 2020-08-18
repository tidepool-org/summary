FROM golang as builder
ARG MODULE
# librdkafka Build from source
RUN git clone https://github.com/edenhill/librdkafka.git

WORKDIR librdkafka

RUN ./configure --prefix /usr

RUN make

RUN make install

# Build go binary

WORKDIR /go/src/github.com/tidepool-org/summary
COPY . .

ENV GO111MODULE=on
RUN go mod download


RUN go build -o dist/main
RUN ls
# final stage
FROM ubuntu
ARG MODULE
COPY --from=builder /usr/lib/pkgconfig /usr/lib/pkgconfig
COPY --from=builder /usr/lib/librdkafka* /usr/lib/
COPY --from=builder /go/src/github.com/tidepool-org/summary/dist /dist
WORKDIR /dist
CMD ["./main"]
