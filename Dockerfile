FROM golang:1.11.1
WORKDIR /go/src/github.com/99heitor/pokemon-quiz-bot/
RUN go get -d -v golang.org/x/net/html  
RUN go get -d -v gopkg.in/telegram-bot-api.v4
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
WORKDIR /root/
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /go/src/github.com/99heitor/pokemon-quiz-bot/app .
COPY --from=0 /go/src/github.com/99heitor/pokemon-quiz-bot/cmd/pokemon.csv .

CMD ["./app"]  