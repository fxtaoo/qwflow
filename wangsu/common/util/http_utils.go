package util

import (
	"io"
	"net/http"
	"qwflow/wangsu/common/model"
	"strings"
)

func Call(requestMsg model.HttpRequestMsg) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest(requestMsg.Method, requestMsg.Url, strings.NewReader(requestMsg.Body))

	for key := range requestMsg.Headers {
		req.Header.Set(key, requestMsg.Headers[key])
	}
	resp, _ := client.Do(req)

	body, _ := io.ReadAll(resp.Body)

	// fmt.Println(resp)
	return body
}
