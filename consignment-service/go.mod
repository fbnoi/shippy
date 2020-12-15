module github.com/fbnoi/shippy/consignment-service

go 1.15

replace github.com/fbnoi/shippy/vessel-service => ../vessel-service

require (
	github.com/fbnoi/shippy/vessel-service v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.3
	github.com/micro/go-micro/v2 v2.9.1
	google.golang.org/protobuf v1.25.0
)
