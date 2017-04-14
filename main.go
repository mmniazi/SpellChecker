package main

import (
	"github.com/sajari/fuzzy"
	"net"
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "./proto"
	"google.golang.org/grpc/reflection"
	"fmt"
	"strings"
	"os"
	"bufio"
	"regexp"
)

const port = ":50051"

var model *fuzzy.Model

type server struct{}

func (s *server) Check(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	words:= strings.Fields(in.Sentence)

	for i, word := range words {
		corrected:= model.SpellCheck(word)
		if corrected != "" {
			words[i] = corrected
		}
	}
	log.Println(words)
	return &pb.Reply{Sentence: strings.Join(words, " ")}, nil
}

func main() {
	model = fuzzy.NewModel()

	// This expands the distance searched, but costs more resources (memory and time).
	// For spell checking, "2" is typically enough, for query suggestions this can be higher
	model.SetDepth(2)


	fmt.Println("Started Training")
	model.Train(LoadSampleEnglish())
	fmt.Println("Finished Training")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server started on localhost:50051")

	s := grpc.NewServer()
	pb.RegisterSpellCheckerServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func LoadSampleEnglish() []string {
	var out []string
	cwd , _ := os.Getwd()
	file, err := os.Open(cwd + "/data/big.txt")
	if err != nil {
		fmt.Println(err)
		return out
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		exp, _ := regexp.Compile("[a-zA-Z]+")
		words := exp.FindAll([]byte(scanner.Text()), -1)
		for _, word := range words {
			if len(word) > 1 {
				out = append(out, strings.ToLower(string(word)))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	return out
}
