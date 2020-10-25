module github.com/micro-community/micro-users

go 1.15

require (
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/micro-community/micro-chat v0.0.0-20201025084715-19ba018fb27c
	github.com/micro/dev v0.0.0-20201023140212-49030ae8a31f
	github.com/micro/micro/v3 v3.0.0-beta.7
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/protobuf v1.25.0
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.29.0
