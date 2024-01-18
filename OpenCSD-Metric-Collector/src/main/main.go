package main

import (
    "context"
    "log"
    "net"
	"fmt"
    "time"
	"os"
	"strconv"
    // grpc
    pb "opencsd/src/config"
    "google.golang.org/grpc"

    // influxdb v1 client
    client "github.com/influxdata/influxdb/client/v2"
)

var(
	//INSTANCE_METRIC_COLLECTOR_IP = "10.0.4.80"
	//INSTANCE_METRIC_COLLECTOR_PORT = "40802"
	
    INFLUXDB_CSD_DB = "opencsd_management_platform"
    // INFLUXDB_CSD_MEASUREMENT = "csd_metric"

    INFLUXDB_ETCD_DB = "opencsd_management_platform"
    // INFLUXDB_ETCD_MEASUREMENT = "etcd_metric"

	INFLUXDB_NODE_DB = "opencsd_management_platform"

)

type server struct {
    pb.UnimplementedCSDMetricServer // csd metric method stub
	pb.UnimplementedNodeMetricServer // node metric method stub
	// pb.UnimplementedEtcdMetricServer // etcd metric method stub

}

type CSDMetricData struct {
	Id           int32           `json:"id"`
	
	TotalCpuCapacity  int32 	 `json:"totalCpuCapacity"`
	CpuUsage		  float64 	 `json:"cpuUsage"`
	CpuUsagePercent	  float64 	 `json:"cpuUsagePercent"`

	TotalMemCapacity  int32      `json:"totalMemCapacity"`
	MemUsage		  int32 	 `json:"memUsage"`
	MemUsagePercent	  float64    `json:"memUsagePercent"`
	
	TotalDiskCapacity int32     `json:"totalDiskCapacity"`
	DiskUsage		  int32     `json:"diskUsage"`
	DiskUsagePercent  float64   `json:"diskUsagePercent"`
	
	NetworkBandwidth  int32     `json:"networkBandwidth"`
	NetworkRxData	  int32     `json:"networkRxData"`
	NetworkTxData	  int32     `json:"networkTxData"`

	CsdMetricScore    float64 	`json:"csdMetricScore"`
	CsdMetricGrade    string	
}

type NodeMetricData struct {
	Id				     int64
	
	NodeIp				 string
	NodeName			 string
	NodeStatus			 string

	TotalCpuCapacity	 int64	 
	CpuUsage			 int64

	TotalMemCapacity	 int64
	MemUsage			 int64

	TotalStorageCapacity int64
	StorageUsage		 int64

	NetworkRxData		 int64
	NetworkTxData		 int64
}

type EtcdMetricData struct {
	ID				int32
}

// Node stub handler server
func (s *server) ReceiveNodeMetric(ctx context.Context, in *pb.NodeMetricRequest) (*pb.MetricResponse, error) {
	id := in.GetId()

	nodeName := in.GetNodeName()
	nodeIp := in.GetNodeIp()
	nodeStatus := in.GetNodeStatus()

	totalCpuCapacity := in.GetTotalCpuCapacity()
    cpuUsage := in.GetCpuUsage()

	totalMemCapacity := in.GetTotalMemCapacity()
    memUsage := in.GetMemUsage()

	totalStorageCapacity := in.GetTotalDiskCapacity()
    storageUsage := in.GetDiskUsage()

	networkRxData := in.GetNetworkRxData()
	networkTxData := in.GetNetworkTxData()

	//  receive log

	var nodeMetricData NodeMetricData

	nodeMetricData.Id = id
	
	nodeMetricData.NodeIp = nodeIp
	nodeMetricData.NodeName = nodeName
	nodeMetricData.NodeStatus = nodeStatus

	nodeMetricData.TotalCpuCapacity = totalCpuCapacity
	nodeMetricData.CpuUsage = cpuUsage

	nodeMetricData.TotalMemCapacity = totalMemCapacity
	nodeMetricData.MemUsage = memUsage

	nodeMetricData.TotalStorageCapacity = totalStorageCapacity
	nodeMetricData.StorageUsage = storageUsage

	nodeMetricData.NetworkRxData = networkRxData
	nodeMetricData.NetworkTxData = networkTxData
	fmt.Printf("%+v\n", nodeMetricData)
	nodeMetricInsert(&nodeMetricData)

	// client stub에게 response
	return &pb.MetricResponse{JsonConfig: `request success`}, nil
}

// CSD stub handler server
func (s *server) ReceiveCSDMetric(ctx context.Context, in *pb.CSDMetricRequest) (*pb.MetricResponse, error) {
	// request data parsing
	id := in.GetId()

    totalCpuCapacity := in.GetTotalCpuCapacity()
    cpuUsage := in.GetCpuUsage()
    cpuUsagePercent := in.GetCpuUsagePercent()

    totalMemCapacity := in.GetTotalMemCapacity()
    memUsage := in.GetMemUsage()
    memUsagePercent := in.GetMemUsagePercent()

    totalDiskCapacity := in.GetTotalDiskCapacity()
    diskUsage := in.GetDiskUsage()
    diskUsagePercent := in.GetDiskUsagePercent()

    networkBandwidth := in.GetNetworkBandwidth()
    networkRxData := in.GetNetworkRxData()
    networkTxData := in.GetNetworkTxData()

	csdMetricScore := in.GetCsdMetricScore()

	
    // 수신 데이터 로그 출력
    // log.Printf("Received data: Id=%d, CpuUsage=%.4f, MemUsage=%.4f, NetworkSpeed=%.4f", id, cpuUsage, memUsage, networkSpeed)

	log.Printf("Received data: TotalCpuCapacity=%d, CpuUsage=%.4f, CpuUsagePercent=%.4f", totalCpuCapacity, cpuUsage, cpuUsagePercent)
	log.Printf("TotalMemCapacity=%d, MemUsage=%d, MemUsagePercent=%.4f", totalMemCapacity, memUsage, memUsagePercent)
	log.Printf("TotalDiskCapacity=%d, DiskUsage=%d, DiskUsagePercent=%.4f", totalDiskCapacity, diskUsage, diskUsagePercent)
	log.Printf("NetworkBandwidth=%d, NetworkRxData=%d, NetworkTxData=%d", networkBandwidth, networkRxData, networkTxData)
	log.Printf("CsdMetricScore=%.4f", csdMetricScore)

	// 메트릭 구조체 생성 및 초기화
	var csdMetricData CSDMetricData
	
	csdMetricData.Id = id

	csdMetricData.TotalCpuCapacity = totalCpuCapacity
	csdMetricData.CpuUsage = cpuUsage
	csdMetricData.CpuUsagePercent = cpuUsagePercent

	csdMetricData.TotalMemCapacity = totalMemCapacity
	csdMetricData.MemUsage = memUsage
	csdMetricData.MemUsagePercent = memUsagePercent

	csdMetricData.TotalDiskCapacity = totalDiskCapacity
	csdMetricData.DiskUsage = diskUsage
	csdMetricData.DiskUsagePercent = diskUsagePercent

	csdMetricData.NetworkBandwidth = networkBandwidth
	csdMetricData.NetworkRxData = networkRxData
	csdMetricData.NetworkTxData = networkTxData

	csdMetricData.CsdMetricScore = csdMetricScore
	csdMetricData.CsdMetricGrade = calcCSDScore(csdMetricScore)
    // csd metric DB에 삽입
    CSDMetricInsert(&csdMetricData)
	
    // // Instance Metric Collector에 metric 전달
	// CSDMetricSender(INSTANCE_METRIC_COLLECTOR_IP, INSTANCE_METRIC_COLLECTOR_PORT, &csdMetricData)

	// client stub에게 response
	return &pb.MetricResponse{JsonConfig: `request success`}, nil
}

func calcCSDScore(score float64) string{
	var grade string
    switch {
		case score >= 80:
			grade = "very good"
		case score >= 60:
			grade = "good"
		case score >= 40:
			grade = "fail"
		case score >= 20:
			grade = "poor"
		default:
			grade = "very poor"
    }
	return grade
}

func CSDMetricInsert(metricData *CSDMetricData) {
	INFLUX_IP := os.Getenv("INFLUX_IP")
    INFLUX_PORT := os.Getenv("INFLUX_PORT")
    INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
    INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")

	//INSTANCE_METRIC_COLLECTOR_IP = "10.0.4.80"
	//INSTANCE_METRIC_COLLECTOR_PORT = "40802"
	
	c, err := client.NewHTTPClient(client.HTTPConfig{ // InfluxDB 연결
		Addr: "http://" + INFLUX_IP + ":" + INFLUX_PORT, 
		Username: INFLUX_USERNAME,
		Password: INFLUX_PASSWORD,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{ // 배치 포인트 생성
		Database:  INFLUXDB_CSD_DB, // DB 이름
 		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	var INFLUXDB_CSD_MEASUREMENT = "csd" + strconv.Itoa(int(metricData.Id)) + "_metric" // csd 번호에 따라 테이블 설정 

	// tags := map[string]string{"flags": "csd-metric"} //etcd 정보인지 csd 정보인지 분기 (etcd : 1, csd metric : 2)
	fields := map[string]interface{}{
		"current_time": time.Now().Format("15:04:05"),
		"id": metricData.Id,

		"cpu_total": metricData.TotalCpuCapacity,
		"cpu_usage": metricData.CpuUsage,
		"cpu_percent": metricData.CpuUsagePercent,

		"memory_total":  metricData.TotalMemCapacity,
		"memory_usage": metricData.MemUsage,
		"memory_percent": metricData.MemUsagePercent,

		
		"disk_total":  metricData.TotalDiskCapacity,
		"disk_usage": metricData.DiskUsage,
		"disk_percent": metricData.DiskUsagePercent,

		"network_bandwidth":  metricData.NetworkBandwidth,
		"network_rx_byte": metricData.NetworkRxData,
		"network_tx_byte": metricData.NetworkTxData,

		"score": metricData.CsdMetricScore,
		"grade": metricData.CsdMetricGrade,
	}

	pt, err := client.NewPoint(INFLUXDB_CSD_MEASUREMENT, nil, fields, time.Now()) //measurements에 tag, field, time insert => 필요 없는 값이면 nil로 설정, time은 nil 설정이 안됨
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt) // 생성해둔 배치 포인트에 새로운 데이터 추가 

	err = c.Write(bp) //influxdb에 write
	if err != nil {
		log.Fatal(err)
	}

}

// IP, Port, username, password 

func nodeMetricInsert(metricData *NodeMetricData){
	INFLUX_IP := os.Getenv("INFLUX_IP") 
	INFLUX_PORT := os.Getenv("INFLUX_PORT")
    INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
    INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")

	c, err := client.NewHTTPClient(client.HTTPConfig{ // InfluxDB 연결
		Addr: "http://" + INFLUX_IP + ":" + INFLUX_PORT, 
		Username: INFLUX_USERNAME,
		Password: INFLUX_PASSWORD,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{ // 배치 포인트 생성
		Database:  INFLUXDB_NODE_DB, // DB 이름
 		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	var INFLUXDB_NODE_MEASUREMENT = "node" + strconv.Itoa(int(metricData.Id)) + "_metric" // csd 번호에 따라 테이블 설정 

	// tags := map[string]string{"flags": "csd-metric"} //etcd 정보인지 csd 정보인지 분기 (etcd : 1, csd metric : 2)
	fields := map[string]interface{}{
		"current_time": time.Now().Format("15:04:05"),
		"id": metricData.Id,

		"node_name": metricData.NodeName,
		"node_ip": metricData.NodeIp,
		"node_status": metricData.NodeStatus,

		"cpu_total": metricData.TotalCpuCapacity,
		"cpu_usage": metricData.CpuUsage,

		"memory_total":  metricData.TotalMemCapacity,
		"memory_usage": metricData.MemUsage,
		
		"storage_total":  metricData.TotalStorageCapacity,
		"storage_usage": metricData.StorageUsage,

		"network_rx_byte": metricData.NetworkRxData,
		"network_tx_byte": metricData.NetworkTxData,
	}

	pt, err := client.NewPoint(INFLUXDB_NODE_MEASUREMENT, nil, fields, time.Now()) //measurements에 tag, field, time insert => 필요 없는 값이면 nil로 설정, time은 nil 설정이 안됨
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt) // 생성해둔 배치 포인트에 새로운 데이터 추가 

	err = c.Write(bp) //influxdb에 write
	if err != nil {
		log.Fatal(err)
	}
}

// //80서버 instance metric collector grpc 서버 접속 및 메트릭 전송
// func CSDMetricSender(ip string, port string, metricData *CSDMetricData) {
//     // gRPC 서버에 연결
//     conn, err := grpc.Dial(ip + ":" + port, grpc.WithInsecure())
//     if err != nil {
//         log.Fatalf("Could not connect: %v", err)
//     }
//     defer conn.Close()

//     // gRPC 클라이언트 생성
//     client := pb.NewCSDMetricClient(conn)

// 	// gRPC 메서드 호출 => csd metric 기반 request 생성 및 전송 후 response 수신
// 	request := &pb.CSDMetricRequest{Id : metricData.ID, CpuUsage : metricData.CPUUSAGE, MemUsage : metricData.MEMUSAGE, NetworkSpeed: metricData.NETWORKUSAGE}
    
// 	// grpc 서버 응답
// 	response, err := client.ReceiveCSDMetric(context.Background(), request)
//     if err != nil {
//         log.Fatalf("Get Response Error From Grpc Server: %v", err)
//     }
//     fmt.Printf("Response From gRPC Server: %s\n", response.JsonConfig) //응답 형식 확인해야 함
// }

// etcd stub handler server
// func (s *server) Get(ctx context.Context, in *pb.EtcdMetricRequest) (*pb.EtcdMetricResponse, error) {
// 	// request data parsing

// 	// 수신 데이터 로그 출력

// 	// etcd metric DB에 삽입
// 	etcdMetricInsert()

// 	// client stub에게 response
// 	return &pb.EtcdMetricResponse{JsonConfig: `request success`}, nil
// }

// func etcdInsert(){

// }

func main() {
    lis, err := net.Listen("tcp", ":40801")

    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    // grpc 서버 생성
    s := grpc.NewServer()

	// stub handler 서버 등록 
	pb.RegisterCSDMetricServer(s, &server{}) // register csd metric server
	pb.RegisterNodeMetricServer(s, &server{}) // register node metric server
	// pb.RegisterEtcdMetricServer(s, &server{}) // register etcd metric server start
    
	fmt.Println("gRPC Server Created [port : 40801]")
	
	if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
