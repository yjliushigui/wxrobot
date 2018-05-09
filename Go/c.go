package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"encoding/xml"
	"math/rand"
	"regexp"
	"io/ioutil"
	qrterminal "github.com/mdp/qrterminal"
	"os"
	"strconv"
	"errors"
	"bytes"
	json2 "encoding/json"
	"net/url"
	reflect2 "reflect"
)

type WxKey struct {
	AppId       string
	RedirectURI string
	Fun         string
	Lang        string
	time        int64
}
type Friend struct {
	ChatRoomId       int
	MemberCount      int
	RemarkName       string
	OwnerUin         int
	PYQuanPin        string
	NickName         string
	AppAccountFlag   int
	Statues          int
	Sex              int
	UserName         string
	HeadImgUrl       string
	ContactFlag      int
	HideInputBarFlag int
	DisplayName      string
	Uin              string
	MemberList       interface{}
	AttrStatus       int
	SnsFlag          int
	City             string
	VerifyFlag       int
	PYInitial        string
	Province         string
	Alias            string
	RemarkPYInitial  string
	RemarkPYQuanPin  string
	StarFriend       int
	IsOwner          int
	Signature        string
	UniFriend        int
	KeyWord          string
	EncryChatRoomId  int
}

var HttpHeader *string
var timeWX = time.Now().UnixNano() / 1000000
var timeWX13 = strconv.FormatInt(timeWX, 10)
var t = time.Now().Unix()
var timeWX9 = strconv.FormatInt(t, 10)
var urlChannel = make(chan string, 200)                                                                        //chan中存入string类型的href属性,缓冲200
var atagRegExp = regexp.MustCompile(`<a[^>]+[(href)|(HREF)]\s*\t*\n*=\s*\t*\n*[(".+")|('.+')][^>]*>[^<]*</a>`) //以Must前缀的方法或函数都是必须保证一定能执行成功的,否则将引发一次panic
var userAgent = [...]string{"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
	"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
	"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
	"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
	"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
	"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
	"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var webwxDataTicket string

var webwxAuthTicket string

var Response *ResponseData

type ResponseData struct {
	XMLName     xml.Name `xml:"error"`
	Ret         string   `xml:"ret"`
	Message     string   `xml:"message"`
	Skey        string   `xml:"skey"`
	Wxsid       string   `xml:"wxsid"`
	Wxuin       string   `xml:"wxuin"`
	PassTicket  string   `xml:"pass_ticket"`
	Isgrayscale string   `xml:"isgrayscale"`
}
var BaseRequest *BaseRequestData

//noinspection ALL
type BaseRequestData struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
}
func GetRandomUserAgent() string {
	return userAgent[r.Intn(len(userAgent))]
}

func GetHref(atag string) (href, content string) {
	inputReader := strings.NewReader(atag)
	decoder := xml.NewDecoder(inputReader)
	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		// 处理元素开始（标签）
		case xml.StartElement:
			for _, attr := range token.Attr {
				attrName := attr.Name.Local
				attrValue := attr.Value
				if strings.EqualFold(attrName, "href") || strings.EqualFold(attrName, "HREF") {
					href = attrValue
				}
			}
			// 处理元素结束（标签）
		case xml.EndElement:
			// 处理字符数据（这里就是元素的文本）
		case xml.CharData:
			content = string([]byte(token))
		default:
			href = ""
			content = ""
		}
	}
	return href, content
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect2.TypeOf(obj)
	v := reflect2.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
func DecodeWxXML(XMLContent []byte) (v *ResponseData, err error) {
	err = xml.Unmarshal(XMLContent, &v)
	if err == nil {

		return v, nil
	}
	return nil, err

}

/*处理cookie*/
func getCookieData(cookies []*http.Cookie) (webwxDataTicket string, webwxAuthTicket string) {
	for _, cookie := range cookies {
		/*获取cookie的webwxDataTicket*/
		rule1 := `webwx_data_ticket=([0-9a-zA-Z+_/@]*)`
		reg1 := regexp.MustCompile(rule1)
		webwxDataTicketCookie := reg1.FindString(cookie.String())
		if webwxDataTicketCookie != "" {
			webwxDataTicket = strings.Replace(webwxDataTicketCookie, "webwx_data_ticket=", "", -1)
		}
		/*获取cookie的webwxAuthTicket*/
		rule2 := `webwx_auth_ticket=([0-9a-zA-Z+_/@]*)`
		reg2 := regexp.MustCompile(rule2)
		webwxAuthTicketCookie := reg2.FindString(cookie.String())
		if webwxAuthTicketCookie != "" {
			webwxAuthTicket = strings.Replace(webwxAuthTicketCookie, "webwx_auth_ticket=", "", -1)
		}
	}
	return webwxDataTicket, webwxAuthTicket

}
func WxRedirect(uuid string) string {
	wxinitUrl := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?uuid=" + uuid + "&tip=0&_=e'" + strconv.FormatInt(timeWX, 10)
	resp, _ := http.Get(wxinitUrl)
	page, _ := ioutil.ReadAll(resp.Body)
	ruleURI := `((http[s]{0,1}|ftp)://[a-zA-Z0-9\.\-]+\.([a-zA-Z]{2,4})(:\d+)?(/[a-zA-Z0-9\.\-~!@#$%^&*+?:_/=<>]*)?)|((www.)|[a-zA-Z0-9\.\-]+\.([a-zA-Z]{2,4})(:\d+)?(/[a-zA-Z0-9\.\-~!@#$%^&*+?:_/=<>]*)?)`
	regURI := regexp.MustCompile(ruleURI)
	resURI := regURI.FindAllStringSubmatch(string(page), -1)
	uriSplit := strings.Split(resURI[2][1], "scan")
	redirectUri := uriSplit[0] + "fun=new&scan=" + timeWX9
	httpRule := `(https://[0-9a-zA-Z]+\.qq\.com)/`
	httpRexp := regexp.MustCompile(httpRule)
	/*获取头部连接类型*/
	HHres := httpRexp.FindStringSubmatch(redirectUri)
	HHres[0] = HHres[0] + "cgi-bin/mmwebwx-bin/"
	HttpHeader = &HHres[0]
	return redirectUri
}
func getDeviceID() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	deviceID := rnd.Int63n(1000000000000000)
	return "e" + strconv.FormatInt(deviceID, 10)
}

func JsonMap(jsonData []byte) (Jmap map[string]interface{}, err error) {
	err = json2.Unmarshal(jsonData, &Jmap)
	return Jmap, err
}

func PostWX(URL string, param map[string]interface{}) (respContent interface{}, err error) {
	pJson, _ := json2.Marshal(param)
	jsonStr := bytes.NewBuffer([]byte(pJson))
	req, _ := http.NewRequest("POST", URL, jsonStr)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	page, _ := ioutil.ReadAll(resp.Body)
	return page, err
}
func WriteFile(filename string, content interface{}) bool {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return false
	}
	f.Write(content.([]byte))
	f.Close()
	return true
}

/**
获取R(取反)
 */
func getR(reverse bool) string {
	if reverse {
		return strconv.FormatInt(^time.Now().Unix()*1000, 10)
	} else {
		return strconv.FormatInt(time.Now().Unix()*1000, 10)
	}
}

func urlEncode(params map[string]interface{}) string {
	u := url.Values{}
	for key, value := range params {
		u.Set(key, value.(string))
	}
	return u.Encode()
}

//TODO
func main() {
	Start()
}

/**
1
 */
func Start() {
	uuid, err := getUuid()
	if err == nil {
		Qrcode(uuid)
		fmt.Println("二维码生成成功")
		fmt.Println("=========")
		fmt.Println("请用手机微信扫描二维码")
		Login(uuid)
	} else {
		fmt.Println(err.Error())
	}
}

/**
2
 */
func getUuid() (Uuid string, err error) {
	errors.New("获取uuid失败")
	wx := WxKey{"wx782c26e4c19acffb", "https://wx.qq.com/cgi-bin?mmwebwx-bin=webwxnewloginpage", "new", "zh_CN", timeWX}
	resp, _ := http.Get("https://login.wx.qq.com/jslogin?appid=" + wx.AppId + "&redirect_uri=" + wx.RedirectURI + "&fun=" + wx.Fun + "&lang=" + wx.Lang + "&_=" + strconv.FormatInt(timeWX, 10))
	page, _ := ioutil.ReadAll(resp.Body)
	ruleCode := `\d+`
	regCode := regexp.MustCompile(ruleCode)
	resCode := regCode.FindSubmatch(page)
	Code := string(resCode[0])
	if Code == "200" {
		/*获取uuid并生成相应的二维码*/
		ruleUuid := `(?sim:["'](.*?)==["'])`
		regUuid := regexp.MustCompile(ruleUuid)
		resUuid := regUuid.FindSubmatch(page)
		Uuid := string(resUuid[1]) + "=="
		return Uuid, nil
	}
	return "", err

}

/**
3
 */
func Qrcode(Uuid string) {
	QRcodeUrl := "https://login.weixin.qq.com/l/" + Uuid
	qrterminal.GenerateHalfBlock(QRcodeUrl, qrterminal.L, os.Stdout)
}

/**
4
 */
func Login(Uuid string) {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		loginUrl := "https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=" + Uuid + "&tip=0&r=" + strconv.FormatInt(^timeWX, 10) + "&_=" + strconv.FormatInt(timeWX, 10)
		resp, _ := http.Get(loginUrl)
		page, _ := ioutil.ReadAll(resp.Body)
		ruleCode := `\d+`
		regCode := regexp.MustCompile(ruleCode)
		resCode := regCode.FindSubmatch(page)
		if string(resCode[0]) == "201" {
			fmt.Println("========")
			fmt.Println("请在手机微信上点击登录！")
		} else if string(resCode[0]) == "200" {
			fmt.Println("========")
			fmt.Println("登录成功")
			ticker.Stop()
			time.Sleep(2 * time.Second)
			/*获取回调接口和cookie*/
			redirectURL := WxRedirect(Uuid)
			redirectPage, _ := http.Get(redirectURL)
			redirectData, _ := ioutil.ReadAll(redirectPage.Body)
			cookies := redirectPage.Cookies()
			webwxDataTicket, webwxAuthTicket = getCookieData(cookies)
			/*获取初始化数据*/
			Response, _ = DecodeWxXML(redirectData)
			BaseRequest = &BaseRequestData{Response.Wxuin, Response.Wxsid, Response.Skey, ""}
			//BaseRequest.Skey = Response.Skey
			//BaseRequest.Sid = Response.Wxsid
			//BaseRequest.Uin = Response.Wxuin

			fmt.Println("========")
			fmt.Println("初始化数据成功")
			fmt.Println(BaseRequest)
			fmt.Println("========")
			ret, _ := strconv.Atoi(Response.Ret)
			WxInit()
			if ret != 0 {
				fmt.Println("========")
				fmt.Println("获取失败")
				Start()
			}
		} else {
			fmt.Println("请用手机微信扫描二维码")
		}
	}
}

/**
5
 */
func WxInit() {
	fmt.Println("初始化微信数据")
	BaseRequest.DeviceID = getDeviceID()
	u := url.Values{}
	u.Set("pass_ticket", Response.PassTicket)
	u.Set("skey", Response.Skey)
	u.Set("r", getR(true))
	WxInitURL := *HttpHeader + "webwxinit?" + u.Encode()
	param := make(map[string]interface{})
	param["BaseRequest"] = *BaseRequest
	page, _ := PostWX(WxInitURL, param)
	respContent, err := JsonMap(page.([]byte))
	if err != nil {
		panic(err)
	}
	BaseResponse := respContent["BaseResponse"].(map[string]interface{})
	if int(BaseResponse["Ret"].(float64)) == 0 {
		//User:=respContent["User"].(map[string]interface{})
		fmt.Println("初始化成功")

	} else {
		fmt.Println("初始化失败")
		err := errors.New(BaseResponse["ErrMsg"].(string))
		panic(err)
	}
	//f, err := os.OpenFile("WXINFO/wxinit_data.txt", os.O_CREATE|os.O_RDWR, os.ModePerm)
	//if err != nil {
	//	panic(err)
	//}
	//f.Write(page)
	//f.Close()
	fmt.Println("开始同步心跳数据")
	synckey, synckeyList, _ := getInitSync(page.([]byte))
	fmt.Println(synckeyList)

	heart(synckey, synckeyList)
	//getContactList()
}

/**
发送心跳
 */
func heart(synckey string, synckeyList []interface{}) {
	var headerUrl string
	if *HttpHeader == "https://wx2.qq.com/cgi-bin/mmwebwx-bin/" {
		headerUrl = "https://webpush.wx2.qq.com/cgi-bin/mmwebwx-bin/"
	} else {
		headerUrl = "https://webpush.wx.qq.com/cgi-bin/mmwebwx-bin/"
	}

	params := make(map[string]interface{})
	params["skey"] = Response.Skey
	params["sid"] = Response.Wxsid
	params["uin"] = Response.Wxuin
	params["synckey"] = synckey
	params["deviceid"] = getDeviceID()
	params["pass_ticket"] = Response.PassTicket
	params["r"] = getR(true)
	params["_"] = getR(false)
	synckeyRequestUrl := headerUrl + "synccheck?" + urlEncode(params)

	resp, _ := http.Get(synckeyRequestUrl)
	page, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("返回的心跳数据")
	fmt.Println(string(page))
	ruleURI := `retcode:"(.*?)",selector:"(.*?)"`
	regURI := regexp.MustCompile(ruleURI)
	resURI := regURI.FindStringSubmatch(string(page))
	retCode, _ := strconv.Atoi(resURI[1])

	if retCode > 0 {
		fmt.Println("心跳同步失败！")
		//Start()
	}

	fmt.Println("心跳同步成功")

	selector, _ := strconv.Atoi(resURI[2])
	if selector > 0 {
		webwxsync(synckeyList)
	} else {
		heart(synckey, synckeyList)
	}

}
func webwxsync(synckeyList interface{}) {
	params := make(map[string]interface{})
	params["sid"] = Response.Wxsid
	params["skey"] = Response.Skey
	params["pass_ticket"] = Response.PassTicket

	wxsyncUrl := *HttpHeader + "webwxsync?" + urlEncode(params)
	postData := make(map[string]interface{})
	postData["BaseRequest"] = *BaseRequest
	MapSynckey := make(map[string]interface{})
	MapSynckey["Count"] = len(synckeyList.([]interface{}))
	MapSynckey["List"] = synckeyList
	postData["SyncKey"] = MapSynckey
	postData["rr"] = getR(true)
	page, _ := PostWX(wxsyncUrl, postData)
	responseData, _ := JsonMap(page.([]byte))
	fmt.Println("消息列表")
	fmt.Println(responseData["AddMsgList"])
	newSyncKeyMap := responseData["SyncKey"].(map[string]interface{})
	newSyncKeyList := newSyncKeyMap["List"].([]interface{})
	newSyncKey := Sync(newSyncKeyList)

	heart(newSyncKey, newSyncKeyList)

}

/**
7 联系人列表
 */
func getContactList() {
	params := make(map[string]interface{})
	params["pass_ticket"] = Response.PassTicket
	params["seq"] = 0
	params["skey"] = Response.Skey
	params["r"] = getR(true)
	ContactListUrl := *HttpHeader + "webwxgetcontact?" + urlEncode(params)

	param := make(map[string]interface{})
	BaseRequest.DeviceID = getDeviceID()
	param["BaseRequest"] = BaseRequest
	fmt.Print(param)
	page, err := PostWX(ContactListUrl, param)
	ContactList, _ := JsonMap(page.([]byte))
	for key, People := range ContactList["MemberList"].([]interface{}) {
		chatList := People.(map[string]interface{})
		fmt.Println("排序:", key, "昵称:", chatList["NickName"], "用户名:", chatList["UserName"])
	}
	if err != nil {
		panic(err)
	}
	//TODO 心跳检测输出最新消息
	//WriteFile("WXINFO/wxcontacnt_data.txt", page)
}

/**
6 心跳
 */
func getInitSync(page []byte) (string, []interface{}, interface{}) {
	//undoJson, _ := ioutil.ReadFile("WXINFO/wxinit_data.txt")
	var decodeJson map[string]interface{}
	json2.Unmarshal(page, &decodeJson)
	BaseResponse := decodeJson["BaseResponse"].(map[string]interface{})

	if int(BaseResponse["Ret"].(float64)) == 0 {
		User := decodeJson["User"].(map[string]interface{})
		fmt.Println(User)
		SyncKey := decodeJson["SyncKey"].(map[string]interface{})
		SyncKeyList := SyncKey["List"].([]interface{})
		fmt.Println("心跳列队")
		reqSync := Sync(SyncKeyList)
		fmt.Println(reqSync)
		fmt.Println("聊天列表")
		fmt.Println("活跃人联系列表")
		ContactList := decodeJson["ContactList"].([]interface{})
		for _, People := range ContactList {
			chatList := People.(map[string]interface{})
			fmt.Println("昵称:", chatList["NickName"])
			fmt.Println("用户名:", chatList["UserName"])
		}

		return reqSync, SyncKeyList, User
	} else {
		err := errors.New(BaseResponse["ErrMsg"].(string))
		panic(err)
	}

}

/**
拼接心跳数据
 */
func Sync(SyncKeyList []interface{}) (string) {
	var p string
	for i := 0; i < len(SyncKeyList); i++ {
		synckey := SyncKeyList[i].(map[string]interface{})
		K := synckey["Key"].(float64)
		key := strconv.FormatFloat(K, 'g', -1, 64)
		val := synckey["Val"].(float64)
		value := fmt.Sprintf("%.0f", val)
		if i == len(SyncKeyList)-1 {
			p += key + "_" + value
		} else {
			p += key + "_" + value + "|"
		}
	}

	return p

}
