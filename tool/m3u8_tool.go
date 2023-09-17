package tool

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	URL "net/url"
	"regexp"
	"strings"
)

const (
	Identify                  = "#EXTM3U"
	NestedPrefix              = "#EXT-X-STREAM-INF"
	M3U8Suffix                = ".m3u8"
	TSSuffix                  = ".ts"
	HTTPPrefix                = "http"
	M3U8HTTPContentTypeHeader = "application/vnd.apple.mpegurl"
	TSHTTPContentTypeHeader   = "application/x-linguist"
	M3U8FileFolderPrefix      = "m3u8/"
	TSFileFolderPrefix        = "ts/"
	EmptyString               = ""
)

func GetFileContent(ctx context.Context, url string, header http.Header) ([]byte, error) {
	// 构建http请求 模仿浏览器发出请求设置请求头
	req, err := newHttpRequest(url, header)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Printf("[services.GetFileContent] resp:%v \n", resp)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetM3U8FileContent(ctx context.Context, url string) ([]string, error) {
	body, err := GetFileContent(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(body), "\n")
	return lines, nil
}

func newHttpRequest(url string, header http.Header) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	setHTTPHeader(req, header)
	return req, nil
}

func setHTTPHeader(req *http.Request, header http.Header) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	for key, value := range header {
		for _, val := range value {
			req.Header.Add(key, val)
		}
	}
}

func IsM3U8URL(url string) bool {
	parseURL, err := URL.Parse(url)
	if err != nil {
		return false
	}
	return strings.HasSuffix(parseURL.Path, M3U8Suffix)
}

func IsTsURL(url string) bool {
	parseURL, err := URL.Parse(url)
	if err != nil {
		return false
	}
	return strings.HasSuffix(parseURL.Path, TSSuffix)
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
		if strings.HasPrefix(str, NestedPrefix) {
			cnt++
		}
	}
	return cnt <= 1
}

func getM3U8IndexURL(url string) string {
	// 使用 strings 包中的 LastIndex 函数查找最后一个 / 的位置
	lastIndex := strings.LastIndex(url, "/")
	// 使用切片操作获取前面部分的 URL
	indexURL := url[:lastIndex+1]
	return indexURL
}

func getM3U8BaseURL(url string) string {
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
		// 判断是不是m3u8链接
		if strings.Contains(url, M3U8Suffix) {
			// 前缀带有 http 说明是完整的 url 不是 uri 不需要拼接
			finalURL = joinURL(url, content[i])
			if IsM3U8URL(finalURL) {
				return finalURL
			}
		}
	}
	return finalURL
}

func convertStringSlice2ByteSlice(content []string) []byte {
	// 将字符串切片转换为一个字符串，每个元素以换行符分隔
	str := ""
	for _, s := range content {
		str += s + "\n"
	}
	return []byte(str)
}

// GenerateKey url对应生成的key是相同的 采用sha256
// 如果是 m3u8 产生的文件名是 m3u8/m3u8文件名
// 如果是 ts 产生的文件名是 ts/m3u8文件名/ts文件名
// 根据 isM3U8 判断生成的是M3U8 key 还是 ts key
func GenerateKey(originURL string, tsURL string, isM3U8 bool) string {
	var byteURL []byte
	if isM3U8 {
		byteURL = []byte(originURL)
	} else {
		byteURL = []byte(tsURL)
	}
	hash := sha256.New()
	hash.Write(byteURL)
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	if isM3U8 {
		return M3U8FileFolderPrefix + hashCode + M3U8Suffix
	} else {
		return TSFileFolderPrefix + GenerateKey(originURL, EmptyString, true) + "/" + hashCode + TSSuffix
	}
}

// TranslateM3U8ContentURL 将嵌套的m3u8 url 变成可访问的url
func translateM3U8ContentURL(content []string, url string) []string {
	for i, s := range content {
		if s == "" {
			continue
		}
		if !strings.HasPrefix(s, "#") {
			content[i] = joinURL(url, s)
		} else {
			// 例子：需要替换uri #EXT-X-KEY:METHOD=AES-128,URI="enc.key",IV=0x00000000000000000000000000000000
			if strings.Contains(s, "URI") {
				// 找到 uri 后面双引号内的左右位置
				uriStart := strings.Index(s, "URI=\"") + 5
				uriEnd := strings.Index(s[uriStart:], "\"") + uriStart
				// 构造替换后的值
				oldURI := s[uriStart:uriEnd]
				newKeyURI := joinURL(url, oldURI)
				// Replace the URI value
				content[i] = s[:uriStart] + newKeyURI + s[uriEnd:]
			}
		}
	}
	return content
}

func ReturnM3U8Content(content []string, url string) []byte {
	content = translateM3U8ContentURL(content, url)
	return convertStringSlice2ByteSlice(content)
}

func joinURL(url string, uri string) string {
	var finalURL string
	if strings.HasPrefix(uri, HTTPPrefix) {
		finalURL = uri
	} else if strings.HasPrefix(uri, "/") {
		finalURL = getM3U8BaseURL(url) + uri
	} else {
		finalURL = getM3U8IndexURL(url) + uri
	}
	return finalURL
}

// GetTSURLFromM3U8 从 m3u8 文件获取 ts url 返回一个string数组代表 ts url 的集合
// 输入值: num 代表要获取m3u8文件的几个ts url , content 代表 m3u8 文件
// 返回值: bool 代表要存储的 ts 个数是否超过了 m3u8 含有的全部 ts 个数
// cnt记录总数， cur记录存储了多少
func GetTSURLFromM3U8(num int, content []string) (tsURLs []string, overNumber bool) {
	cnt := 0
	cur := 0
	for _, url := range content {
		if !IsTsURL(url) {
			continue
		}
		cnt++
		if cur < num {
			tsURLs = append(tsURLs, url)
			cur++
		}
	}
	result := float32(cur) / float32(cnt)
	return tsURLs, result >= 0.75
}

// GetAllFinalM3U8URL 从一个 m3u8 一级文件获取所有二级 m3u8 文件的url链接
func GetAllFinalM3U8URL(content []string, url string) (finalURLs []string) {
	for _, s := range content {
		if strings.Contains(s, M3U8Suffix) {
			finalURL := joinURL(url, s)
			if IsM3U8URL(finalURL) {
				finalURLs = append(finalURLs, finalURL)
			}
		}
	}
	return finalURLs
}

var linePattern = regexp.MustCompile(`([a-zA-Z-]+)=("[^"]+"|[^",]+)`)

// ParseLineParameters extra parameters in string `line`
func ParseLineParameters(line string) map[string]string {
	r := linePattern.FindAllStringSubmatch(line, -1)
	params := make(map[string]string)
	for _, arr := range r {
		params[arr[1]] = strings.Trim(arr[2], "\"")
	}
	return params
}
