package tool

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	URL "net/url"
	"path"
	"strings"
)

const Identify = "#EXTM3U"
const NestedPrefix = "#EXT-X-STREAM-INF"
const M3U8Suffix = "m3u8"
const HTTPPrefix = "http"

func GetM3U8FileContent(url string) ([]string, error) {

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("x-use-ppe", "1")
	req.Header.Set("x-tt-env", "ppe_13156482")
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
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

func GetM3U8BaseURL(url string) string {
	// 解析 URL
	u, err := URL.Parse(url)
	if err != nil {
		panic(err)
	}
	// 获取完整的域名
	domain := u.Scheme + "://" + u.Host
	return domain
}

func GetM3U8Filename(url string) string {
	u, err := URL.Parse(url)
	if err != nil {
		panic(err)
	}
	filename := path.Base(u.Path)
	return filename
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

func ConvertStringSlice2ByteSlice(content []string) []byte {
	// 将字符串切片转换为一个字符串，每个元素以换行符分隔
	str := ""
	for _, s := range content {
		str += s + "\n"
	}
	return []byte(str)
}

func WriteToFile(content []string) {
	// 将字符串写入文件
	err := ioutil.WriteFile("output.txt", ConvertStringSlice2ByteSlice(content), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File written successfully")
}

func GenerateKey(url string) string {
	byteURL := []byte(url)
	hash := sha256.New()
	//输入数据
	hash.Write(byteURL)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	return hashCode
}
