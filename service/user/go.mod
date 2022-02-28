module user

go 1.15

require (
	github.com/asim/go-micro/plugins/registry/consul/v3 v3.7.0
	github.com/asim/go-micro/v3 v3.7.1
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v1.8.8
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/tedcy/fdfs_client v0.0.0-20200106031142-21a04994525a
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
