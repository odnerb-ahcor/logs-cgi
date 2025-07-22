FROM golang:1.23

WORKDIR /app

COPY . .
RUN go get golang.org/x/text/encoding/charmap
RUN go get golang.org/x/text/transform
RUN go get github.com/fsnotify/fsnotify
RUN go build -o server

EXPOSE 5000

CMD [ "./server" ]