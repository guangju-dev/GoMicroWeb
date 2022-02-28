package utils

import (
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
)

func InitMicro() micro.Service {
	reg := consul.NewRegistry()
	return micro.NewService(
		micro.Registry(reg),
	)
}
