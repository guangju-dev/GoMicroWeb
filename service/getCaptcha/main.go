package main

import (
	"getCaptcha/handler"
	pb "getCaptcha/proto"

	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	// "github.com/micro/micro/v3/service"
	// "github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	srv := micro.NewService(
		micro.Name("getCaptcha"),
		micro.Version("latest"),
		micro.Registry(reg),
	)
	// Register handler
	pb.RegisterGetCaptchaHandler(srv.Server(), new(handler.GetCaptcha))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
