package api

import (
	"net/http"
	"scheduler/params/resp"
	"scheduler/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListTag 用于获取标签统计信息。
// @param c gin.Context 上下文，可选 status 查询参数。
// @return 返回标签统计列表或错误信息。
func ListTag(c *gin.Context) {
	statusStr := c.Query("status")
	var status int
	if statusStr != "" {
		var err error
		status, err = strconv.Atoi(statusStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, resp.CommonDTO{
				Msg: "invalid status parameter",
			})
			return
		}
	}
	l := service.GetCountSrv()
	list, err := l.ListTagCount(c, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, list)
}

// ListSpec 用于获取 Spec 统计信息。
// @param c gin.Context 上下文，可选 status 查询参数。
// @return 返回 Spec 统计列表或错误信息。
func ListSpec(c *gin.Context) {
	statusStr := c.Query("status")
	var status int
	if statusStr != "" {
		var err error
		status, err = strconv.Atoi(statusStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, resp.CommonDTO{
				Msg: "invalid status parameter",
			})
			return
		}
	}
	l := service.GetCountSrv()
	list, err := l.ListSpecCount(c, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, list)
}
