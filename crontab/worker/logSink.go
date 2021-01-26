package worker

import (
	"context"
	"crontabInit/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//mongodb存储日志
type LogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var (
	GLogSink *LogSink
)

func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)

	//建立连接
	clientOptions := options.Client().ApplyURI(GConfig.MongodbUri)
	mongoTime := time.Duration(GConfig.MongodbConnectTimeout)
	clientOptions.ConnectTimeout = &mongoTime
	if client, err = mongo.Connect(
		context.TODO(),
		clientOptions); err != nil {
		return
	}

	//赋值
	GLogSink = &LogSink{
		client:         client,
		logCollection:  client.Database("cron").Collection("log"),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}
	go GLogSink.writeLoop()
	return
}

//批量写入日志
func (l *LogSink) saveLogs(batch *common.LogBatch) {
	l.logCollection.InsertMany(context.TODO(), batch.Logs)
}

//消费日志的协程
func (l *LogSink) writeLoop() {
	var (
		log          *common.JobLog
		logBatch     *common.LogBatch
		commitTimer  *time.Timer
		timeoutBatch *common.LogBatch //超时的日志
	)

	for {
		select {
		case log = <-l.logChan:
			//取出日志写入mongodb
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				//批次自动超时提交
				commitTimer = time.AfterFunc(time.Duration(GConfig.JobLogCommitTimeout)*time.Millisecond, func(batch *common.LogBatch) func() {
					return func() {
						l.autoCommitChan <- batch
					}
				}(logBatch),
				)
			}
			//把日志追加到批次中
			logBatch.Logs = append(logBatch.Logs, log)
			//如果批次满100条就发送
			if len(logBatch.Logs) >= GConfig.JobLogBatchSize {
				l.saveLogs(logBatch)
			}
			//清空
			logBatch = nil
			//取消定时器
			commitTimer.Stop()
		case timeoutBatch = <-l.autoCommitChan: //过期的批次写入
			if timeoutBatch != logBatch { //跳过已经被提交的批次
				continue
			}
			l.saveLogs(timeoutBatch)
			//清空logBatch
			logBatch = nil
		}
	}
}

//发送日志
func (l *LogSink) Append(jobLog *common.JobLog) {
	if jobLog == nil {
		return
	}
	select {
	case l.logChan <- jobLog:
	default:
		//满了就丢弃
	}
}
