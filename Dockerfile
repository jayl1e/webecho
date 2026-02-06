FROM golang AS builder
WORKDIR /app
COPY . .
RUN make build-dbg

FROM debian
COPY --from=builder /app/webecho /app/webecho
ENV PATH="/app:${PATH}"
EXPOSE 8080/tcp
CMD ["webecho"]
