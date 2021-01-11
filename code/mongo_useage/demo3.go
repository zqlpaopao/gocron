package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//任务的执行时间点
type TimePointss struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// 一条日志
type LogRecordss struct {
	JobName   string      `bson:"jobName"`   // 任务名
	Command   string      `bson:"command"`   // shell命令
	Err       string      `bson:"err"`       // 脚本错误
	Content   string      `bson:"content"`   // 脚本输出
	TimePoint TimePointss `bson:"timePoint"` // 执行时间点
}

type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client *mongo.Client
		//result *mongo.InsertOneResult
		err       error
		deleteRes *mongo.DeleteResult
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

	//删除开始时间遭遇当前时间的所有日志
	DeleteConds := &DeleteCond{BeforeCond: TimeBeforeCond{
		Before: time.Now().Unix(),
	}}

	deleteRes, err = collection.DeleteMany(context.TODO(), DeleteConds)

	fmt.Println(deleteRes.DeletedCount)
}
