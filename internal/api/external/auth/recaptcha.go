package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

type AliRes struct {
	Detail    string `json:"detail"`
	RiskLevel string `json:"risk_level"`
	RequestId string `json:"request_id"`
	Msg       string `json:"msg"`
	Code      int    `json:"code"`
}

func SendRequest(httpMethod string, path string, query map[string]string, token, sessionID, sig, remoteIP string) (resCode int, err error) {
	aliyunBaseUrl := "https://afs.aliyuncs.com"

	if query == nil {
		query = make(map[string]string)
	}
	canonicalizedString := SignRequest(httpMethod, path, query, token, sessionID, sig, remoteIP)
	uri := aliyunBaseUrl + path + "?" + canonicalizedString

	targetUrl, err := url.Parse(uri)
	if err != nil {
		return 0, err
	}

	request := &http.Request{
		Method:     httpMethod,
		ProtoMajor: 1,
		ProtoMinor: 1,
		URL:        targetUrl,
	}
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	log.Println(string(bodyBytes))

	var aliRes AliRes
	err = json.Unmarshal(bodyBytes, &aliRes)
	if err != nil {
		return 0, nil
	}

	if response.StatusCode != 200 {
		log.Println(response.Status)
		return aliRes.Code, errors.New(response.Status)
	}

	return aliRes.Code, nil
}

func SignRequest(httpMethod string, path string, query map[string]string, token, sessionID, sig, remoteIP string) string {
	accKey := config.C.AliyunRecaptcha.AccessKey
	accKeySecret := config.C.AliyunRecaptcha.AccSecretKey
	appKey := config.C.AliyunRecaptcha.AppKey
	scene := config.C.AliyunRecaptcha.Scene

	// https://help.aliyun.com/document_detail/29745.html?spm=a2c4g.11186623.6.618.51e67ebb2zKj8K
	// set common parameters
	query["Format"] = "json"                                   // 返回值的類型，支持JSON與XML。默認為XML
	query["Version"] = "2018-01-12"                            // API版本號，為日期形式：YYYY-MM-DD，本版本對應為2018-01-12
	query["AccessKeyId"] = accKey                              // 阿里雲頒發給用戶的訪問服務所用的密鑰ID
	query["SignatureMethod"] = "HMAC-SHA1"                     // 簽名方式
	query["Timestamp"] = time.Now().UTC().Format(time.RFC3339) // 請求的時間戳。日期格式按照ISO8601標準表示，並需要使用UTC時間。
	query["SignatureVersion"] = "1.0"                          // 簽名算法版本
	query["SignatureNonce"] = GenerateGUID()                   // 唯一隨機數，用於防止網絡重放攻擊。用戶在不同請求間要使用不同的隨機數值

	query["AppKey"] = appKey
	query["RemoteIp"] = remoteIP
	query["Scene"] = scene

	query["Token"] = token
	query["SessionId"] = sessionID
	query["Sig"] = sig

	// 1.使用請求參數構造規範化的請求字符串（Canonicalized Query String）
	// a) 按照參數名稱的字典順序對請求中所有的請求參數（包括文檔中描述的“公共請求參數”和給定了的請求接口的自定義參數
	// 但不能包括“公共請求參數”中提到Signature參數本身）進行排序。
	// 註意：此排序嚴格大小寫敏感排序。
	keys := []string{}
	for key, _ := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// b) 對每個請求參數的名稱和值進行編碼。名稱和值要使用UTF-8字符集進行URL編碼
	// c) 對編碼後的參數名稱和值使用英文等號（=）進行連接。
	// d) 再把英文等號連接得到的字符串按參數名稱的字典順序依次使用&符號連接，即得到規範化請求字符串。
	var buffer bytes.Buffer
	for _, key := range keys {
		if buffer.Len() > 0 {
			buffer.WriteByte('&')
		}
		buffer.WriteString(percentEncode(key))
		buffer.WriteByte('=')
		buffer.WriteString(percentEncode(query[key]))
	}

	canonicalizedString := buffer.String()

	// 使用上一步構造的規範化字符串按照下面的規則構造用於計算簽名的字符串：
	// StringToSign= HTTPMethod + “&” + percentEncode(“/”) + ”&” + percentEncode(CanonicalizedQueryString)
	// 其中HTTPMethod是提交請求用的HTTP方法，比如GET。
	buffer.Reset()
	buffer.WriteString(httpMethod)
	buffer.WriteByte('&')
	buffer.WriteString(percentEncode(path))
	buffer.WriteByte('&')
	buffer.WriteString(percentEncode(canonicalizedString))

	stringToSign := buffer.Bytes()

	// 按照RFC2104的定義，使用上面的用於簽名的字符串計算簽名HMAC值。
	// 註意：計算簽名時使用的Key就是用戶持有的Access Key Secret並加上一個“&”字符(ASCII:38)，使用的哈希算法是SHA1。
	key := []byte(accKeySecret + "&")
	mac := hmac.New(sha1.New, key)
	mac.Write(stringToSign)
	bytes := mac.Sum(nil)

	// 按照Base64編碼規則把上面的HMAC值編碼成字符串，即得到簽名值（Signature）。
	signature := base64.RawStdEncoding.EncodeToString(bytes)

	// 將得到的簽名值作為Signature參數添加到請求參數中，即完成對請求簽名的過程。
	// 註意：得到的簽名值在作為最後的請求參數值提交給DNS服務器的時候，要和其他參數一樣，按照RFC3986的規則進行URL編碼）
	query["Signature"] = signature

	return canonicalizedString + "&Signature=" + percentEncode(signature) + "%3D"
}

func GenerateGUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.WithError(err).Error("Cannot get random number")
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// RFC3986
func percentEncode(str string) string {
	// URL編碼的編碼規則是：
	// i. 對於字符 A-Z、a-z、0-9以及字符“-”、“_”、“.”、“~”不編碼；
	// ii. 對於其他字符編碼成“%XY”的格式，其中XY是字符對應ASCII碼的16進制表示。比如英文的雙引號（”）對應的編碼就是%22
	// iii. 對於擴展的UTF-8字符，編碼成“%XY%ZA…”的格式；
	// iv. 需要說明的是英文空格（ ）要被編碼是%20，而不是加號（+）。
	// 註：一般支持URL編碼的庫（比如Java中的java.net.URLEncoder）都是按照“application/x-www-form-urlencoded”的MIME類型的規則進行編碼的。
	// 實現時可以直接使用這類方式進行編碼，把編碼後的字符串中加號（+）替換成%20、星號（*）替換成%2A、%7E替換回波浪號（~），即可得到上述規則描述的編碼字符串。
	encoded := url.QueryEscape(str)

	encoded = strings.Replace(encoded, "+", "%20", -1)
	encoded = strings.Replace(encoded, "*", "%2A", -1)
	encoded = strings.Replace(encoded, "%7E", "~", -1)

	return encoded
}
