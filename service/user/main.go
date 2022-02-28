package main

import (
	"user/handler"
	"user/model"
	pb "user/proto"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
)

func main() {
	model.InitRedis()
	model.InitDb()
	// Create service
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	srv := micro.NewService(
		micro.Name("user"),
		micro.Version("latest"),
		micro.Registry(reg),
	)
	// Register handler
	pb.RegisterUserHandler(srv.Server(), new(handler.User))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
