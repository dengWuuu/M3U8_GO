package tool

import (
	"fmt"
	"strings"
	"testing"
)

func TestURLTool(t *testing.T) {
	url := "https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/index.m3u8"
	fmt.Println(strings.EqualFold(GetM3U8IndexURL(url), "https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/"))
}

func TestGetM3U8(t *testing.T) {
	url := "https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/index.m3u8"
	fmt.Printf("读取的URL: %v \n", url)

	content, err := GetM3U8FileContent(url)

	if err != nil {
		fmt.Printf("获取文件内容失败 %v", err)
		return
	}

	if !IsM3U8(content[0]) {
		fmt.Println("不是m3u8文件")
		return
	}
	fmt.Println(content)
	if !IsNested(content) {
		fmt.Println("不是嵌套m3u8文件")
		return
	}

	if !IsSimpleSourceM3U8(content) {
		fmt.Println("多源的m3u8文件")
		return
	}
	finalURL := GetFinalURL(content, url)
	if strings.EqualFold(finalURL, "") {
		fmt.Println("未获取到嵌套二级URL")
	}
	fmt.Println(finalURL)

	content, err = GetM3U8FileContent(finalURL)
	if err != nil {
		fmt.Printf("获取嵌套文件内容失败 %v", err)
		return
	}

	fmt.Println(content)

	WriteToFile(content)

}
