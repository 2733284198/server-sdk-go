/**
 * 融云 Server API go 客户端
 * create by RongCloud
 * create datetime : 2018-11-28
 *
 * v3.0.0
 */

package RCServerSDK

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/astaxie/beego/httplib"
)

const (
	// RONGCLOUDSMSURI 公有云 SMS API 地址 私有云请手动修改
	RONGCLOUDSMSURI = "http://172.29.202.3:18082"
	// RONGCLOUDURI 公有云 API 地址 私有云请手动修改
	RONGCLOUDURI = "http://172.29.202.3:18081"
	// UTF8 字符编码
	UTF8 = "UTF-8"
	// ReqType body类型
	ReqType = "json"
	// USER_AGENT sdk 名称
	USER_AGENT = "rc-go-sdk/3.0"
	// DEFAULTTIMEOUT 超时时间
	DEFAULTTIMEOUT = 30
)

// RongCloud ak sk
type RongCloud struct {
	appKey    string
	appSecret string
	*RongCloudExtra
}

type RongCloudExtra struct {
	RongCloudURI    string
	RongCloudSMSURI string
	TimeOut         time.Duration
}

type CodeReslut struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"errorMessage"`
}

//本地生成签名
//Signature (数据签名)计算方法：将系统分配的 App Secret、Nonce (随机数)、
//Timestamp (时间戳)三个字符串按先后顺序拼接成一个字符串并进行 SHA1 哈希计算。如果调用的数据签名验证失败，接口调用会返回 HTTP 状态码 401。
// getSignature
func (rc *RongCloud) getSignature() (nonce, timestamp, signature string) {
	nonceInt := rand.Int()
	nonce = strconv.Itoa(nonceInt)
	timeInt64 := time.Now().Unix()
	timestamp = strconv.FormatInt(timeInt64, 10)
	h := sha1.New()
	io.WriteString(h, rc.appSecret+nonce+timestamp)
	signature = fmt.Sprintf("%x", h.Sum(nil))
	return
}

// FillHeader 在http header 增加API签名
func (rc *RongCloud) FillHeader(req *httplib.BeegoHTTPRequest) {
	nonce, timestamp, signature := rc.getSignature()
	req.Header("App-Key", rc.appKey)
	req.Header("Nonce", nonce)
	req.Header("Timestamp", timestamp)
	req.Header("Signature", signature)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Header("Rc-Sdk", "server-sdk-go")
	req.Header("User-Agent", USER_AGENT)
}

// fillJSONHeader josn格式
func FillJSONHeader(req *httplib.BeegoHTTPRequest) {
	req.Header("Content-Type", "application/json")
}

func NewRongCloud(appKey, appSecret string, extra *RongCloudExtra) *RongCloud {

	defaultextra := RongCloudExtra{
		RongCloudURI:    RONGCLOUDURI,
		RongCloudSMSURI: RONGCLOUDSMSURI,
		TimeOut:         DEFAULTTIMEOUT,
	}
	if extra == nil {
		rc := RongCloud{
			appKey:         appKey,    //appkey
			appSecret:      appSecret, //appsecret
			RongCloudExtra: &defaultextra,
		}
		return &rc
	} else {
		if extra.TimeOut == 0 {
			extra.TimeOut = DEFAULTTIMEOUT
		}
		if extra.RongCloudSMSURI == "" || extra.RongCloudURI == "" {
			extra.RongCloudURI = RONGCLOUDURI
			extra.RongCloudSMSURI = RONGCLOUDSMSURI
		}
		rc := RongCloud{
			appKey:         appKey,    //appkey
			appSecret:      appSecret, //appsecret
			RongCloudExtra: extra,
		}
		return &rc
	}
}
