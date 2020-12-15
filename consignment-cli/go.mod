module consignment-cli

go 1.15

require (
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/fbnoi/shippy/consignment-service v0.0.0-00010101000000-000000000000
	github.com/micro/go-micro/v2 v2.9.1
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/sys v0.0.0-20201211090839-8ad439b19e0f // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/genproto v0.0.0-20201211151036-40ec1c210f7a // indirect
)

replace github.com/fbnoi/shippy/consignment-service => ../consignment-service
replace github.com/fbnoi/shippy/vessel-service => ../vessel-service

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
