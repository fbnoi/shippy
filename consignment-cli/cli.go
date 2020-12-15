package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/micro/go-micro/v2"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/fbnoi/shippy/consignment-service/proto/consignment"
)

const (
	defaultFile = "consignment.json"
)

// 读取 consignment.json 中记录的货物信息
func parseFile(fileName string) (*pb.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, nil
}

func main() {

	service := micro.NewService(micro.Name("shippy.cli.consignment"))
	service.Init()

	client := pb.NewConsignmentService("shippy.service.consignment", service.Client())

	infoFile := defaultFile

	if len(os.Args) > 1 {
		infoFile = os.Args[1]
	}
	consignment, err := parseFile(infoFile)
	// 解析货物信息
	if err != nil {
		log.Fatalf("parse info file error: %v", err)
	}

	// 调用 RPC
	// 将货物存储到我们自己的仓库里
	resp, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("create consignment error: %v", err)
	}
	// 新货物是否托运成功
	log.Printf("created: %t", resp.Created)

	resp, err = client.GetConsignments(context.Background(), &pb.GetRequest{})

	if err != nil {
		log.Fatalf("get consignments error: %v", err)
	}

	for _, consignment := range resp.GetConsignments() {
		log.Printf("%+v", consignment)
	}
}
