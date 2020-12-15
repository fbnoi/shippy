package main

import (
	"context"
	vesselProto "github.com/fbnoi/shippy/vessel-service/proto/vessel"
	"log"
	"strconv"
	"sync"

	pb "github.com/fbnoi/shippy/consignment-service/proto/consignment"
	micro "github.com/micro/go-micro/v2"
)

const (
	port = ":8089"
)

// 仓库接口
type iRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

//
// 存放货物的仓库， 实现了 iRepository 接口
//
type repository struct {
	iRepository
	consignments []*pb.Consignment
	lastID       int
	mux          sync.Mutex
}

// 将托运的货物放入仓库
func (r *repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	r.mux.Lock()
	id := strconv.Itoa(r.lastID)
	r.lastID++
	consignment.Id = id
	r.consignments = append(r.consignments, consignment)
	r.mux.Unlock()
	return consignment, nil
}

// 获取仓库内所有的货物
func (r *repository) GetAll() []*pb.Consignment {
	return r.consignments
}

// ********************************
// 托运服务 实现了 consignmentService 接口
//
// ********************************
type consignmentService struct {
	repo *repository
	vesselClient vesselProto.VesselService
}

func (s *consignmentService) Init() {
	s.repo = &repository{}
	s.repo.lastID = 1
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
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

// 服务入口
func main() {
	servConsignment := consignmentService{}
	servConsignment.Init()

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
