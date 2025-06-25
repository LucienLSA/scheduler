package api

import (
	"fmt"
	"net/http"
	"os"
	"scheduler/job"
	"scheduler/params/resp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	if job.Shutdown {
		c.JSON(http.StatusServiceUnavailable, resp.CommonDTO{
			Msg: "closed",
		})
		return
	}

	c.JSON(http.StatusOK, resp.CommonDTO{Msg: "running"})
}

func Shutdown(c *gin.Context) {
	waitStr := c.Query("wait")
	if waitStr == "" {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: "wait parameter is required",
		})
		return
	}

	wait, err := strconv.Atoi(waitStr)
	if err != nil || wait <= 0 {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: "wait must be greater than 0",
		})
		return
	}

	if job.Shutdown {
		c.JSON(http.StatusBadRequest, resp.CommonDTO{
			Msg: "shutdown request has been triggered",
		})
		return
	}

	host := c.Request.Host
	if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") || strings.Contains(host, "0.0.0.0") {
		job.Shutdown = true
		go func() {
			time.Sleep(time.Duration(wait) * time.Second)
			os.Exit(1)
		}()
		c.JSON(http.StatusOK, resp.CommonDTO{Msg: fmt.Sprintf("shutdown after %v seconds", wait)})
		return
	}

	c.JSON(http.StatusForbidden, resp.CommonDTO{
		Msg: "only local shutdown request are accepted",
	})
}
