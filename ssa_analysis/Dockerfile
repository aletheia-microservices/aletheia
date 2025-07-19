FROM golang:1.22.4

WORKDIR /app

COPY . .

WORKDIR /app/pointerto

RUN go mod tidy

# CMD ["go", "run", "main.go"]

# docker build -t pointerto . && docker run pointerto > output/app.ptrs
