package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "code/gRPCExample/echo"


)

const (
	address = "localhost:50052"
)

func main(){
	//setup connection to gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil{
		log.Fatalf("err connecting %v", err)
	}

	defer conn.Close()
	//creates new client
	c := pb.NewEchoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &pb.EchoRequest{
		Text: "This is a text",
	})
	if err != nil {
		log.Fatalf("couldnt Get from server %v", err)
	}

	log.Printf("Message from server: %s ", r.Text)
}