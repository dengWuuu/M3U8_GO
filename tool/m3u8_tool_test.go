package tool

import (
	"fmt"
	"testing"
)

func TestGetM3U8(t *testing.T) {
	url := "https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/index.m3u8"
	fileByte := GetM3U8File(url)

	fmt.Println(string(fileByte))
}
