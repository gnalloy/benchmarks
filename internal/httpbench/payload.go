package httpbench

// ResponseBody 返回各服务端共用的确定性响应体。
func ResponseBody(size int) []byte {
	body := make([]byte, size)
	for index := range body {
		body[index] = byte(index)
	}
	return body
}
