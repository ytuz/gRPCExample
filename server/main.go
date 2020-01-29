package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"log"
	"net"
	"net/http"

	//"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "code/gRPCExample/echo"
	//"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	port = ":50052"
)

// server is used to implement echo.EchoServer.
type server struct{}

func (s *server) Get(ctx context.Context, r *pb.EchoRequest) (*pb.EchoResponse, error){
	return &pb.EchoResponse{
		Text:                 "Hello Get",
	}, nil
}

func (s *server)Update(ctx context.Context, r *pb.EchoRequest) (*pb.EchoResponse, error){
	t := "Hello update"
	if r.Text != ""{
		t = r.Text
	}
	return &pb.EchoResponse{
		Text:     t,
	}, nil
}

func runReverseProxy() error{
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//mux
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterEchoHandlerFromEndpoint(ctx, mux, "localhost:50052", opts)
	if err != nil{
		log.Fatal(err)
	}

	return http.ListenAndServe(":8082", mux)
}

func main(){
	lis, err := net.Listen("tcp", port)
	if err != nil{
		log.Fatalf("failed to listen %v", err)
	}

	//use golang provided gRPC server
	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(loggingMiddleware(), authMiddleware())))
	pb.RegisterEchoServer(s, &server{})

	log.Printf("Listening on address: %s", lis.Addr())

	go runReverseProxy()

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve server %v", err)
	}

}

func loggingMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Print("In interceptor", req)

		return handler(ctx, req)
	}
}

func authMiddleware() grpc.UnaryServerInterceptor{
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Print("In Auth middleware")

		return handler(ctx, req)
	}
}