package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"scheduler/db/model"
	"scheduler/service"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var Shutdown = false

// 调度任务
func scheduled(ctx context.Context, spec string) {
	//获取该cron下的所有任务
	l := service.GetTaskSrv()
	taskList, err := l.ListStartedTaskBySpec(ctx, spec)
	if err != nil {
		zap.L().Error("service.ListStartedTaskBySpec failed, err:%v\n", zap.Error(err))
	}
	//如果任务列表长度为0，则删除该工作者
	if len(taskList) == 0 {
		worker, _ := workerFactory.Load(spec)
		worker.(*cron.Cron).Stop()
		workerFactory.Delete(spec)
		worker = nil
		zap.L().Info("a invalid worker is deleted")
	}
	//循环请求
	for _, task := range taskList {
		go func(task model.Task) {
			if Shutdown {
				return
			}
			if time.Now().UnixMilli() <= task.TimeLock {
				return
			}
			yes, err := l.TryExecuteTask(ctx, task)
			if err != nil {
				zap.L().Error("service.TryExecuteTask failed, err:%v\n", zap.Error(err))
				return
			}
			if yes == 0 {
				return
			}
			// 使用主url发起请求
			if Execute(ctx, task, task.Url, 0) {
				return
			}
			// 如果主url请求失败，且有备用url，使用备用url发起请求
			if len(task.BackupUrl) > 0 {
				Execute(ctx, task, task.BackupUrl, 1)
			}
		}(task)
	}
}

// Execute 执行任务
func Execute(ctx context.Context, task model.Task, url string, isBackup int32) bool {
	l := service.GetRecordSrv()
	for i := 0; i <= int(task.RetryMax); i++ {
		record := model.Record{
			TaskId:     task.Id,
			IsBackup:   isBackup,
			RetryCount: int32(i),
			ExecutedAt: time.Now(),
		}
		code, timeCost, result := httpSend(task.Method, url, task.Body, task.Header)
		record.Result = result
		record.Code = int32(code)
		record.TimeCost = int32(timeCost)
		err := l.AddRecord(ctx, record)
		if err != nil {
			zap.L().Error("service.AddRecord failed, err:%v\n", zap.Error(err))
		}
		if record.Code >= 200 && record.Code < 300 {
			return true
		}
		if task.RetryCycle > 0 {
			time.Sleep(time.Duration(task.RetryCycle) * time.Millisecond)
		}
	}
	return false
}

// 发送http请求
func httpSend(method, url, body, header string) (int, int64, string) {
	if method != "POST" && method != "GET" && method != "PUT" && method != "PATCH" && method != "DELETE" {
		return -1, 0, "http method is not match"
	}
	if len(url) == 0 {
		return -1, 0, "http url is empty"
	}
	payload := strings.NewReader(body)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return -1, 0, fmt.Sprintf("http build request error: %s", err.Error())
	}
	if len(header) > 2 {
		var headerMap map[string]string
		err = json.Unmarshal([]byte(header), &headerMap)
		if err != nil {
			return -1, 0, fmt.Sprintf("http header error: %s", err.Error())
		}
		if len(headerMap) > 0 {
			for k, v := range headerMap {
				req.Header.Add(k, v)
			}
		}
	}
	startTime := time.Now().UnixMilli()
	response, err := http.DefaultClient.Do(req)
	endTime := time.Now().UnixMilli()
	timeCost := endTime - startTime
	if err != nil {
		return -1, timeCost, fmt.Sprintf("http send error: %s", err.Error())
	}
	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			zap.L().Error("body.Close failed, err:%v\n", zap.Error(err))
		}
	}(response.Body)
	resultBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, timeCost, fmt.Sprintf("http build response error: %s", err.Error())
	}
	return response.StatusCode, timeCost, string(resultBytes)
}
