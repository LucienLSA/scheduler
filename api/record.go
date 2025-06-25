package api

import (
	"net/http"
	"scheduler/params/req"
	"scheduler/params/resp"
	"scheduler/service"

	"github.com/gin-gonic/gin"
)

// ListRecord 用于获取任务执行记录列表。
// @param c gin.Context 上下文，包含查询参数：shard、taskId、code、startTime、endTime、pageIndex、pageSize。
// @return 返回分页的记录列表和总数，出错时返回错误信息。
func ListRecord(c *gin.Context) {
	var query req.RecordQuery
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
	l := service.GetRecordSrv()
	list, total, err := l.ListRecord(c, query.Shard, query.TaskId, query.Code, query.StartTime, query.EndTime, query.PageIndex, query.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp.NewPageDTO().BuildWithRecord(list, total))
}
