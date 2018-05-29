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
}

//群组
var List []*map[string]interface{}
var RecentGroup []*map[string]interface{}
var GroupList []*map[string]interface{}
//联系人
var FriendList []*map[string]interface{}

var HttpHeader *string

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var webwxDataTicket string

var webwxAuthTicket string

var Owner *string

var Response *ResponseData

var User *map[string]interface{}

var StatusNotifyUserName *string

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

func main3() {
	//var s = "window.code=200;window.redirect_uri='https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage?ticket=AacYxqXsgigYZ_C-vsawz_Aj@qrticket_0&uuid=Qd3-9RHUaA==&lang=zh_CN&scan=1520228140';"
	var s = "https://wx2.qq.com/cgi"
	ruleURI := `(https://[0-9a-zA-Z]+\.qq\.com)/`
	//ruleURI := `((http[s]{0,1}|ftp)://[a-zA-Z0-9\.\-]+\.([a-zA-Z]{2,4})(:\d+)?(/[a-zA-Z0-9\.\-~!@#$%^&*+?:_/=<>]*)?)|((www.)|[a-zA-Z0-9\.\-]+\.([a-zA-Z]{2,4})(:\d+)?(/[a-zA-Z0-9\.\-~!@#$%^&*+?:_/=<>]*)?)`
	regURI := regexp.MustCompile(ruleURI)
	resURI := regURI.FindStringSubmatch(s)
	//url := strings.Split(resURI[2][1],"scan")
	fmt.Println(resURI)
	//
	//url := `https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage?ticket=AVLdWMJ9X-I7SKwXTfzMgEO0@qrticket_0&uuid=gY5QOs1sXg==&lang=zh_CN&scan=1520231578&fun=new&lang=zh_CN`;
	//resp, _ := http.Get(url)
	//
	//page, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(page))
	/*cookie分类*/
	//str := "cookie: wxuin=1449338181; Path=/; Domain=wx.qq.com; Expires=Mon, 05 Mar 2018 19:14:09 GMT " +
	//	"cookie: wxsid=m4lusmStRvwo6/7j; Path=/; Domain=wx.qq.com; Expires=Mon, 05 Mar 2018 19:14:09 GMT " +
	//	"cookie: wxloadtime=1520234049; Path=/; Domain=wx.qq.com; Expires=Mon, 05 Mar 2018 19:14:09 GMT cookie: mm_lang=zh_CN; Path=/; Domain=wx.qq.com; Expires=Mon, 05 Mar 2018 19:14:09 GMT " +
	//	"cookie: webwx_data_ticket=gSfVf3nMmzWxr8ztb+rY7YNf; Path=/; Domain=qq.com; Expires=Mon, 05 Mar 2018 19:14:09 GMT " +
	//	"cookie: webwxuvid=5338542d4d1d7a49844371eb3aca31f5415f946a7e24fedfdeab5e2ac2ec168678d7446d758ab8b9b23757e2ac05dd77; Path=/; Domain=wx.qq.com; Expires=Thu, 02 Mar 2028 07:14:09 GMT " +
	//	"cookie: webwx_auth_ticket=CIsBENLIhe4GGoABLJKYn+0AT956om9TnWOBSCQdwmzuxHcjYxIMqHpz2jTLkc6WfqgwPV9LdQpGrNKL0vPWXWNCmoV2Lu88ORKxnuawJkKQtBU7RFdmlKpom+XObAK35BNXO1eVtebcWo0nUXXAk6TnkrcLvSAt8GYHbcU4MjEzLKLivYeWo4/51Po=; Path=/; Domain=wx.qq.com; Expires=Thu, 02 Mar 2028 07:14:09 GMT"
	//rule := `(cookie: [0-9a-zA-Z_+/=]*=[0-9a-zA-Z_+/=]*)`
	//reg := regexp.MustCompile(rule)
	//res := reg.FindAllStringSubmatch(str, -1)
	//for i := 0; i < len(res); i++ {
	//	cookieRP := strings.Replace(res[i][0], "cookie: ", "", -1)
	//	/*获取cookie的webwxDataTicket*/
	//	rule1 := `webwx_data_ticket=([0-9a-zA-Z+_/@]*)`
	//	reg1 := regexp.MustCompile(rule1)
	//	webwxDataTicket := reg1.FindString(cookieRP)
	//	if webwxDataTicket != "" {
	//		webwxDataTicket = strings.Replace(webwxDataTicket, "webwx_data_ticket=", "", -1)
	//	}
	//
	//	/*获取cookie的webwxAuthTicket*/
	//	rule2 := `webwx_auth_ticket=([0-9a-zA-Z+_/@]*)`
	//	reg2 := regexp.MustCompile(rule2)
	//	webwxAuthTicket := reg2.FindString(cookieRP)
	//	if webwxAuthTicket != "" {
	//		webwxAuthTicket = strings.Replace(webwxAuthTicket, "webwx_auth_ticket=", "", -1)
	//	}
	//}
}

/**
Strucct TO MAP
 */
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect2.TypeOf(obj)
	v := reflect2.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

/**
解析XML
 */
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

/**
获取回调URL
 */
func WxRedirect(uuid string) string {
	wxinitUrl := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?uuid=" + uuid + "&tip=0&_=e'" + getR(false)
	resp, _ := http.Get(wxinitUrl)
	page, _ := ioutil.ReadAll(resp.Body)
	ruleURI := `((http[s]{0,1}|ftp)://[a-zA-Z0-9\.\-]+\.([a-zA-Z]{2,4})(:\d+)?(/[a-zA-Z0-9\.\-~!@#$%^&*+?:_/=<>]*)?)|((www.)|[a-zA-Z0-9\.\-]+\.([a-zA-Z]{2,4})(:\d+)?(/[a-zA-Z0-9\.\-~!@#$%^&*+?:_/=<>]*)?)`
	regURI := regexp.MustCompile(ruleURI)
	resURI := regURI.FindAllStringSubmatch(string(page), -1)
	uriSplit := strings.Split(resURI[2][1], "scan")
	redirectUri := uriSplit[0] + "fun=new&scan=" + getR(false)
	httpRule := `(https://[0-9a-zA-Z]+\.qq\.com)/`
	httpRexp := regexp.MustCompile(httpRule)
	/*获取头部连接类型*/
	HHres := httpRexp.FindStringSubmatch(redirectUri)
	HHres[0] = HHres[0] + "cgi-bin/mmwebwx-bin/"
	HttpHeader = &HHres[0]
	return redirectUri
}

/**
获取deviceid
 */
func getDeviceID() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	deviceID := rnd.Int63n(1000000000000000)
	return "e" + strconv.FormatInt(deviceID, 10)
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

var BaseRequest *BaseRequestData

//noinspection ALL
type BaseRequestData struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
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

func getMsgID() string {
	id := strconv.FormatInt(time.Now().UnixNano()/100, 10)
	return id
}

func checkIsGroup(UserName interface{}) bool {
	rule := "@@"
	res, _ := regexp.MatchString(rule, UserName.(string))
	return res
}

//TODO

func main() {
	Start()
}

func tulinv2() {
	URL := "http://openapi.tuling123.com/openapi/api/v2"
	postData := make(map[string]interface{})
	postData["reqType"] = "0"
	perception := make(map[string]interface{})
	inputText := make(map[string]interface{})
	inputText["text"] = "附近的酒店"
	perception["inputText"] = inputText
	postData["perception"] = perception
	userInfo := make(map[string]interface{})
	userInfo["apiKey"] = "04a857e6426946ea84b9fbbac3a40e2a"
	userInfo["userId"] = "c26fe0df3b5539be"
	postData["userInfo"] = userInfo
	res, _ := PostWX(URL, postData)
	data, _ := json2.Marshal(res)
	fmt.Println(string(data))
}
func tulin(text interface{}, UserName interface{},NickName interface{}) {
	postData := make(map[string]interface{})
	params := make(map[string]interface{})
	params["key"] = "04a857e6426946ea84b9fbbac3a40e2a"
	params["info"] = text
	URL := "http://www.tuling123.com/openapi/api?" + urlEncode(params)
	data, _ := PostWX(URL, postData)
	Mdata, _ := JsonMap(data.([]byte))
	sendMsg(UserName, Mdata["text"])
	fmt.Println("你对",NickName,"说->",Mdata["text"])
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
		fmt.Println(errors.New("获取UUID失败"))
	}
}

/**
2
 */
func getUuid() (Uuid string, err error) {
	errors.New("获取uuid失败")
	wx := WxKey{"wx782c26e4c19acffb", "https://wx.qq.com/cgi-bin?mmwebwx-bin=webwxnewloginpage", "new", "zh_CN"}
	resp, _ := http.Get("https://login.wx.qq.com/jslogin?appid=" + wx.AppId + "&redirect_uri=" + wx.RedirectURI + "&fun=" + wx.Fun + "&lang=" + wx.Lang + "&_=" + getR(false))
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
		loginUrl := "https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=" + Uuid + "&tip=0&r=" + getR(true) + "&_=" + getR(false)
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

			fmt.Println("=====================================================================")
			fmt.Println("初始化数据成功")
			fmt.Println("=====================================================================")
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
		UserMap := respContent["User"].(map[string]interface{})
		User = &UserMap
	} else {
		fmt.Println("初始化失败")
		Start()
	}
	fmt.Println("开始同步心跳数据")

	//格式化心跳的请求数据
	synckey, synckeyList, _ := getInitSync(page.([]byte))

	ch := make(chan bool)
	//发送第一次消息为了获得所有的群组
	statusNotify()

	synckey, synckeyList = firstHeart(synckey, synckeyList)

	go heart(synckey, synckeyList, ch)

	getContactList()

	operation()

	<-ch

	fmt.Println("退出程序")
	//getContactList()
}

/**
初始化发送第一次消息
 */
func statusNotify() {
	fmt.Println("发送第一次初始化信息...")
	params := make(map[string]interface{})
	params["lang"] = "zh_CN"
	params["pass_ticket"] = Response.PassTicket

	postData := make(map[string]interface{})
	postData["BaseRequest"] = *BaseRequest
	postData["Code"] = 3
	postData["FromUserName"] = Owner
	postData["ToUserName"] = Owner
	postData["ClientMsgId"] = getMsgID()
	statusUrl := *HttpHeader + "webwxstatusnotify?" + urlEncode(params)
	data, err := PostWX(statusUrl, postData)
	respContent, err := JsonMap(data.([]byte))
	BaseResponse := respContent["BaseResponse"].(map[string]interface{})
	if int(BaseResponse["Ret"].(float64)) != 0 {
		fmt.Println("初始化获取最近联系人信息失败.......")
		Start()
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("发送第一次初始化信息成功...")

}

func firstHeart(synckey string, synckeyList []interface{}) (string, []interface{}) {
	fmt.Println("开始第一次心跳......")
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

	ruleURI := `retcode:"(.*?)",selector:"(.*?)"`
	regURI := regexp.MustCompile(ruleURI)
	resURI := regURI.FindStringSubmatch(string(page))
	retCode, _ := strconv.Atoi(resURI[1])
	if retCode > 0 {
		fmt.Println("第一次心跳同步失败！")
		Start()
	}
	fmt.Println("第一次心跳成功")
	return firstWebwxsync(synckeyList)

}

func firstWebwxsync(synckeyList interface{}) (string, []interface{}) {
	fmt.Println("获取初始化联系人信息列表......")
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

	newSyncKeyMap := responseData["SyncKey"].(map[string]interface{})

	newSyncKeyList := newSyncKeyMap["List"].([]interface{})

	newSyncKey := Sync(newSyncKeyList)

	for _, msg := range responseData["AddMsgList"].([]interface{}) {

		Msg := msg.(map[string]interface{})

		if Msg["StatusNotifyUserName"] != nil {

			str := Msg["StatusNotifyUserName"].(string)

			StatusNotifyUserName = &str

			webwxbatchgetcontact(str)
		} else {
			Start()
		}
	}
	return newSyncKey, newSyncKeyList
}

/**
发送心跳
 */
func heart(synckey string, synckeyList []interface{}, ch chan bool) {
	for {
		select {
		default:
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

			data, _ := ioutil.ReadAll(resp.Body)

			fmt.Println("心跳同步中...")
			ruleURI := `retcode:"(.*?)",selector:"(.*?)"`
			regURI := regexp.MustCompile(ruleURI)
			resURI := regURI.FindStringSubmatch(string(data))
			retCode, _ := strconv.Atoi(resURI[1])

			if retCode > 0 {
				fmt.Println("心跳同步失败！")
				Start()
			}

			fmt.Println("心跳同步成功...")

			selector, _ := strconv.Atoi(resURI[2])
			if selector > 0 {
				webwxsync(synckeyList, ch)
			} else {
				heart(synckey, synckeyList, ch)
			}
		}
	}

}

func webwxsync(synckeyList interface{}, ch chan bool) {
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
	fmt.Println("获取到的消息")
	res, _ := json2.Marshal(page)
	fmt.Println(string(res))
	responseData, _ := JsonMap(page.([]byte))

	fmt.Println("获取消息列表")
	fmt.Println(responseData["AddMsgList"])

	for _, msg := range responseData["AddMsgList"].([]interface{}) {

		Msg := msg.(map[string]interface{})
		if checkIsGroup(Msg["FromUserName"]) {

		} else {
			if Msg["FromUserName"] ==*Owner {
				ToUser, _ := getUserInList(Msg["ToUserName"])
				fmt.Println("你对", ToUser["NickName"], "说->", Msg["Content"])
			} else {
				FromUser, err := getUserInList(Msg["FromUserName"])

				if err != nil {
					fmt.Println("您有新的消息：", "未知", "->", Msg["Content"])
				} else {
					fmt.Println("您有新的消息：", FromUser["NickName"], "->", Msg["Content"])
				}
				tulin(Msg["Content"], Msg["FromUserName"],FromUser["NickName"])

			}

		}
	}
	newSyncKeyMap := responseData["SyncKey"].(map[string]interface{})
	newSyncKeyList := newSyncKeyMap["List"].([]interface{})
	newSyncKey := Sync(newSyncKeyList)

	heart(newSyncKey, newSyncKeyList, ch)

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
		OwnerUserName := User["UserName"].(string)
		Owner = &OwnerUserName

		SyncKey := decodeJson["SyncKey"].(map[string]interface{})
		SyncKeyList := SyncKey["List"].([]interface{})
		reqSync := Sync(SyncKeyList)
		//fmt.Println("心跳列队")
		//fmt.Println(reqSync)
		//fmt.Println("聊天列表")
		//fmt.Println("活跃人联系列表")
		//ContactList := decodeJson["ContactList"].([]interface{})
		//for _, People := range ContactList {
		//	chatList := People.(map[string]interface{})
		//	fmt.Println("昵称:", chatList["NickName"])
		//	fmt.Println("用户名:", chatList["UserName"])
		//}

		return reqSync, SyncKeyList, User
	} else {
		err := errors.New(BaseResponse["ErrMsg"].(string))
		panic(err)
	}

}

/**
7 联系人列表
 */
func getContactList() {
	params := make(map[string]interface{})
	params["pass_ticket"] = Response.PassTicket
	params["seq"] = "0"
	params["skey"] = Response.Skey
	params["r"] = getR(true)
	ContactListUrl := *HttpHeader + "webwxgetcontact?" + urlEncode(params)

	param := make(map[string]interface{})
	BaseRequest.DeviceID = getDeviceID()
	param["BaseRequest"] = BaseRequest
	page, err := PostWX(ContactListUrl, param)
	ContactList, _ := JsonMap(page.([]byte))
	for _, People := range ContactList["MemberList"].([]interface{}) {
		chatList := People.(map[string]interface{})
		if checkIsGroup(chatList["UserName"]) {
			GroupList = append(GroupList, &chatList)
		} else {
			FriendList = append(FriendList, &chatList)
		}
	}
	if len(RecentGroup) > 0 {
		if len(GroupList) > 0 {
			for _, value := range RecentGroup {
				PValue := *value
				for _, Rvalue := range GroupList {
					PRvalue := *Rvalue
					if PValue["UserName"] != PRvalue["UserName"] {
						GroupList = append(GroupList, &PRvalue)
					}
				}
			}
		} else {
			GroupList = RecentGroup
		}
	}
	if err != nil || len(GroupList) < 1 {
		fmt.Println("获取联系人失败")
		fmt.Println("")
		fmt.Println("--------------")
		fmt.Println("1:继续    ", "|")
		fmt.Println("2:退出    ", "|")
		fmt.Println("--------------")
		var s int
		fmt.Scanf("%d", &s)
		if s == 2 {
			os.Exit(0)
		}
	}
}

/**
获取活跃人聊天列表
 */
func webwxbatchgetcontact(str string) {
	list := strings.Split(str, ",")
	params := make(map[string]interface{})
	params["type"] = "ex"
	params["lang"] = "zh_CN"
	params["r"] = getR(true)
	params["pass_ticket"] = Response.PassTicket
	PostData := make(map[string]interface{})
	PostData["BaseRequest"] = *BaseRequest
	PostData["Count"] = len(list)

	var MG []map[string]string
	for _, name := range list {
		M := make(map[string]string)
		M["UserName"] = string(name)
		M["EncryChatRoomId"] = ""
		MG = append(MG, M)
	}
	PostData["List"] = MG

	URL := *HttpHeader + "webwxbatchgetcontact?" + urlEncode(params)

	data, _ := PostWX(URL, PostData)
	Result, _ := JsonMap(data.([]byte))
	ContactList := Result["ContactList"]
	for _, People := range ContactList.([]interface{}) {
		chatList := People.(map[string]interface{})
		if checkIsGroup(chatList["UserName"]) {
			RecentGroup = append(RecentGroup, &chatList)
		}
	}
}

/**
发送文字消息
 */
func sendMsg(UserName interface{}, Msg interface{}) {
	param := make(map[string]interface{})
	param["pass_ticket"] = Response.PassTicket
	sendUrl := *HttpHeader + "webwxsendmsg?" + urlEncode(param)
	params := make(map[string]interface{})
	MsgMap := make(map[string]interface{})
	params["BaseRequest"] = *BaseRequest
	params["Scene"] = "0"
	MsgMap["Type"] = "1"
	MsgMap["Content"] = Msg
	MsgMap["FromUserName"] = *Owner
	MsgMap["ToUserName"] = UserName
	ID := getMsgID()
	MsgMap["LocalID"] = ID
	MsgMap["ClientMsgId"] = ID
	params["Msg"] = MsgMap
	_, err := PostWX(sendUrl, params)
	if err != nil {
		panic(err)
	}
}

/**
获取某个用户的信息
 */
func getUserInList(UserName interface{}) (User map[string]interface{}, err error) {
	for _, People := range FriendList {
		User = *People
		if User["UserName"] == UserName{
			return User, nil
		}
	}
	return nil, errors.New("未知联系人")
}

func operation() {
	var a int
	fmt.Println("-----------------")
	fmt.Println("1:发送消息       ", "|")
	fmt.Println("2:查看联系人列表  ", "|")
	fmt.Println("3:查看详细信息    ", "|")
	fmt.Println("-----------------")
	fmt.Println("请输入：")
	fmt.Scanf("%d", &a)
	switch a {
	case 1:
		fmt.Println("选择群或者联系人")
		fmt.Println("-----------------")
		fmt.Println("1:群组列表    ", "|")
		fmt.Println("2:联系人列表  ", "|")
		fmt.Println("-----------------")
		var k int
		fmt.Scanf("%d", &k)
		fmt.Println("联系人列表")
		var isGroup bool
		switch k {
		case 1:
			List = GroupList
			isGroup = true
			break
		case 2:
			List = FriendList
			isGroup = false
			break
		}
		fmt.Println("-------------------------------------------------------------------------------")
		for key, People := range List {
			chatList := *People
			fmt.Println("排序|", key, "|昵称|", chatList["NickName"], "|用户名|", chatList["UserName"], "|")
			fmt.Println("-------------------------------------------------------------------------------")
		}
		if !(len(List) > 0) {
			fmt.Println("你选择的列表为空")
			operation()
		}
		var f int
		fmt.Scanf("%d", &f)
		var s string

		chooseData := *List[f]

		if isGroup {
			//getGroupUserList(chooseData["UserName"])
		}
		fmt.Println("请输入你想对", chooseData["NickName"], "说的内容:")
		fmt.Scanf("%s", &s)
		fmt.Println("我->", chooseData["NickName"], ":", s)
		sendMsg(chooseData["UserName"].(string), s)
		fmt.Println("发送成功")
		operation()
		break
	case 2:
		fmt.Println("开发中........")
		operation()
		break
	default:
		fmt.Println("开发中........")
		operation()
		break

	}
}
