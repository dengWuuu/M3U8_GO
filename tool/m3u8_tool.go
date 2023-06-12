package tool

import (
	"fmt"
	"io"
	"net/http"
)

func GetM3U8File(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("io.ReadCloser 关闭失败")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("获取m3u8文件内容失败")
	}
	return body
}
