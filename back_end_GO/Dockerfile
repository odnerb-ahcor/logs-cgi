FROM golang:1.19

WORKDIR /app

COPY . .

RUN go build -o server

EXPOSE 5000

CMD [ "./server" ]