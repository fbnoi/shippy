package main

import (
	"context"
	vesselProto "github.com/fbnoi/shippy/vessel-service/proto/vessel"
	"log"
	"os"

	pb "github.com/fbnoi/shippy/consignment-service/proto/consignment"
	micro "github.com/micro/go-micro/v2"
)

type config struct {
	host 		string
	port 		string
	database 	string
	charset 	string
	username 	string
	password 	string
}

type consignmentService struct {
	repo *ConsignmentRepository
	vesselClient vesselProto.VesselService
}

func (s *consignmentService) Init(conf config) {
	session, err := CreateSession(conf.host, conf.port, conf.database, conf.charset, conf.username, conf.password)
	if err != nil {
		panic(err)
	}
	s.repo = &ConsignmentRepository{session}
	service := micro.NewService(micro.Name("shippy.cli.consignment"))
	s.vesselClient = vesselProto.NewVesselService("shippy.service.vessel", service.Client())
}

// 添加一个托运服务
func (s *consignmentService) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	//查找符合条件的货船
	spec := &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity: int32(len(req.GetContainers())),
	}
	vessel, err := s.vesselClient.FindAvailableAndReserve(ctx, spec)
	if err != nil {
		return err
	}
	req.VesselId = vessel.GetVessel().GetId()
	// 接收承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *consignmentService) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	if consignments, err := s.repo.GetAll(); err != nil {
		return err
	} else {
		res.Consignments = consignments
		return nil
	}
}

// 服务入口
func main() {
	conf := config{}

	if os.Getenv("host") 		== "" {conf.host 		= "192.168.0.194"} 	else {conf.host 	= os.Getenv("host")}
	if os.Getenv("port") 		== "" {conf.port 		= "3306"} 			else {conf.port 	= os.Getenv("port")}
	if os.Getenv("database") 	== "" {conf.database 	= "consignments"} 	else {conf.database = os.Getenv("database")}
	if os.Getenv("charset")	== "" {conf.charset 	= "utf8mb4"} 		else {conf.charset 	= os.Getenv("charset")}
	if os.Getenv("username") 	== "" {conf.username 	= "linac"} 			else {conf.username = os.Getenv("username")}
	if os.Getenv("password") 	== "" {conf.password 	= "whut123"} 		else {conf.password = os.Getenv("password")}

	servConsignment := consignmentService{}
	servConsignment.Init(conf)

	service := micro.NewService(micro.Name("shippy.service.consignment"))

	service.Init()

	// Register service
	if err := pb.RegisterConsignmentServiceHandler(service.Server(), &servConsignment); err != nil {
		log.Panic(err)
	}
	// Run the server
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
