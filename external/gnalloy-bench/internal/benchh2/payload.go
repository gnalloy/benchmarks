package benchh2

// ResponseBody 返回固定响应体，确保所有框架返回同样字节数。
func ResponseBody(size int) []byte {
	body := make([]byte, size)
	for i := range body {
		body[i] = byte(i)
	}
	return body
}
