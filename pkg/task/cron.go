package task

import (
	"errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

func init() {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Development = false
	Logger, _ = loggerConfig.Build()
}

type cronLogger struct {
}

func (this_ *cronLogger) Info(msg string, keysAndValues ...interface{}) {
	Logger.Info(msg, zap.Any("keysAndValues", keysAndValues))

}
func (this_ *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	Logger.Error(msg, zap.Any("keysAndValues", keysAndValues), zap.Error(err))
}

var (
	taskCron    *cron.Cron
	cronTaskMap map[string]*CronTask
)

func init() {
	taskCron = cron.New(cron.WithSeconds(), cron.WithLogger(&cronLogger{}))
	taskCron.Start()
	cronTaskMap = map[string]*CronTask{}
}

func addTaskCache(task *CronTask) (err error) {
	if task.Key == "" {
		err = errors.New("任务属性Key不能为空")
		return
	}
	if task.Spec == "" {
		err = errors.New("任务属性Spec不能为空")
		return
	}
	if cronTaskMap[task.Key] != nil {
		err = errors.New("任务Key[" + task.Key + "]已存在")
		return
	}
	cronTaskMap[task.Key] = task

	return
}

func removeTaskCache(task *CronTask) {
	delete(cronTaskMap, task.Key)
}

// AddTask 添加定时任务
func AddTask(task *CronTask) (err error) {
	err = addTaskCache(task)
	if err != nil {
		return
	}

	task.cronEntryID, err = taskCron.AddFunc(task.Spec, task.run)
	if err != nil {
		removeTaskCache(task)
		return
	}
	return
}
