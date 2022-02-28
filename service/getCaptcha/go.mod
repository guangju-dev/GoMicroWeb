module getCaptcha

go 1.15

require (
	github.com/afocus/captcha v0.0.0-20191010092841-4bd1f21c8868
	github.com/asim/go-micro/plugins/registry/consul/v3 v3.7.0
	github.com/asim/go-micro/v3 v3.7.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v1.8.8
	github.com/hashicorp/golang-lru v0.5.3 // indirect
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
