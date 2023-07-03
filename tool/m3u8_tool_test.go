package tool

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestURLTool(t *testing.T) {
	url := "https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/index.m3u8"
	fmt.Println(strings.EqualFold(getM3U8IndexURL(url), "https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/"))
}

func TestJson(t *testing.T) {
	params := "{\"tea_event_index\":2164,\"app_start_times_by_version\":14,\"ttweb_media_channel\":\"browser_activity\",\"ttweb_media_time_to_play_ready\":0,\"nt\":4,\"webview_id\":56,\"pid\":23996,\"ttweb_media_source\":\"https://110.42.2.247:9092/c/iqiyi_301/772c6a04e81bce39ab55cec4d65c73f7.m3u8?vkey=5b4aVZkk0anbeNJP33k0n1pfdNc6myvxR5G-XGXahZxB\",\"app_host_abi\":\"64\",\"ttweb_media_frame_host\":\"www.wuaitao.net\",\"player_id\":4,\"os_api\":30,\"ttweb_media_has_ever_played\":\"true\",\"loadso\":\"0881130050506\",\"sdk_update_version_code\":\"119205\",\"webview_count\":0,\"ttweb_media_is_movie\":\"true\",\"ttweb_media_player_type\":\"MediaPlayer_ttmp\",\"ttweb_media_time_to_play_event\":14095,\"sdk_scc_version\":20,\"ttweb_media_frame_url\":\"http://www.wuaitao.net/vod-play-id-3817-src-1-num-1.html\",\"kernel_scc_version\":20,\"ttweb_media_time_to_first_frame\":0,\"os_version\":\"11\",\"processname\":\"com.cat.readall\",\"ttweb_media_frame_title\":\"《伙头军客栈》-第01集-在线观看-百度云-吾爱淘电影\",\"ttweb_media_duration\":0,\"ttweb_media_is_video\":\"true\",\"ttweb_media_time_to_metadata\":361,\"app_start_times\":14,\"ttweb_media_time_to_played\":0.101224,\"ttweb_media_stay_time\":14456,\"ttweb_media_has_reached_have_enough\":\"false\",\"ttweb_media_video_width\":0,\"sdk_aar_version\":\"0621130048\",\"ttweb_media_error_code\":-1094995529,\"ttweb_media_extra_data\":\"{\\\"search_position\\\":\\\"top_bar\\\",\\\"position\\\":\\\"search\\\",\\\"query\\\":\\\"伙头军客栈 电视剧\\\",\\\"show_rank\\\":22,\\\"enter_from\\\":\\\"click_search\\\",\\\"origin_url\\\":\\\"http:\\\\/\\\\/www.wuaitao.net\\\\/vod-detail-id-3817.html\\\",\\\"ad_id\\\":0,\\\"gd_enter_from\\\":\\\"click_search\\\",\\\"category_name\\\":\\\"__search__\\\",\\\"search_id\\\":\\\"20230618191415BC90E76A8924E00FC679\\\",\\\"search_result_id\\\":\\\"-6619194857463894961\\\",\\\"req_id\\\":\\\"202306181915019A5B61417AD303B195AF\\\",\\\"enter_group_id\\\":\\\"\\\",\\\"query_id\\\":\\\"6668994262828324110\\\",\\\"search_subtab_name\\\":\\\"synthesis\\\",\\\"log_pb\\\":{\\\"impr_id\\\":\\\"202306181915019A5B61417AD303B195AF\\\",\\\"is_incognito\\\":0},\\\"cell_type\\\":67,\\\"result_type\\\":\\\"toutiao_web\\\",\\\"rank\\\":22,\\\"source\\\":\\\"album_tag\\\",\\\"from\\\":\\\"album_tag\\\",\\\"offset\\\":30,\\\"db_name\\\":\\\"L1\\\",\\\"page_type\\\":\\\"1\\\",\\\"query_type\\\":\\\"SearchAggregationQueryType\\\",\\\"search_parent_from\\\":\\\"album_tag\\\",\\\"enter_time_mills\\\":1687086957182,\\\"origin_title\\\":\\\"《伙头军客栈》全集高清视频-在线观看-百度云-吾爱淘电影\\\"}\",\"ttweb_media_video_height\":0}"
	paramsMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(params), &paramsMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(params)
	fmt.Println(paramsMap)
}

func TestIsM3U8URL(t *testing.T) {
	url := "https://hnzy.bfvvs.com/play/pen8nwYb/index.m3u8"
	fmt.Println(IsM3U8URL(url))
}

func TestKey(t *testing.T) {
	urls := []string{
		"https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/index.m3u8",
		"https://hnzy.bfvvs.com/play/Ddw29WRb/index.m3u8?query=aksjdhfklashdfk",
		"https://hnzy.bfvvs.com/play/pen8nwYb/index.m3u8",
		"https://hnzy.bfvvs.com/play/mep7royd/index.m3u8",
		"https://hnzy.bfvvs.com/play/nelDK47e/index.m3u8",
		"https://hnzy.bfvvs.com/play/DdwZPPwb/index.m3u8",
		"https://hnzy.bfvvs.com/play/DbDPBBnb/index.m3u8",
		"https://hnzy.bfvvs.com/play/RdGzKV0b/index.m3u8",
		"https://hnzy.bfvvs.com/play/RdGzKV0b/index.m3u8",
		"https://hnzy.bfvvs.com/play/YaOLoLRb/index.m3u8",
		"https://hnzy.bfvvs.com/play/rb22vGjb/index.m3u8",
	}
	for _, url := range urls {
		fmt.Println(GenerateKey(url, EmptyString, true))
	}
}

func TestURLBase(t *testing.T) {
	url := "https://vip.lz-cdn1.com/20230613/21106_a7408872/index.m3u8"
	fmt.Println(getM3U8BaseURL(url))
}

func TestGetSimpleM3U8(t *testing.T) {
	url := "https://hnzy.bfvvs.com/play/Ddw29WRb/index.m3u8"
	fmt.Printf("读取的URL: %v \n", url)
	content, err := GetM3U8FileContent(nil, url)

	if err != nil {
		fmt.Printf("获取文件内容失败 %v", err)
		return
	}

	if !IsM3U8(content[0]) {
		fmt.Println("不是m3u8文件")
		return
	}
	if !IsNested(content) {
		content = translateM3U8ContentURL(content, url)
		fmt.Println("不是嵌套m3u8文件", content)
		return
	}

	if !IsSimpleSourceM3U8(content) {
		content = translateM3U8ContentURL(content, url)
		fmt.Println("多源的m3u8文件")
		return
	}
	finalURL := GetFinalURL(content, url)
	if strings.EqualFold(finalURL, "") {
		fmt.Println("未获取到嵌套二级URL")
	}
	fmt.Println(finalURL)

	content, err = GetM3U8FileContent(nil, finalURL)
	if err != nil {
		fmt.Printf("获取嵌套文件内容失败 %v", err)
		return
	}
	byteSlice := ReturnM3U8Content(content, finalURL)

	fmt.Println(string(byteSlice))
}

func TestTranslate(t *testing.T) {
	content := []string{"#EXT-X-KEY:METHOD=AES-128,URI=\"enc.key\",IV=0x00000000000000000000000000000000"}
	translateM3U8ContentURL(content, "https://example.com/index.m3u8")
}

func TestGetTSURLFromM3U8(t *testing.T) {
	content := []string{"#EXTM3U",
		"#EXT-X-VERSION:3,",
		"#EXT-X-PLAYLIST-TYPE:VOD,",
		"#EXT-X-MEDIA-SEQUENCE:0,",
		"#EXT-X-TARGETDURATION:9,",
		"#EXT-X-DISCONTINUITY,",
		"#EXTINF:6.083333,",
		"https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/2000k/hls/56207c2dfad000000.ts",
		"#EXTINF:4.125000,",
		"https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/2000k/hls/56207c2dfad000001.ts",
		"#EXTINF:4.125000,",
		"https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/2000k/hls/56207c2dfad000002.ts",
		"#EXTINF:4.166667,",
		"https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/2000k/hls/56207c2dfad000003.ts",
		"#EXTINF:3.333333,",
		"https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/2000k/hls/56207c2dfad000004.ts",
		"#EXTINF:4.166667,",
		"https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/2000k/hls/56207c2dfad000005.ts",
		"#EXTINF:2.583333,",
		"https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/2000k/hls/56207c2dfad000006.ts",
		"#EXTINF:4.166667,",
		"https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/2000k/hls/56207c2dfad000007.ts",
		"#EXTINF:5.416667,",
		"https://hd.lz-cdn18.com/20230610/3269_5e4f3eae/2000k/hls/56207c2dfad000008.ts",
	}
	ls, isOverNum := GetTSURLFromM3U8(60, content)
	assert.Equal(t, len(ls), 9)
	assert.Equal(t, isOverNum, true)
}

func TestGetAllFinalM3U8URL(t *testing.T) {
	url := ""
	content := []string{
		"#EXTM3U",
		"#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=200000",
		"https://devimages.apple.com/iphone/samples/bipbop/gear1/prog_index.m3u8",
		"#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=311111",
		"https://devimages.apple.com/iphone/samples/bipbop/gear2/prog_index.m3u8",
		"#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=484444",
		"https://devimages.apple.com/iphone/samples/bipbop/gear3/prog_index.m3u8",
		"#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=737777",
		"https://devimages.apple.com/iphone/samples/bipbop/gear4/prog_index.m3u8",
	}
	allFinalM3U8URL := GetAllFinalM3U8URL(content, url)
	assert.Equal(t, len(allFinalM3U8URL), 4)
}
