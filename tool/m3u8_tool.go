package tool

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const Identify = "#EXTM3U"
const NestedPrefix = "#EXT-X-STREAM-INF"
const M3U8Suffix = "m3u8"
const HTTPPrefix = "http"

func GetM3U8FileContent(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("获取m3u8文件内容失败")
		return nil, err
	}

	lines := strings.Split(string(body), "\n")
	return lines, nil
}

func IsM3U8(identify string) bool {
	return identify == Identify
}
func IsNested(content []string) bool {
	for _, str := range content {
		if strings.HasPrefix(str, NestedPrefix) {
			return true
		}
	}
	return false
}

func IsSimpleSourceM3U8(content []string) bool {
	cnt := 0
	for _, str := range content {
		if strings.HasSuffix(str, M3U8Suffix) {
			cnt++
		}
	}
	return cnt <= 1
}

func GetM3U8IndexURL(url string) string {
	// 使用 strings 包中的 LastIndex 函数查找最后一个 / 的位置
	lastIndex := strings.LastIndex(url, "/")
	// 使用切片操作获取前面部分的 URL
	indexURL := url[:lastIndex+1]
	return indexURL
}

func GetFinalURL(content []string, url string) string {
	finalURL := ""

	// 理论上m3u8文件地址会在数组后面从后开始遍历更快返回二级URL字符串
	for i := len(content) - 1; i >= 0; i-- {
		if strings.HasSuffix(content[i], M3U8Suffix) {
			// 如果前缀带有 http 说明是完整的 url 不是 uri 不需要拼接
			if strings.HasPrefix(content[i], HTTPPrefix) {
				finalURL = content[i]
				break
			} else {
				finalURL = GetM3U8IndexURL(url) + content[i]
				break
			}
		}
	}
	return finalURL
}

func WriteToFile(content []string) {

	// 将字符串切片转换为一个字符串，每个元素以换行符分隔
	str := ""
	for _, s := range content {
		str += s + "\n"
	}

	// 将字符串写入文件
	err := ioutil.WriteFile("output.txt", []byte(str), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File written successfully")
}
