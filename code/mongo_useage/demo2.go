package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//任务的执行时间点
type TimePoints struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// 一条日志
type LogRecords struct {
	JobName   string     `bson:"jobName"`   // 任务名
	Command   string     `bson:"command"`   // shell命令
	Err       string     `bson:"err"`       // 脚本错误
	Content   string     `bson:"content"`   // 脚本输出
	TimePoint TimePoints `bson:"timePoint"` // 执行时间点
}

type FindByJobName struct {
	JobName string `bson:"jobName"`
}

func main() {
	var (
		client *mongo.Client
		//result *mongo.InsertOneResult
		err    error
		cursor *mongo.Cursor
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

	//4, 查询记录
	cron := &FindByJobName{JobName: "job10"}
	Skip := int64(0)
	limit := int64(2)
	option := &options.FindOptions{
		Skip:  &Skip,
		Limit: &limit,
	}

	//查询
	if cursor, err = collection.Find(context.TODO(), cron, option); nil != err {
		fmt.Println(err)
	}

	defer cursor.Close(context.TODO())
	//遍历结果集
	for cursor.Next(context.TODO()) {
		log := &LogRecords{}
		if err = cursor.Decode(log); nil != err {
			fmt.Println(err)
		}
		fmt.Println(log)
	}

}
