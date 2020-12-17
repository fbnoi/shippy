package main

import (
	"database/sql"
	pb "github.com/fbnoi/shippy/consignment-service/proto/consignment"
	"strconv"
)

type iRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error)
	GetAll()([]*pb.Consignment, error)
}

type ConsignmentRepository struct {
	session *sql.DB
}


func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	if err := repo.session.Ping(); err != nil {
		return nil, err
	}
	stmt, err := repo.session.Prepare("INSERT INTO consignments (`description`,`weight`, `vesselId`) VALUES (?,?,?)")
	defer stmt.Close()
	if  err != nil {
		return nil, err
	}
	ret, err := stmt.Exec(consignment.GetDescription(), consignment.GetWeight(), consignment.GetVesselId());
	if  err != nil {
		return nil, err
	}
	if lastInsertID, err := ret.LastInsertId(); err != nil {
		return nil, err
	} else {
		consignment.Id = strconv.FormatInt(lastInsertID, 10)
	}
	subStmt, err := repo.session.Prepare("INSERT INTO containers (`customerId`,`origin`, `userId`, `consignmentId`) VALUES (?,?,?,?)")
	defer subStmt.Close()
	if err != nil {
		return nil, err
	}
	for _, container := range consignment.GetContainers() {
		ret, err := subStmt.Exec(container.GetCustomerId(), container.GetOrigin(), container.GetUserId(), consignment.GetId())
		if (err != nil) {
			return nil, err
		}
		if  lastInsertID,err := ret.LastInsertId(); err==nil{
			container.Id = strconv.FormatInt(lastInsertID, 10)
		} else {
			return nil, err
		}
	}
	return consignment, nil;
}

func (repo *ConsignmentRepository) GetAll () ([]*pb.Consignment, error) {
	if err := repo.session.Ping(); err != nil {
		return nil, err
	}
	stmt, err := repo.session.Prepare("select * from consignments order by id desc")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	subStmt, err := repo.session.Prepare("select id, customerId, origin, userId from containers where `consignmentId`=? order by id desc limit 100 ")
	if err != nil {
		return nil, err
	}
	defer subStmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	consignments := make([]*pb.Consignment, 0)
	for rows.Next() {
		consignment := &pb.Consignment{}
		rows.Scan(&consignment.Id, &consignment.Description, &consignment.Weight, &consignment.VesselId)
		subRows, err := subStmt.Query(consignment.GetId())
		if err != nil {
			return nil, err
		}
		consignment.Containers = make([]*pb.Container, 0)
		for subRows.Next() {
			container := &pb.Container{}
			rows.Scan(&container.Id, &container.CustomerId, &container.Origin, &container.UserId)
			consignment.Containers = append(consignment.Containers, container)
		}
		consignments = append(consignments, consignment)
	}
	return consignments, nil
}

func (repo *ConsignmentRepository) Close() error {
	return repo.session.Close()
}