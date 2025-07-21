FROM golang AS builder
WORKDIR /app
COPY . .
RUN make build-slim

FROM scratch
COPY --from=builder /app/netecho /netecho
ENV PATH="/"
EXPOSE 8080/tcp
ENTRYPOINT ["/netecho"]
