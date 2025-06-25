package tools

import (
	"errors"
	"net/url"
)

// verifyPageQueryParams 校验分页参数
func VerifyPageQueryParams(pageIndex int, pageSize int) error {
	if pageIndex <= 0 {
		return errors.New("page index must be greater than 0")
	}
	if pageSize < 0 {
		return errors.New("page size must be greater than or equal to 0")
	}
	return nil
}

// verifyUrl 校验 URL 格式
func VerifyUrl(sourceUrl string) bool {
	if len(sourceUrl) < 6 {
		return false
	}
	_, err := url.ParseRequestURI(sourceUrl)
	if err != nil {
		return false
	}
	u, err := url.Parse(sourceUrl)
	if err != nil || len(u.Scheme) == 0 || len(u.Host) == 0 {
		return false
	}
	// Check if the URL has a valid scheme (http or https)
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return true
}
