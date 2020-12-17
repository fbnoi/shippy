package main

import (
	"context"
	pb "github.com/fbnoi/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro/v2"
	"log"
	"os"
)

type vesselService struct{
	repo *VesselRepository
}

type config struct {
	host 		string
	port 		string
	database 	string
	charset 	string
	username 	string
	password 	string
}

func (s *vesselService) FindAvailableAndReserve(ctx context.Context, spec *pb.Specification, res *pb.Response) error {
	vessel, err := s.repo.FindAvailable(spec)
	if err != nil {
		return err
	}
	res.Vessel = vessel
	return nil
}

func (s *vesselService) Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error {
	if vessel, err := s.repo.Create(req); err != nil {
		return err
	} else {
		res.Created = true
		res.Vessel = vessel
		return nil
	}
}

func (s *vesselService) Init(conf config) {
	if session, err := CreateSession(conf.host, conf.port, conf.database, conf.charset, conf.username, conf.password); err != nil {
		log.Fatalln(err)
		panic(err)
	} else {
		s.repo = &VesselRepository{session}
		vessels := []*pb.Vessel{
			&pb.Vessel{Name: "Boaty McBoatface", MaxWeight: 13080, Capacity: 378, Available: true},
			&pb.Vessel{Name: "Marry May", MaxWeight: 23080, Capacity: 480, Available: true},
			&pb.Vessel{Name: "Chaclate", MaxWeight: 7800, Capacity: 180, Available: true},
		}
		for _, vessel := range vessels {
			if _, err := s.repo.Create(vessel); err != nil {
				log.Fatalln(err)
				panic(err)
			}
		}
	}
}

func main() {

	conf := config{}

	if os.Getenv("host") 		== "" {conf.host 		= "192.168.0.194"} 	else {conf.host 	= os.Getenv("host")}
	if os.Getenv("port") 		== "" {conf.port 		= "3306"} 			else {conf.port 	= os.Getenv("port")}
	if os.Getenv("database") 	== "" {conf.database 	= "vessels"} 		else {conf.database = os.Getenv("database")}
	if os.Getenv("charset")	== "" {conf.charset 	= "utf8mb4"} 		else {conf.charset 	= os.Getenv("charset")}
	if os.Getenv("username") 	== "" {conf.username 	= "linac"} 			else {conf.username = os.Getenv("username")}
	if os.Getenv("password") 	== "" {conf.password 	= "whut123"} 		else {conf.password = os.Getenv("password")}

	service := micro.NewService(micro.Name("shippy.service.vessel"))
	service.Init()
	vesselServ := &vesselService{}
	vesselServ.Init(conf)

	if err := pb.RegisterVesselServiceHandler(service.Server(), vesselServ); err != nil {
		log.Panic(err)
	}
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}