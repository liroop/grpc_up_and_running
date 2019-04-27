// Go to ${Advanced-gRPC}/samples/ch01/uc01
// Optional: Execute protoc --go_out=plugins=grpc:golang/product_info product_info.proto
// Execute go get -v github.com/advanced-grpc/samples/ch01/uc01/golang/product_info
// Execute go run golang/product_info_server/main.go

package main

import (
	"context"
	"errors"
	"log"
	"net"

	pb "github.com/advanced-grpc/samples/ch01/productinfo/go/product_info"
	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/google/uuid"
)

const (
	port = ":50051"
)

// server is used to implement ecommerce/product_info.
type server struct{
	productMap map[string]*pb.Product
}

// AddProduct implements ecommerce.AddProduct
func (s *server) AddProduct(ctx context.Context, in *pb.Product) (*wrapper.StringValue, error) {
	out, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	in.Id = out.String()
	if (s.productMap == nil) {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[in.Id] = in
	return &wrapper.StringValue{Value: in.Id}, nil
}

// GetProduct implements ecommerce.GetProduct
func (s *server) GetProduct(ctx context.Context, in *wrapper.StringValue) (*pb.Product, error) {
	value, exists := s.productMap[in.Value]
	if exists {
		return value, nil
	}
	return nil, errors.New("Product does not exist for the ID" + in.Value)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
