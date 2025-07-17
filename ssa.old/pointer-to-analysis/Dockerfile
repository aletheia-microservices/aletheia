FROM golang:1.22.4

WORKDIR /app

COPY . .

RUN go build -o pointer-to-analysis

CMD ["./pointer-to-analysis"]
