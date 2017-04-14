FROM golang:latest
RUN go get github.com/sajari/fuzzy
RUN go get google.golang.org/grpc
CMD cd /go/src/spellchecker/ && go run main.go
