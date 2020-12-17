package main

import (
	pb "github.com/fbnoi/shippy/vessel-service/proto/vessel"
	"gorm.io/gorm"
)

type iRepository interface {
	Create(*pb.Vessel) (*pb.Vessel, error)
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	db *gorm.DB
}

func (repo *VesselRepository) Create(vessel *pb.Vessel) (*pb.Vessel, error){
	if result := repo.db.Create(vessel); result.Error != nil {
		return nil, result.Error
	} else {
		return vessel, nil
	}
}

func (repo VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	vessel := &pb.Vessel{}
	if result := repo.db.Where("maxWeight > ? AND capacity > ? AND available = 1", spec.GetMaxWeight(), spec.GetCapacity()).Take(vessel); result.Error != nil {
		return nil, result.Error
	} else {
		return vessel, nil
	}
}