package main

import (
	"context"
	"net"
	"os"
	"time"

	lrf "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	pb "github.com/marceloaguero/microservices-istio/pb"
)

type greetingServiceServer struct {
}

var (
	greetings []*pb.Greeting
)

func (s *greetingServiceServer) Greeting(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	greetings = nil

	tmpGreeting := pb.Greeting{
		Id:      uuid.New().String(),
		Service: "Service-C",
		Message: "Konnichiwa, from Service-C!",
		Created: time.Now().Local().String(),
	}

	greetings = append(greetings, &tmpGreeting)

	callMongoDB(tmpGreeting)

	return &pb.GreetingResponse{
		Greeting: greetings,
	}, nil
}

func callMongoDB(greeting pb.Greeting) {
	log.Info(greeting)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_CONN")))
	if err != nil {
		log.Error(err)
	}

	defer client.Disconnect(nil)

	collection := client.Database("service-g").Collection("greetings")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, greeting)
	if err != nil {
		log.Error(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func init() {
	formatter := lrf.Formatter{ChildFormatter: &log.JSONFormatter{}}
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(getEnv("LOG_LEVEL", "info"))
	if err != nil {
		log.Error(err)
	}
	log.SetLevel(level)
}

func main() {
	lis, err := net.Listen("tcp", getEnv("PORT", ":50051"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreetingServiceServer(s, &greetingServiceServer{})
	log.Fatal(s.Serve(lis))
}
