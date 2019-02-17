FROM golang:1.11.1 as builder
WORKDIR /go/src/github.com/99heitor/pokemon-quiz-bot/
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make static-build

FROM scratch
WORKDIR /root/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/99heitor/pokemon-quiz-bot/bot .
COPY --from=builder /go/src/github.com/99heitor/pokemon-quiz-bot/pokemon.csv .

CMD ["./bot"]  