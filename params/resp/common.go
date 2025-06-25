package resp

// CreatedDTO 用于返回创建成功后的 ID
type CreatedDTO struct {
	Id int64 `json:"id"`
}

// CommonDTO 用于通用消息返回
type CommonDTO struct {
	Msg string `json:"msg"`
}
