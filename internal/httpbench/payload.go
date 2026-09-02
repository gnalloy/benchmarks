package httpbench

import (
	"bytes"
	"strconv"
	"strings"
)

var headerTerminator = []byte("\r\n\r\n")

// ResponseBody 返回各服务端共用的确定性响应体。
func ResponseBody(size int) []byte {
	body := make([]byte, size)
	for index := range body {
		body[index] = byte(index)
	}
	return body
}

// ResponseBytes 返回 HTTP/1.1 keep-alive 固定响应帧。
func ResponseBytes(payload int) []byte {
	body := ResponseBody(payload)
	header := "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " +
		strconv.Itoa(len(body)) + "\r\nConnection: keep-alive\r\n\r\n"
	out := make([]byte, 0, len(header)+len(body))
	out = append(out, header...)
	out = append(out, body...)
	return out
}

// RequestBytes 返回固定 GET 请求，payload 仅描述响应体大小。
func RequestBytes(host string) []byte {
	host = strings.TrimSpace(host)
	if host == "" {
		host = "127.0.0.1"
	}
	return []byte("GET /bench HTTP/1.1\r\nHost: " + host + "\r\nUser-Agent: gnalloy-bench\r\nAccept: */*\r\nConnection: keep-alive\r\n\r\n")
}

// ServerState 保存服务端连接上的半包 HTTP 请求头。
type ServerState struct {
	pending []byte
}

// AppendAndCountRequests 追加新字节并返回完整 HTTP/1 请求数。
func (s *ServerState) AppendAndCountRequests(data []byte) int {
	if len(data) > 0 {
		s.pending = append(s.pending, data...)
	}
	count := 0
	for {
		index := bytes.Index(s.pending, headerTerminator)
		if index < 0 {
			return count
		}
		count++
		consumed := index + len(headerTerminator)
		copy(s.pending, s.pending[consumed:])
		s.pending = s.pending[:len(s.pending)-consumed]
	}
}
