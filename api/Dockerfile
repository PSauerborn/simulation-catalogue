FROM golang:1.25 AS tests

WORKDIR /app/tests

COPY go.mod go.sum ./
RUN go mod download

RUN go install gotest.tools/gotestsum@latest

COPY tests ./tests
COPY *.go ./

CMD ["gotestsum", "--format", "testdox"]

FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o api .

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update && apt-get install -y \
    libgfortran5 \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/api .

COPY etc ./etc

CMD ["./api"]
