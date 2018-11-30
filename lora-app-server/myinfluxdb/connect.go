package myinfluxdb

import (
	"fmt"
	"os"
	"sync"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

//DataStruct 通道结构体
type DataStruct struct {
	Tag  map[string]string
	Data map[string]interface{}
}

//DataSaveChan 数据通道
var DataSaveChan = make(chan DataStruct, 1000)

//GetDataSavechan 返回传输管道
func GetDataSavechan() chan DataStruct {
	return DataSaveChan
}

//Connect 连接数据库，等待数据并保存
func Connect() {
	go MQTTRX()
	Conn, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "",
		Password: "",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't connect to myinluxdb!!!!")
		return
	}
	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  "entre",
		Precision: "ms",
	})

	if err != nil {
		// log.Fatal(err)
		fmt.Println("inluxdb error : ", err)
		return
	}
	var lock sync.Mutex
	dataNum := 0
	dataNum1 := 0
	go func() {
		for {
			var buffer [512]byte
			_, err := os.Stdin.Read(buffer[:])
			if err != nil {
				fmt.Println("read error:", err)
			}
			fmt.Println("num is clear:", dataNum1)
			dataNum1 = 0
			fmt.Println("num is clear:", dataNum1)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * 1)
			if dataNum != 0 {
				lock.Lock()
				//这里可能有问题 当网络断开的时候
				if err := Conn.Write(bp); err != nil {
					fmt.Println("inluxdb error : ", err)
					lock.Unlock()
					continue
				}
				dataNum = 0
				lock.Unlock()
			}
		}
	}()
	for data := range DataSaveChan {
		dataNum1++
		chanOK <- "OK"
		lock.Lock()
		dataNum++
		pt, err := influxdb.NewPoint("test", data.Tag, data.Data, time.Now())
		if err != nil {
			// log.Fatal(err)
			fmt.Println("inluxdb error : ", err)
			return
		}
		bp.AddPoint(pt)
		lock.Unlock()
	}
	fmt.Println("over")
}
