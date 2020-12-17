module github.com/fbnoi/shippy/vessel-service

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.9.1
	google.golang.org/protobuf v1.25.0
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.8
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
