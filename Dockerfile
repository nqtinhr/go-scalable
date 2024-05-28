FROM golang:alpine3.19

RUN addgroup tester && adduser -S -G tester tester

USER tester

WORKDIR /backend

COPY . .

RUN go mod download

RUN go build -o main .

CMD ["./main"]
# CMD ["go", "run", "main.go"]
