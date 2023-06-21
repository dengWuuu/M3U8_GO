package tool

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	URL "net/url"
	"strings"
)

const (
	Identify              = "#EXTM3U"
	NestedPrefix          = "#EXT-X-STREAM-INF"
	M3U8Suffix            = "m3u8"
	HTTPPrefix            = "http"
	HTTPContentTypeHeader = "application/vnd.apple.mpegurl"
	FileFolderPrefix      = "m3u8/"
)

func GetM3U8FileContent(url string) ([]string, error) {
	req, err := http.NewRequest("GET", url, nil)
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
func NewM3U8HttpRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	SetHTTPHeader(req)
	return req, nil
}

func SetHTTPHeader(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
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
		return ""
	}
	// 获取完整的域名
	domain := u.Scheme + "://" + u.Host
	return domain
}

// GetFinalURL 获取嵌套在m3u8文件中的URL
func GetFinalURL(content []string, url string) string {
	finalURL := ""
	// 理论上m3u8文件地址会在数组后面从后开始遍历更快返回二级URL字符串
	for i := len(content) - 1; i >= 0; i-- {
		// 字符串带m3u8后缀 说明这个是一个m3u8的链接
		if strings.HasSuffix(content[i], M3U8Suffix) {
			// 前缀带有 http 说明是完整的 url 不是 uri 不需要拼接
			finalURL = JoinURL(url, content[i])
			break
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

// GenerateKey url对应生成的key是相同的 采用sha256
func GenerateKey(url string) string {
	byteURL := []byte(url)
	hash := sha256.New()
	hash.Write(byteURL)
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	// 返回哈希值 + .m3u8的后缀 			加上m3u8/的前缀，在TOS m3u8/ 代表m3u8目录下存储这个文件
	return FileFolderPrefix + hashCode + "." + M3U8Suffix
}

// TranslateM3U8ContentURL 将嵌套的m3u8 url 变成可访问的url
func TranslateM3U8ContentURL(content []string, url string) []string {
	for i, s := range content {
		if !strings.HasPrefix(s, "#") {
			content[i] = JoinURL(url, s)
		} else {
			// 例子：需要替换uri #EXT-X-KEY:METHOD=AES-128,URI="enc.key",IV=0x00000000000000000000000000000000
			if strings.Contains(s, "URI") {
				// 找到 uri 后面双引号内的左右位置
				uriStart := strings.Index(s, "URI=\"") + 5
				uriEnd := strings.Index(s[uriStart:], "\"") + uriStart
				// 构造替换后的值
				oldURI := s[uriStart:uriEnd]
				newKeyURI := JoinURL(url, oldURI)
				// Replace the URI value
				content[i] = s[:uriStart] + newKeyURI + s[uriEnd:]
			}
		}
	}
	return content
}

func JoinURL(url string, uri string) string {
	finalURL := ""
	if strings.HasPrefix(uri, HTTPPrefix) {
		finalURL = uri
	} else if strings.HasPrefix(uri, "/") {
		finalURL = GetM3U8BaseURL(url) + uri
	} else {
		finalURL = GetM3U8IndexURL(url) + uri
	}
	return finalURL
}
