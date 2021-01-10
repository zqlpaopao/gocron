module etcdCode1

go 1.14

replace github.com/coreos/bbolt v1.3.5 => go.etcd.io/bbolt v1.3.5

replace google.golang.org/grpc v1.29.1 => google.golang.org/grpc v1.26.0

require (
	github.com/coreos/bbolt v1.3.5 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/google/uuid v1.1.4 // indirect
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/grpc v1.29.1 // indirect
)