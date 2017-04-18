FROM golang:latest
RUN go get github.com/sajari/fuzzy
RUN go get google.golang.org/grpc
ADD ./ /go/src/spellchecker/
CMD cd /go/src/spellchecker/ && go run main.go
