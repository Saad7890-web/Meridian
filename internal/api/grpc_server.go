package api

import (
	"log"
	"net"
	"strconv"

	pb "github.com/Saad7890-web/meridian/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPCServer(port int, kvService *KVService) {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	pb.RegisterKVServer(server, kvService)

	
	reflection.Register(server)

	log.Printf("gRPC server running on :%d\n", port)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}