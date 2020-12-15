package main

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/fbnoi/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro/v2"
	"log"
)

type iRepository interface {
	findAvailableReserve(*pb.Specification) (*pb.Vessel, error)
}

type vesselRepository struct{
	vessels []*pb.Vessel
}

func (r *vesselRepository) findAvailableReserve (spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range r.vessels {
		// 这里可能会出现脏读
		// 比如说有新的订单进入并选择了此轮船，可能导致此是这个轮船此时的状态并不能满足工作需要
		// 为了保证在读取轮船信息时，没有其他进程修改了此轮船的信息，要对此行数据添加读写锁
		if vessel.Available == true &&
			vessel.GetCapacity() >= spec.GetCapacity() &&
			vessel.GetMaxWeight() >= spec.GetMaxWeight() {
			vessel.Capacity -= spec.GetCapacity()
			vessel.MaxWeight -= spec.GetMaxWeight()
			if vessel.Capacity <= 10 || vessel.MaxWeight <= 500 {
				vessel.Available = false
			}
			//更新船的目前最大容量
			return vessel, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no vessel found for the spec: %v, vessels: %v", spec, r.vessels))
}

type vesselService struct{
	repo *vesselRepository
}

func (s *vesselService) FindAvailableAndReserve(ctx context.Context, spec *pb.Specification, res *pb.Response) error {
	vessel, err := s.repo.findAvailableReserve(spec)
	if err != nil {
		return err
	}
	res.Vessel = vessel
	return nil
}

func (s *vesselService) Init() {
	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500, Available: true},
	}
	s.repo = &vesselRepository{vessels}
}

func main() {
	service := micro.NewService(micro.Name("shippy.service.vessel"))
	service.Init()
	vesselServ := &vesselService{}
	vesselServ.Init()

	if err := pb.RegisterVesselServiceHandler(service.Server(), vesselServ); err != nil {
		log.Panic(err)
	}
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}