package api

import (
	"net/http"
	"scheduler/db/model"
	"scheduler/job"
	"strconv"

	"scheduler/params/req"
	"scheduler/params/resp"
	"scheduler/service"

	"github.com/gin-gonic/gin"
)

// ListTask 用于获取任务列表。
// @param c gin.Context 上下文，包含查询参数。
// @return 返回分页的任务列表和总数，出错时返回错误信息。
func ListTask(c *gin.Context) {
	var query req.TaskQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}
	err := query.ConversionAndVerify()
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}
	l := service.GetTaskSrv()
	list, total, err := l.ListTask(c, &query)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, resp.NewPageDTO().BuildWithTask(list, total))
}

// AddTask 用于新增任务。
// @param c gin.Context 上下文，JSON 请求体为任务参数。
// @return 返回新建任务的ID，或错误信息。
func AddTask(c *gin.Context) {
	var cmd req.TaskCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}
	err := cmd.ConversionAndVerifyWithAdd()
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}
	err = job.VerifySpec(cmd.Spec)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}
	l := service.GetTaskSrv()
	id, err := l.AddTask(c, *model.NewTask().Build(0, cmd))
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &resp.CreatedDTO{Id: id})
}

// GetTask 用于获取单个任务详情。
// @param c gin.Context 上下文，路径参数为任务ID。
// @return 返回任务详情，或错误信息。
func GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: "invalid id parameter",
		})
		return
	}
	l := service.GetTaskSrv()
	task, err := l.GetTask(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, resp.NewTaskDTO().Build(task))
}

// EditTask 用于编辑任务。
// @param c gin.Context 上下文，路径参数为任务ID，JSON 请求体为任务参数。
// @return 无返回体，出错时返回错误信息。
func EditTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: "invalid id parameter",
		})
		return
	}

	var cmd req.TaskCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}

	err = cmd.ConversionAndVerifyWithEdit()
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}
	if len(cmd.Spec) > 0 {
		err = job.VerifySpec(cmd.Spec)
		if err != nil {
			c.JSON(http.StatusBadRequest, resp.CommonDTO{
				Msg: err.Error(),
			})
			return
		}
	}
	l := service.GetTaskSrv()
	err = l.EditTask(c, *model.NewTask().Build(id, cmd))
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteTask 用于删除任务。
// @param c gin.Context 上下文，路径参数为任务ID。
// @return 无返回体，出错时返回错误信息。
func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: "invalid id parameter",
		})
		return
	}
	l := service.GetTaskSrv()
	err = l.DeleteTask(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ExecuteTask 用于立即执行任务。
// @param c gin.Context 上下文，路径参数为任务ID。
// @return 返回任务执行启动信息，或错误信息。
func ExecuteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: "invalid id parameter",
		})
		return
	}
	l := service.GetTaskSrv()
	task, err := l.GetTask(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}

	go func() {
		// 使用主url发起请求
		if job.Execute(c, task, task.Url, 0) {
			return
		}
		// 如果主url请求失败，且有备用url，使用备用url发起请求
		if len(task.BackupUrl) > 0 {
			job.Execute(c, task, task.BackupUrl, 1)
		}
	}()

	c.JSON(http.StatusOK, resp.CommonDTO{Msg: "task execution started"})
}
