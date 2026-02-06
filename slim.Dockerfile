FROM golang AS builder
WORKDIR /app
COPY . .
RUN make build-slim

FROM scratch
COPY --from=builder /app/webecho /webecho
ENV PATH="/"
EXPOSE 8080/tcp
ENTRYPOINT ["/webecho"]
