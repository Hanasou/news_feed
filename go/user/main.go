package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/Hanasou/news_feed/go/common/grpc/userpb"
	"github.com/Hanasou/news_feed/go/user/config"
	"github.com/Hanasou/news_feed/go/user/core"
	"github.com/Hanasou/news_feed/go/user/server/grpc_server"
	"google.golang.org/grpc"
)

func createServer(config *config.UserServiceConfig, userService *core.UserService) {
	switch config.Server.Type {
	case "grpc":
		createGrpcServer(config, userService)
	default:
		log.Fatalf("Unsupported server: %s", config.Server.Type)
	}

}

func createGrpcServer(config *config.UserServiceConfig, userService *core.UserService) {
	serviceUrl := config.Server.Host + ":" + strconv.Itoa(config.Server.Port)
	lis, err := net.Listen("tcp", serviceUrl)
	if err != nil {
		log.Fatalln("Failed to listen to grpc service: ", err)
	}
	log.Println("Connected to: ", serviceUrl)

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, grpc_server.NewGrpcUserServer(userService))

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to serve: ", err)
	}
}

func main() {
	fmt.Println("Hello, from Users service!")
	config, err := config.InitConfig()
	if err != nil {
		log.Fatalln("Could not initialize configuraiton file: ", err)
	}
	userService, err := core.InitializeService(config)
	if err != nil {
		log.Fatalln("Could not initialize user service: ", err)
	}
	createServer(config, userService)
}
