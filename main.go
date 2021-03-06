/*
 *@author 小王
 *@Version 0.1
 *@Update time 2014.09.04
 *@golang 微信公众平台GOlang SDK


 */
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/hprose/hprose-go/hprose"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type WebWeChat struct {
	Token   string
	cookies []*http.Cookie
	Err     string
}
type User_info struct {
	Nick_name         string
	User_name         string
	Original_username string
	Signature         string
	Email             string
	Tencent_id        string
	Tencent_nick      string
	Fake_id           int
	Is_vip            int
	Is_wx_verify      int
	Is_dev_user       int
}

func main() {
	service := hprose.NewHttpService()
	service.AddFunction("WcMessageInfo", WcMessageInfo)
	service.AddFunction("WcOwnInfo", WcOwnInfo)
	service.AddFunction("WcSendMsg_Text", WcSendMsg_Text)
	service.AddFunction("WcGroupSendMsg_Text", WcGroupSendMsg_Text)
	service.AddFunction("WcBand", WcBand)
	service.AddFunction("GetQrcode", GetQrcode)
	service.AddFunction("GetAvatar", GetAvatar)
	service.AddFunction("GetFollowUserGroup", GetFollowUserGroup)
	var port string = "1245"
	fmt.Println("开始监听" + port + "端口")
	http.ListenAndServe(":"+port, service)
}

//获取网页Token
func GetToken(u, p string) (WCahtReqR WebWeChat) {
	var ReqUrl string = "https://mp.weixin.qq.com/cgi-bin/login"
	pwdmd5 := md5.New()
	pwdmd5.Write([]byte(p))
	hex.EncodeToString(pwdmd5.Sum(nil))
	var data string = "username=" + u + "&pwd=" + hex.EncodeToString(pwdmd5.Sum(nil)) + "&imgcode=&f=json"
	req, _ := http.NewRequest("POST", ReqUrl, strings.NewReader(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Referer", "https://mp.weixin.qq.com/")
	loc := []string{}
	redirect := func(req *http.Request, via []*http.Request) error {
		loc = append(loc, req.URL.Path)
		return fmt.Errorf("重定向取消")
	}
	tr := &http.Transport{}
	client := &http.Client{
		Transport:     tr,
		CheckRedirect: redirect,
	}
	resp, _ := client.Do(req)
	respbodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	respjson, _ := simplejson.NewJson(respbodyByte)
	base := respjson.Get("base_resp")
	errinfo := base.Get("err_msg").MustString()
	if errinfo == "need verify code" {
		WCahtReqR.Err = errinfo
		return WCahtReqR
	}
	WCahtReqR.cookies = resp.Cookies()
	WCahtReqR.Token = strings.Split(respjson.Get("redirect_url").MustString(), "=")[3]
	return
}

//获取用户二维码
func GetQrcode(u, p, fakeid string) (respbodyByte []byte) {
	WCahtReqR := GetToken(u, p)
	var ReqUrl string = "https://mp.weixin.qq.com/misc/getqrcode?fakeid=" + fakeid + "&token=" + WCahtReqR.Token + "&style=1"
	req, _ := http.NewRequest("GET", ReqUrl, nil)
	for i := range WCahtReqR.cookies {
		req.AddCookie(WCahtReqR.cookies[i])
	}
	loc := []string{}
	redirect := func(req *http.Request, via []*http.Request) error {
		loc = append(loc, req.URL.Path)
		return fmt.Errorf("重定向取消")
	}
	tr := &http.Transport{}
	client := &http.Client{
		Transport:     tr,
		CheckRedirect: redirect,
	}
	resp, _ := client.Do(req)
	respbodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

//获取用户头像
func GetAvatar(u, p, fakeid string) (respbodyByte []byte) {
	WCahtReqR := GetToken(u, p)
	var ReqUrl string = "https://mp.weixin.qq.com/misc/getheadimg?fakeid=" + fakeid + "&token=" + WCahtReqR.Token + "&style=1"
	req, _ := http.NewRequest("GET", ReqUrl, nil)
	for i := range WCahtReqR.cookies {
		req.AddCookie(WCahtReqR.cookies[i])
	}
	loc := []string{}
	redirect := func(req *http.Request, via []*http.Request) error {
		loc = append(loc, req.URL.Path)
		return fmt.Errorf("重定向取消")
	}
	tr := &http.Transport{}
	client := &http.Client{
		Transport:     tr,
		CheckRedirect: redirect,
	}
	resp, _ := client.Do(req)
	respbodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

//给制定fakeid发送消息
func WcSendMsg_Text(u, p, content, tofakeid string) bool {
	WCahtReqR := GetToken(u, p)
	var ReqUrl string = "https://mp.weixin.qq.com/cgi-bin/singlesend?t=ajax-response&f=json&token=" + WCahtReqR.Token + "&lang=zh_CN"
	var data string = "token=" + WCahtReqR.Token + "&lang=zh_CN&f=json&ajax=1&type=1&random=0.037916635" + RandM() + "6162031&type=1&content=" + content + "&tofakeid=" + tofakeid + "&imgcode="
	req, _ := http.NewRequest("POST", ReqUrl, strings.NewReader(data))
	for i := range WCahtReqR.cookies {
		req.AddCookie(WCahtReqR.cookies[i])
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Referer", "https://mp.weixin.qq.com/cgi-bin/singlesendpage?t=message/send&action=index&tofakeid=807054600&token="+WCahtReqR.Token+"&lang=zh_CN")
	loc := []string{}
	redirect := func(req *http.Request, via []*http.Request) error {
		loc = append(loc, req.URL.Path)
		return fmt.Errorf("重定向取消")
	}
	tr := &http.Transport{}
	client := &http.Client{
		Transport:     tr,
		CheckRedirect: redirect,
	}
	resp, _ := client.Do(req)
	respbodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	respjson, _ := simplejson.NewJson(respbodyByte)
	base := respjson.Get("base_resp")
	if base.Get("err_msg").MustString() == "ok" {
		return true
	}
	return false
}

//发送群消息
func WcGroupSendMsg_Text(u, p, content string) bool {
	var operation_seq string
	WCahtReqR := GetToken(u, p)
	var ReqUrl string = "https://mp.weixin.qq.com/cgi-bin/masssendpage?t=mass/send&token=" + WCahtReqR.Token + "&lang=zh_CN&f=json"
	req, _ := http.NewRequest("GET", ReqUrl, nil)
	for i := range WCahtReqR.cookies {
		req.AddCookie(WCahtReqR.cookies[i])
	}
	req.Header.Set("Referer", "https://mp.weixin.qq.com/cgi-bin/masssendpage?t=mass/send&token="+WCahtReqR.Token+"&lang=zh_CN")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		respdatabyte, _ := ioutil.ReadAll(resp.Body)
		respjson, err := simplejson.NewJson(respdatabyte)
		if err != nil {
			fmt.Println(err.Error())
		}
		operation_seqint := respjson.Get("operation_seq").MustInt()
		operation_seq = strconv.Itoa(operation_seqint)
	}

	ReqUrl = "https://mp.weixin.qq.com/cgi-bin/masssend?t=ajax-response&token=" + WCahtReqR.Token + "&lang=zh_CN"
	var data string = "token=" + WCahtReqR.Token + "&lang=zh_CN&f=json&ajax=1&random=0.9823722" + RandM() + "99422729&type=1&content=" + content + "&cardlimit=&sex=&groupid=&synctxweibo=0&country=&province=&city=&imgcode=&operation_seq=" + operation_seq
	req, _ = http.NewRequest("POST", ReqUrl, strings.NewReader(data))
	for i := range WCahtReqR.cookies {
		req.AddCookie(WCahtReqR.cookies[i])
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Referer", "https://mp.weixin.qq.com/cgi-bin/masssendpage?t=mass/send&token="+WCahtReqR.Token+"&lang=zh_CN")

	loc := []string{}
	redirect := func(req *http.Request, via []*http.Request) error {
		loc = append(loc, req.URL.Path)
		return fmt.Errorf("重定向取消")
	}
	tr := &http.Transport{}
	client = &http.Client{
		Transport:     tr,
		CheckRedirect: redirect,
	}
	resp, _ = client.Do(req)
	respbodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	respjson, _ := simplejson.NewJson(respbodyByte)
	base := respjson.Get("base_resp")
	if base.Get("err_msg").MustString() == "ok" {
		return true
	}
	return false
}

//获取用户信息
func WcOwnInfo(u, p string) (userinfo User_info) {
	WCahtReqR := GetToken(u, p)
	var WcOwnInfoUrlData string = "t=setting/index&action=index&token=" + WCahtReqR.Token + "&lang=zh_CN&f=json"
	var WcOwnInfoUrl string = "https://mp.weixin.qq.com/cgi-bin/settingpage?" + WcOwnInfoUrlData
	req, _ := http.NewRequest("GET", WcOwnInfoUrl, nil)
	for i := range WCahtReqR.cookies {
		req.AddCookie(WCahtReqR.cookies[i])
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		respdatabyte, _ := ioutil.ReadAll(resp.Body)
		respjson, err := simplejson.NewJson(respdatabyte)
		if err != nil {
			fmt.Println(err.Error())
		}
		userjsonobj, _ := respjson.CheckGet("user_info")
		settingjsonobj, _ := respjson.CheckGet("setting_info")
		setting_emailjsonobj, _ := settingjsonobj.CheckGet("bind_email")
		setting_blogjsonobj, _ := settingjsonobj.CheckGet("micro_blog")
		setting_introjsonobj, _ := settingjsonobj.CheckGet("intro")
		var userinfo User_info
		userinfo.Fake_id = userjsonobj.Get("fake_id").MustInt()
		userinfo.Nick_name = userjsonobj.Get("nick_name").MustString()
		userinfo.User_name = userjsonobj.Get("user_name").MustString()
		userinfo.Signature = setting_introjsonobj.Get("signature").MustString()
		userinfo.Is_dev_user = userjsonobj.Get("is_dev_user").MustInt()
		userinfo.Is_vip = userjsonobj.Get("is_vip").MustInt()
		userinfo.Is_wx_verify = userjsonobj.Get("is_wx_verify").MustInt()
		userinfo.Email = setting_emailjsonobj.Get("account").MustString()
		userinfo.Original_username = settingjsonobj.Get("original_username").MustString()
		userinfo.Tencent_id = setting_blogjsonobj.Get("tencent_id").MustString()
		userinfo.Tencent_nick = setting_blogjsonobj.Get("tencent_nick").MustString()
		return userinfo
	}
	return
}

/*
返回的结构
{
    "msg_item": [
        {
            "id": 200716384,
            "type": 1,
            "fakeid": "1168435081",
            "nick_name": "妖精的魅惑",
            "date_time": 1410941778,
            "content": "额，不好意思，那我取消关注",
            "source": "",
            "msg_status": 4,
            "has_reply": 0,
            "refuse_reason": "",
            "multi_item": [],
            "to_uin": 3012020639,
            "send_stat": {
                "total": 0,
                "succ": 0,
                "fail": 0
            }
        },
        {
            "id": 200716252,
            "type": 1,
            "fakeid": "1168435081",
            "nick_name": "妖精的魅惑",
            "date_time": 1410939154,
            "content": "/:,@P记不住全名",
            "source": "",
            "msg_status": 4,
            "has_reply": 0,
            "refuse_reason": "",
            "multi_item": [],
            "to_uin": 3012020639,
            "send_stat": {
                "total": 0,
                "succ": 0,
                "fail": 0
            }
        }
    ]
}
*/
//最近接受到的消息
func WcMessageInfo(u, p, count, day string) (msgstring string) { //获取近期信息
	WCahtReqR := GetToken(u, p)
	if WCahtReqR.Err != "" {
		return "err"
	}
	var WcMessageInfoUrlData string = "t=message/list&count=" + count + "&day=" + day + "&token=" + WCahtReqR.Token + "&lang=zh_CN&f=json"
	var WcMessageInfoUrl string = "https://mp.weixin.qq.com/cgi-bin/message?" + WcMessageInfoUrlData
	req, _ := http.NewRequest("GET", WcMessageInfoUrl, nil)
	for i := range WCahtReqR.cookies {
		req.AddCookie(WCahtReqR.cookies[i])
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		respdatabyte, _ := ioutil.ReadAll(resp.Body)
		respjson, err := simplejson.NewJson(respdatabyte)
		if err != nil {
			fmt.Println(err.Error())
		}
		msgstring := respjson.Get("msg_items").MustString()
		return msgstring
	}
	return
}

//用户绑定相关
func wcbandresp(WCahtReqR WebWeChat) (msgstring string) { //结合绑定使用
	var ReqUrl string = "https://mp.weixin.qq.com/advanced/advanced?action=interface&t=advanced/interface&token=" + WCahtReqR.Token + "&lang=zh_CN&f=json"
	req, _ := http.NewRequest("GET", ReqUrl, nil)
	for i := range WCahtReqR.cookies {
		req.AddCookie(WCahtReqR.cookies[i])
	}
	req.Header.Set("Referer", "https://mp.weixin.qq.com/advanced/advanced?action=dev&t=advanced/dev&token="+WCahtReqR.Token+"&lang=zh_CN")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		respdatabyte, _ := ioutil.ReadAll(resp.Body)
		respjson, err := simplejson.NewJson(respdatabyte)
		if err != nil {
			fmt.Println(err.Error())
		}
		msgstring := respjson.Get("operation_seq").MustInt()
		return strconv.Itoa(msgstring)
	}
	return
}
func WcBand(u, p, url, token string) bool {
	WCahtReqR := GetToken(u, p)
	if WCahtReqR.Err != "" {
		fmt.Println(u + "登录错误")
		return false
	}
	var operation_seq string = wcbandresp(WCahtReqR)
	var ReqUrl string = "https://mp.weixin.qq.com/advanced/callbackprofile?t=ajax-response&token=" + WCahtReqR.Token + "&lang=zh_CN"
	var data string = "callback_token=" + token + "&url=" + url + "&operation_seq=" + operation_seq
	req, _ := http.NewRequest("POST", ReqUrl, strings.NewReader(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://mp.weixin.qq.com/advanced/advanced?action=interface&t=advanced/interface&token="+WCahtReqR.Token+"&lang=zh_CN")
	loc := []string{}
	for i := range WCahtReqR.cookies {
		req.AddCookie(WCahtReqR.cookies[i])
	}
	redirect := func(req *http.Request, via []*http.Request) error {
		loc = append(loc, req.URL.Path)
		return fmt.Errorf("重定向取消")
	}
	tr := &http.Transport{}
	client := &http.Client{
		Transport:     tr,
		CheckRedirect: redirect,
	}
	resp, _ := client.Do(req)
	if resp.StatusCode == 200 {
		respbodyByte, _ := ioutil.ReadAll(resp.Body)
		respjson, err := simplejson.NewJson(respbodyByte)
		if err != nil {
			fmt.Println(err.Error())
		}
		msgstring := respjson.Get("ret").MustString()
		if msgstring == "0" {
			return true
		}
	}
	return false
}

func GetFollowUser(u, p, page, pageidx string) { //pageidx为第几页,从0开始
	//WCahtReqR := GetToken(u, p)
	//var ReqUrl string = "https://mp.weixin.qq.com/cgi-bin/contactmanage?t=user/index&pagesize=" + page + "&pageidx=0&type=0&token=" + WCahtReqR.Token + "&lang=zh_CN&f=json"

}

/*
返回信息
{
    "groups": [
        {
            "id": 0,
            "name": "未分组",
            "cnt": 6
        },
        {
            "id": 1,
            "name": "黑名单",
            "cnt": 0
        },
        {
            "id": 2,
            "name": "星标组",
            "cnt": 0
        },
        {
            "id": 100,
            "name": "内测用户",
            "cnt": 1
        }
    ]
}
*/
func GetFollowUserGroup(u, p, page string) string {
	WCahtReqR := GetToken(u, p)
	var ReqUrl string = "https://mp.weixin.qq.com/cgi-bin/contactmanage?t=user/index&pagesize=" + page + "&pageidx=0&type=0&token=" + WCahtReqR.Token + "&lang=zh_CN&f=json"
	req, _ := http.NewRequest("GET", ReqUrl, nil)
	for i := range WCahtReqR.cookies {
		req.AddCookie(WCahtReqR.cookies[i])
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		respdatabyte, _ := ioutil.ReadAll(resp.Body)
		respjson, err := simplejson.NewJson(respdatabyte)
		if err != nil {
			fmt.Println(err.Error())
		}
		group_list := respjson.Get("group_list").MustString()
		return group_list
	}
	return "error"
}
func RandM() string {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10)
	return strconv.Itoa(r)
}
