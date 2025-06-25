package api

import (
	"net/http"
	"scheduler/params/req"
	"scheduler/params/resp"
	"scheduler/service"

	"github.com/gin-gonic/gin"
)

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
}
