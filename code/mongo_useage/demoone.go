package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	//"time"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"time"
)

//任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// 一条日志
type LogRecord struct {
	JobName   string    `bson:"jobName"`   // 任务名
	Command   string    `bson:"command"`   // shell命令
	Err       string    `bson:"err"`       // 脚本错误
	Content   string    `bson:"content"`   // 脚本输出
	TimePoint TimePoint `bson:"timePoint"` // 执行时间点
}

func main() {
	var (
		client *mongo.Client
		result *mongo.InsertOneResult
		err    error
	)

	//1.建立连接
	// 建立mongodb连接
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	if client, err = mongo.Connect(
		context.TODO(), clientOptions); err != nil {
		fmt.Println(err)
		return
	}
	// 2, 选择数据库my_db
	database := client.Database("ichunt")

	// 3, 选择表my_collection
	collection := database.Collection("cron_log")

	//4, 插入记录(bson)
	record := &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	if result, err = collection.InsertOne(context.TODO(), record); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)

	//// _id: 默认生成一个全局唯一ID, ObjectID：12字节的二进制
	docId := result.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID:", docId.Hex())
	fmt.Println(client)
}
