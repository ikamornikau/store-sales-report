FROM golang:1.27.3 AS build

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/server ./cmd

FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=build --chown=nonroot:nonroot /out/server /server
COPY --chown=nonroot:nonroot configs/example.yaml /secrets/config.yaml

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/server"]
