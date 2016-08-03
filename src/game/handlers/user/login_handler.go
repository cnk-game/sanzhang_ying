package user

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"config"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	domainGame "game/domain/game"
	newUserTask "game/domain/newusertask"
	"game/domain/offlineMsg"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"

	"pb"
	"strconv"
	"strings"
	"time"
	"util"
)

const (
	retryTimes       = 5
	retryMillisecond = 500
	AppID            = "100026251"
	IOS_CHANNEL      = "178"
	CAOHUA_CHANNEL   = "184"
	CaoHuaAppId      = "144"
	CaoHuaAppKey     = "AB2CBFCF1FE60C5900504D9FECA2A868"
	LESHI_CHANNEL    = "173"
	LeshiAppId       = "500098"
	LeshiAppKey      = "d4d7855c08de4ceead84ff28e132f1c3"

	XUNLEI_CHANNEL = "186"
	XunleiAppId    = "050283"
	XunleiAppKey   = "zdJKkuCNsd5luaurCSIJTi2e52nunMPA"

	HAIMA_IOS_CHANNEL = "187"
	HaimaIosAppId     = "a022b0bda1fbba6b65f8ca2d291e2f1b"
	HaimaIosAppKey    = "017e2413864f7c4e807afee9f2435dc0"

	JINLI_CHANNEL  = "212"
	jinliApikey    = "918AD85BD7214741B04B0E2ABA1AD3AD"
	jinliSecretKey = "437D75979E9F4F75882453CBEA67BA52"
	jinliHost      = "id.gionee.com"
	jinliPort      = "443"
	jinliMethod    = "POST"
	jinliUrl       = "/account/verify.do"

	MEIZU_CHANNEL  = "213"
	MeizuApiID     = "2906651"
	MeizuAppSecret = "xIF0GIyisU3E1j3GeyKDNI864friYVQ6"
	MeizuHost      = "https://api.game.meizu.com/%20game%20/security/checksession"

	LIANXIANG_CHANNEL = "214"

	KUPAI_CHANNEL  = "215"
	KupaiApiID     = "5000002955"
	KupaiSecretKey = "52886cb546764689b87de1ad747cc6a3"
	KupaiHost      = "https://openapi.coolyun.com/oauth2/token"

	IOS51_CHANNEL = "216"
	IOS51_APPID   = "100001045"
	IOS51Host     = "http://api.51pgzs.com/passport/checkLogin.php"
)

func genPassword(pwd string) string {
	h := md5.New()
	io.WriteString(h, config.ControlKey+pwd)
	return fmt.Sprintf("%x", h.Sum(nil))
}

type IOSCheckLoginRes struct {
	Ret   int    `json:"ret"`
	Error string `json:"error"`
}

func calcMd5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func IOSCheckLogin(uid, token string) (int, string) {
	param := "appid=" + AppID + "&token=" + token + "&uid=" + uid
	body := bytes.NewBuffer([]byte(param))
	res, err := http.Post("http://passport.xyzs.com/checkLogin.php", "application/x-www-form-urlencoded;", body)
	if err != nil {
		glog.Info("IOSCheckLogin err:", err)
		return -1, "POST接口错误"
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		glog.Info("IOSCheckLogin err:", err)
		return -1, "接口返回值错误"
	}

	resmsg := IOSCheckLoginRes{}
	err = json.Unmarshal(result, &resmsg)
	if err != nil {
		glog.Info("IOSCheckLogin err:", err)
		return -1, "接口返回值解析错误"
	}

	return resmsg.Ret, resmsg.Error
}

func CaoHuaLogin(uid, token string) bool {
	nowStr := fmt.Sprintf("%v", time.Now().Unix())
	tempStr := CaoHuaAppId + uid + token + nowStr + CaoHuaAppKey
	calcSign := calcMd5(tempStr)
	calcSign = strings.ToUpper(calcSign)
	sendStr := "AppId=" + CaoHuaAppId + "&PUserID=" + uid + "&Token=" + token + "&Times=" + nowStr + "&Sign=" + calcSign
	sendStr = "http://api.caohua.com/Api/TokenCheck.ashx?" + sendStr

	res, err := http.Get(sendStr)
	if err != nil {
		glog.Info("CaoHuaLogin err:", err)
		return false
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		glog.Info("CaoHuaLogin err:", err)
		return false
	}

	resultSrt := fmt.Sprintf("%v", result)
	glog.Info("CaoHuaLogin result:", resultSrt)
	index := strings.IndexAny(resultSrt, uid)
	if index != -1 {
		return true
	} else {
		glog.Info("CaoHuaLogin result index", index)
		return false
	}
}

func hmacSHA1Encrypt(encryptKey string, encryptText string) string {
	key := []byte("437D75979E9F4F75882453CBEA67BA52")
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(encryptText))
	fmt.Printf("%x\n", mac.Sum(nil))
	return string(mac.Sum(nil))
}

func macSig(host string, port string, macKey string, timestamp string, nonce string, method string, uri string) string {
	// 1. build mac string
	// 2. hmac-sha1
	// 3. base64-encoded

	text := timestamp + "\n" + nonce + "\n" + strings.ToUpper(method) + "\n" + uri + "\n" + strings.ToLower(host) + "\n" + port + "\n" + "\n"

	ciphertext := hmacSHA1Encrypt(macKey, text)

	sig := make([]byte, 28)

	base64.StdEncoding.Encode(sig, []byte(ciphertext))

	return string(sig)
}

func builderAuthorization() string {

	rand.Seed(time.Now().UnixNano())

	mac := macSig(jinliHost, jinliPort, jinliSecretKey, strconv.FormatInt(time.Now().Unix(), 10), strconv.Itoa(rand.Int())[0:8], jinliMethod, jinliUrl)
	mac = strings.Replace(mac, "\n", "", -1)

	authStr := "MAC " + "id=" + jinliApikey + ",ts=" + strconv.FormatInt(time.Now().Unix(), 10) + ",nonce=" + strconv.Itoa(rand.Int())[0:8] + ",mac=" + mac

	return authStr
}

func jinLiLogin(uid, token string) bool {
	//nowStr := fmt.Sprintf("%v", time.Now().Unix())

	if uid == "" || token == "" {
		return false
	}
	return true

	pool := x509.NewCertPool()
	caCertPath := "./handlers/user/certs/cert_server/ca.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return false
	}

	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair("./handlers/user/certs/cert_server/client.crt", "./handlers/user/certs/cert_server/client.key")
	if err != nil {
		fmt.Println("Loadx509keypair err:", err)
		return false
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}

	client := &http.Client{Transport: tr}
	//resp, err := client.Post("https://id.gionee.com:443/account/verify.do","application/x-www-form-urlencoded",
	//strings.NewReader(string.format("Authorization=\"%s\"",builderAuthorization())));

	//resp, err := client.PostForm("https://id.gionee.com:443/account/verify.do",	url.Values{"Authorization" : {builderAuthorization()}});

	req, err := http.NewRequest("POST", "https://id.gionee.com:443/account/verify.do", strings.NewReader(token))
	if err != nil {
		fmt.Println("Get error:", err)
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", builderAuthorization())

	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	resultSrt := fmt.Sprintf("%v", body)
	glog.Info("JInLiLogin result:", resultSrt)
	index := strings.IndexAny(resultSrt, uid)
	if index != -1 {
		return true
	} else {
		glog.Info("JInLiLogin result index", index)
		return false
	}
}

type KupaiRes struct {
	AccessToken  interface{} `json:"access_token"`
	Expiresin    interface{} `json:"expires_in"`
	RefreshToken interface{} `json:"refresh_token"`
	Openid       interface{} `json:"openid"`
}

type MeizuRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

type MeizuAccessReq struct {
	AppID     int64  `json:"app_id"`
	SessionId string `json:"session_id"`
	Uid       int64  `json:"uid"`
	Ts        int64  `json:"ts"`
	SignType  string `json:"sign_type"`
	Sign      string `json:"sign"`
}

type Ios51Res struct {
	Ret   string `json:"ret"`
	Error string `json:"error"`
}

func Ios51Login(uid, token string) bool {

	//	v := url.Values{}
	//	v.Set("appid", IOS51_APPID)
	//	v.Set("uid", uid)
	//	v.Set("token", token)
	//	body := ioutil.NopCloser(strings.NewReader(v.Encode())) //把form数据编下码
	//	client := &http.Client{}

	//	glog.Info("reqstr=", body)

	//	req, err := http.NewRequest("POST", IOS51Host, body)
	//	if err != nil {
	//		fmt.Println("IOS51 post error:", err)
	//		return false
	//	}

	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	fmt.Printf("%+v\n", req)
	//	resp, err := client.Do(req)

	param := "appid=" + IOS51_APPID + "&token=" + token + "&uid=" + uid
	body := bytes.NewBuffer([]byte(param))
	resp, err := http.Post(IOS51Host, "application/x-www-form-urlencoded;", body)

	defer resp.Body.Close()

	resbody, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resbody))

	resultSrt := fmt.Sprintf("%v", resbody)
	glog.Info("Ios51Login result:", resultSrt)

	res1 := &Ios51Res{}
	err = json.Unmarshal(resbody, res1)
	if err != nil {
		glog.Info("Ios51Login result Unmarshal err:", err)
		return false
	}

	ret := fmt.Sprintf("%v", res1.Ret)
	glog.Info("Ios51Login ret= ", ret)

	if ret == "0" {
		glog.Info("Ios51Login success ")
		return true
	} else {
		glog.Info("Ios51Login fail ")
		return false
	}

	//return true
}

func MeizuLogin(uid, token string) bool {

	//	client := &http.Client{}

	//	reqstr := MeizuAccessReq{}
	//	reqstr.AppID = 2906651
	//	reqstr.SessionId = token
	//	intuid, _ := strconv.ParseInt(uid, 10, 64)
	//	reqstr.Uid = intuid
	//	reqstr.Ts = time.Now().Unix()
	//	reqstr.SignType = "md5"

	//	h := md5.New()
	//	signdata := ""
	//	signdata += "app_id=" + MeizuApiID
	//	signdata += "&session_id=" + token
	//	signdata += "&uid=" + uid
	//	signdata += "&ts=" + fmt.Sprintf("%v", reqstr.Ts)
	//	signdata += ":" + MeizuAppSecret

	//	io.WriteString(h, signdata)
	//	localSign := fmt.Sprintf("%x", h.Sum(nil))

	//	reqstr.Sign = localSign

	//	glog.Info("reqstr=", reqstr)

	//	sendStr, _ := json.Marshal(reqstr)

	//	req, err := http.NewRequest("POST", MeizuHost, strings.NewReader(string(sendStr)))
	//	if err != nil {
	//		fmt.Println("Meizu post error:", err)
	//		return false
	//	}

	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	req.Header.Set("Authorization", builderAuthorization())

	//	resp, err := client.Do(req)

	//	defer resp.Body.Close()
	//	body, err := ioutil.ReadAll(resp.Body)
	//	fmt.Println(string(body))

	//	resultSrt := fmt.Sprintf("%v", body)
	//	glog.Info("MeizuLogin result:", resultSrt)

	//	res1 := &MeizuRes{}
	//	err = json.Unmarshal(body, res1)
	//	if err != nil {
	//		glog.Info("meizuLogin result Unmarshal err:", err)
	//		return false
	//	}

	//	Code := fmt.Sprintf("%v", res1.Code)
	//	glog.Info("MeizuLogin Code= ", Code)

	//	if Code != "200" {
	//		glog.Info("MeizuLogin success ")
	//		return true
	//	} else {
	//		glog.Info("MeizuLogin fail ")
	//		return false
	//	}

	return true
}

func kupaiLogin(uid, token string) (bool, string, string) {

	sendStr := KupaiHost + "?grant_type=authorization_code" + "&client_id=" + KupaiApiID + "&redirect_uri=" + KupaiSecretKey +
		"&client_secret=" + KupaiSecretKey + "&code=" + token
	glog.Info("sendStr=", sendStr)

	res, err := http.Get(sendStr)
	if err != nil {
		glog.Info("kupaiLogin err:", err)
		return false, "", ""
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		glog.Info("kupaiLogin err:", err)
		return false, "", ""
	}

	res1 := &KupaiRes{}
	err = json.Unmarshal(result, res1)
	if err != nil {
		glog.Info("kupaiLogin result Unmarshal err:", err)
		return false, "", ""
	}

	Openid := fmt.Sprintf("%v", res1.Openid)
	AccessToken := fmt.Sprintf("%v", res1.AccessToken)
	glog.Info("kupaiLogin Openid ", Openid)
	glog.Info("kupaiLogin access_token ", AccessToken)

	if Openid != "" && res1.Openid != nil && AccessToken != "" && res1.AccessToken != nil {
		glog.Info("kupaiLogin success ")
		return true, Openid, AccessToken

	} else {
		glog.Info("kupaiLogin fail ")
		return false, "", ""
	}
}

type LeshiRes struct {
	Code interface{} `json:"code"`
	Msg  interface{} `json:"msg"`
}

func LeshiLogin(token string) (bool, string) {
	sendStr := "http://youxi.letv.com/api/member/verify.html?access_token=" + token
	res, err := http.Get(sendStr)
	if err != nil {
		glog.Info("LeshiLogin err:", err)
		return false, ""
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		glog.Info("LeshiLogin err:", err)
		return false, ""
	}

	//result = fmt.Sprintf("%v", result)
	//glog.Info("LeshiLogin result:", result)

	res1 := &LeshiRes{}
	err = json.Unmarshal(result, res1)
	if err != nil {
		glog.Info("LeshiLogin Unmarshal err:", err)
		return false, ""
	}

	code := fmt.Sprintf("%v", res1.Code)
	glog.Info("LeshiLogin code ", code)
	resultMsg := fmt.Sprintf("%v", res1.Msg)
	if code == "0" {
		glog.Info("LeshiLogin success ")
		return true, ""

	} else {
		glog.Info("LeshiLogin fail ", resultMsg)
		return false, resultMsg
	}
}

func XunleiLogin(uid, token string) (bool, string) {
	sendStr := "http://websvr.niu.xunlei.com/checkAppUser.gameUserInfo?gameid=050283&customerid=" + uid + "&customerKey=" + token
	res, err := http.Get(sendStr)
	if err != nil {
		glog.Info("XunleiLogin err:", err)
		return false, ""
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		glog.Info("XunleiLogin err:", err)
		return false, ""
	}

	res1 := &LeshiRes{}
	err = json.Unmarshal(result, res1)
	if err != nil {
		glog.Info("XunleiLogin Unmarshal err:", err)
		return false, ""
	}

	code := fmt.Sprintf("%v", res1.Code)
	glog.Info("XunleiLogin code ", code)
	resultMsg := fmt.Sprintf("%v", res1.Msg)
	if code == "0" {
		glog.Info("XunleiLogin success ")
		return true, ""

	} else {
		glog.Info("XunleiLogin fail ", resultMsg)
		return false, resultMsg
	}
}

func HaiMaIosLogin(uid, token string) (bool, string) {
	param := "appid=" + HaimaIosAppId + "&t=" + token + "&uid=" + uid
	body := bytes.NewBuffer([]byte(param))
	res, err := http.Post("http://api.haimawan.com/index.php?m=api&a=validate_token", "application/x-www-form-urlencoded;", body)
	if err != nil {
		glog.Info("HaiMaIosLogin err:", err)
		return false, "POST接口错误"
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		glog.Info("HaiMaIosLogin err:", err)
		return false, "接口返回值错误"
	}

	resultStr := fmt.Sprintf("%v", result)
	resultStr = string(result)
	glog.Info("HaiMaIosLogin resultStr:", resultStr)
	if strings.Contains(resultStr, "success") {
		return true, ""
	} else {
		return false, ""
	}
}

func LoginHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	msg := &pb.MsgLoginReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.MsgLoginRes{}
	vs := msg.GetVersionName()
	imeiReq := msg.GetDeviceId()
	kupaiOpenid := ""
	if strings.EqualFold(vs, "1.0.1") {
		res.Code = pb.MsgLoginRes_UPDATE_NOTIFY.Enum()
		res.Reason = proto.String("您的游戏版本太低，请升级！")
		server.BuildClientMsg(m.GetMsgId(), res)
	}

	if strings.EqualFold(vs, "1.0.0") {
		res.Code = pb.MsgLoginRes_UPDATE_FORCE.Enum()
		res.Reason = proto.String("您的游戏版本太低，请升级！")
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	player := domainUser.NewPlayer() //game_player.go

	glog.Info("===>用户登录消息sess:", sess, " msg:", msg)
	glog.Info("===>用户登录消息 uid: ", msg.GetUid())
	glog.Info("===>用户登录消息 token: ", msg.GetToken())
	if msg.GetChannelId() == IOS_CHANNEL {
		ret, err_str := IOSCheckLogin(msg.GetUid(), msg.GetToken())
		if ret != 0 {
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			res.Reason = proto.String(err_str)
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
		userName, pwd, err := domainUser.FindIOSUser(msg.GetUid())
		if err != nil && err != mgo.ErrNotFound {
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			res.Reason = proto.String("系统错误")
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		if err == mgo.ErrNotFound {
			userName = "IOS" + msg.GetUid()
			pwd = userName
		}
		msg.Username = proto.String(userName)
		msg.Userpwd = proto.String(pwd)
	} else if msg.GetChannelId() == CAOHUA_CHANNEL {
		bLogin := CaoHuaLogin(msg.GetUid(), msg.GetToken())
		if bLogin == false {
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			res.Reason = proto.String("系统错误了")
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {
			userName, pwd, err := domainUser.FindCaoHuaUser(msg.GetUid())
			if err != nil && err != mgo.ErrNotFound {
				res.Code = pb.MsgLoginRes_FAILED.Enum()
				res.Reason = proto.String("系统错误")
				return server.BuildClientMsg(m.GetMsgId(), res)
			}

			if err == mgo.ErrNotFound {
				userName = "CAOHUA" + msg.GetUid()
				pwd = userName
			}
			msg.Username = proto.String(userName)
			msg.Userpwd = proto.String(pwd)
		}
	} else if msg.GetChannelId() == LESHI_CHANNEL {
		bLogin, resultMsg := LeshiLogin(msg.GetToken())
		if bLogin == false {
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			res.Reason = proto.String(resultMsg)
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {
			userName, pwd, err := domainUser.FindLeshiUser(msg.GetUid())
			if err != nil && err != mgo.ErrNotFound {
				res.Code = pb.MsgLoginRes_FAILED.Enum()
				res.Reason = proto.String("系统错误")
				return server.BuildClientMsg(m.GetMsgId(), res)
			}

			if err == mgo.ErrNotFound {
				userName = "LESHI" + msg.GetUid()
				pwd = userName
			}
			msg.Username = proto.String(userName)
			msg.Userpwd = proto.String(pwd)
		}
	} else if msg.GetChannelId() == XUNLEI_CHANNEL {
		bLogin, resultMsg := XunleiLogin(msg.GetUid(), msg.GetToken())
		if bLogin == false {
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			res.Reason = proto.String(resultMsg)
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {
			userName, pwd, err := domainUser.FindXunleiUser(msg.GetUid())
			if err != nil && err != mgo.ErrNotFound {
				res.Code = pb.MsgLoginRes_FAILED.Enum()
				res.Reason = proto.String("系统错误")
				return server.BuildClientMsg(m.GetMsgId(), res)
			}

			if err == mgo.ErrNotFound {
				userName = "XUNLEI" + msg.GetUid()
				pwd = userName
			}
			msg.Username = proto.String(userName)
			msg.Userpwd = proto.String(pwd)
		}
	} else if msg.GetChannelId() == HAIMA_IOS_CHANNEL {
		bLogin, resultMsg := HaiMaIosLogin(msg.GetUid(), msg.GetToken())
		if bLogin == false {
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			res.Reason = proto.String(resultMsg)
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {
			userName, pwd, err := domainUser.FindHMIosUser(msg.GetUid())
			if err != nil && err != mgo.ErrNotFound {
				res.Code = pb.MsgLoginRes_FAILED.Enum()
				res.Reason = proto.String("系统错误")
				return server.BuildClientMsg(m.GetMsgId(), res)
			}

			if err == mgo.ErrNotFound {
				userName = "HMIOS" + msg.GetUid()
				pwd = userName
			}
			msg.Username = proto.String(userName)
			msg.Userpwd = proto.String(pwd)
		}
	} else if msg.GetChannelId() == JINLI_CHANNEL {
		bLogin := jinLiLogin(msg.GetUid(), msg.GetToken())
		if bLogin == false {
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			res.Reason = proto.String("系统错误了")
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {
			userName, pwd, err := domainUser.FindJinLiUser(msg.GetUid())
			if err != nil && err != mgo.ErrNotFound {
				res.Code = pb.MsgLoginRes_FAILED.Enum()
				res.Reason = proto.String("系统错误")
				return server.BuildClientMsg(m.GetMsgId(), res)
			}

			if err == mgo.ErrNotFound {
				userName = "JINLI" + msg.GetUid()
				pwd = userName
			}

			msg.Username = proto.String(userName)
			msg.Userpwd = proto.String(pwd)
		}
	} else if msg.GetChannelId() == KUPAI_CHANNEL {
		bLogin, Openid, acctoken := kupaiLogin(msg.GetUid(), msg.GetToken())
		kupaiOpenid = Openid
		if bLogin == false {
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			res.Reason = proto.String("系统错误了")
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {

			userName, pwd, err := domainUser.FindKupaiUser(kupaiOpenid)
			if err != nil && err != mgo.ErrNotFound {
				res.Code = pb.MsgLoginRes_FAILED.Enum()
				res.Reason = proto.String("系统错误")
				return server.BuildClientMsg(m.GetMsgId(), res)
			}

			if err == mgo.ErrNotFound {
				userName = "KUPAI" + msg.GetUid()
				pwd = userName
			}
			msg.Username = proto.String(userName)
			msg.Userpwd = proto.String(pwd)

			res.Userid = proto.String(kupaiOpenid)
			res.Token = proto.String(acctoken)
		}
	} else if msg.GetChannelId() == MEIZU_CHANNEL {
		bLogin := MeizuLogin(msg.GetUid(), msg.GetToken())

		if bLogin == false {
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			res.Reason = proto.String("系统错误了")
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {

			userName, pwd, err := domainUser.FindMeizuUser(msg.GetUid())
			if err != nil && err != mgo.ErrNotFound {
				res.Code = pb.MsgLoginRes_FAILED.Enum()
				res.Reason = proto.String("系统错误")
				return server.BuildClientMsg(m.GetMsgId(), res)
			}

			if err == mgo.ErrNotFound {
				userName = "MEIZU" + msg.GetUid()
				pwd = userName
			}
			msg.Username = proto.String(userName)
			msg.Userpwd = proto.String(pwd)

		}
	} else if msg.GetChannelId() == IOS51_CHANNEL {
		bLogin := Ios51Login(msg.GetUid(), msg.GetToken())

		if bLogin == false {
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			res.Reason = proto.String("系统错误了")
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {

			userName, pwd, err := domainUser.FindIos51User(msg.GetUid())
			if err != nil && err != mgo.ErrNotFound {
				res.Code = pb.MsgLoginRes_FAILED.Enum()
				res.Reason = proto.String("系统错误")
				return server.BuildClientMsg(m.GetMsgId(), res)
			}

			if err == mgo.ErrNotFound {
				userName = "MEIZU" + msg.GetUid()
				pwd = userName
			}
			msg.Username = proto.String(userName)
			msg.Userpwd = proto.String(pwd)

		}
	}

	if msg.GetUsername() == "" {
		glog.Info("===>登录无效,username为空sess:", sess)
		return nil
	}

	// 用户名登录
	u, err := domainUser.FindByUserName(msg.GetUsername()) //user.go
	player.User = u
	glog.Info("==>查找用户username:", msg.GetUsername(), " err:", err, " u:", u)
	if err != nil {
		if err == mgo.ErrNotFound {
			if msg.GetIsSwitchAccount() {
				res.Code = pb.MsgLoginRes_USERNAME_OR_PASSWORD_ERROR.Enum()
				return server.BuildClientMsg(m.GetMsgId(), res)
			}

			imeiCount, imeiEr := domainUser.GetDeviceUserCount(imeiReq)
			if imeiEr != mgo.ErrNotFound {
				if imeiCount >= 3 {
					res.Code = pb.MsgLoginRes_DEVICE_MORE_USER.Enum()
					return server.BuildClientMsg(m.GetMsgId(), res)
				}
			}
			nickName := msg.GetNickname()
			if nickName == "" {
				nickName = msg.GetModel()
			}
			if nickName == "" {
				glog.Info("==>登录结果:Nickname不能为空!sess:", sess)
				res.Code = pb.MsgLoginRes_FAILED.Enum()
				return server.BuildClientMsg(m.GetMsgId(), res)
			}
			if msg.GetRobotKey() == config.RobotKey {
				// 机器人
				u.IsRobot = true
				if msg.GetRobotGender() == 0 {
					u.Gender = int(pb.Gender_BOY)
				} else {
					u.Gender = int(pb.Gender_GIRL)
				}
				u.PhotoUrl = msg.GetRobotPhoto()
				u.RobotVipLevel = int(msg.GetRobotVip())
			} else {
				if rand.Float64() < 0.5 {
					u.Gender = int(pb.Gender_BOY)
					u.PhotoUrl = fmt.Sprintf("%v", 6+rand.Int()%3)
				} else {
					u.Gender = int(pb.Gender_GIRL)
					u.PhotoUrl = fmt.Sprintf("%v", 1+rand.Int()%5)
				}
			}
			// 用户不存在，直接创建
			//u.UserId = bson.NewObjectId().Hex()
			u.UserId = domainUser.GetNewUserId()
			u.UserName = msg.GetUsername()

			u.Password = genPassword(msg.GetUserpwd())
			u.Nickname = nickName
			u.CreateTime = time.Now()
			u.ChannelId = msg.GetChannelId()
			u.Model = msg.GetModel()

			err = domainUser.SaveUser(u) //user.go
			if err != nil {
				glog.Info("保存用户失败err:", err, " user:", u)
				res.Code = pb.MsgLoginRes_FAILED.Enum()
				return server.BuildClientMsg(m.GetMsgId(), res)
			} else {
				domainUser.SaveUserNameIdByUserName(u.UserId, u.UserName)
				domainUser.SaveDeivceUserCount(imeiReq, imeiCount+1)
			}

			if msg.GetChannelId() == IOS_CHANNEL {
				domainUser.InsertIOSUser(msg.GetUid(), u.UserName, msg.GetUserpwd())
			} else if msg.GetChannelId() == CAOHUA_CHANNEL {
				domainUser.InsertChaoHuaUser(msg.GetUid(), u.UserName, msg.GetUserpwd())
			} else if msg.GetChannelId() == LESHI_CHANNEL {
				domainUser.InsertLeshiUser(msg.GetUid(), u.UserName, msg.GetUserpwd())
			} else if msg.GetChannelId() == XUNLEI_CHANNEL {
				domainUser.InsertXunleiUser(msg.GetUid(), u.UserName, msg.GetUserpwd())
			} else if msg.GetChannelId() == HAIMA_IOS_CHANNEL {
				domainUser.InsertHMIosUser(msg.GetUid(), u.UserName, msg.GetUserpwd())
			} else if msg.GetChannelId() == JINLI_CHANNEL {
				domainUser.InsertJinLiUser(msg.GetUid(), u.UserName, msg.GetUserpwd())
			} else if msg.GetChannelId() == KUPAI_CHANNEL {
				domainUser.InsertKupaiUser(kupaiOpenid, u.UserName, msg.GetUserpwd())
			} else if msg.GetChannelId() == MEIZU_CHANNEL {
				domainUser.InsertMeizuUser(msg.GetUid(), u.UserName, msg.GetUserpwd())
			} else if msg.GetChannelId() == IOS51_CHANNEL {
				domainUser.InsertIos51User(msg.GetUid(), u.UserName, msg.GetUserpwd())
			}

			onCreatePlayer(player, msg.GetRobotWinTimes(), msg.GetRobotLoseTimes(), msg.GetRobotCurDayEarnGold(), msg.GetRobotCurWeekEarnGold(), msg.GetRobotMaxCards())
		} else {
			// 查找用户失败
			glog.Error("查找用户失败username:", msg.GetUsername(), " err:", err)
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
	} else {
		// 用户存在，校验密码
		if msg.GetRobotKey() != config.RobotKey && u.Password != genPassword(msg.GetUserpwd()) {
			glog.Info("===>检验密码失败username:", msg.GetUsername(), " pwd:", u.Password, "recvPwd:", genPassword(msg.GetUserpwd()))
			res.Code = pb.MsgLoginRes_USERNAME_OR_PASSWORD_ERROR.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
	}

	if u.IsLocked {
		glog.Info("==>userId:", u.UserId, "账号被锁定")
		res.Code = pb.MsgLoginRes_LOCKED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	glog.V(2).Info("===>登录成功player userId:", player.User.UserId)

	if sess.LoggedIn {
		glog.Info("===>user:", msg.GetUsername(), "已经登录!")
		if sess.OnLogout != nil {
			sess.OnLogout()
		}
	}

	player.SessKey = bson.NewObjectId().Hex()

	if !domainUser.GetPlayerManager().AddItem(player.User.UserId, player.User.IsRobot, sess) {
		glog.Info("===>玩家已在线，踢掉userId:", player.User.UserId, " IP:", sess.IP, " sessKey:", player.SessKey, " sess:", sess)
		msg := &pb.MsgShowTips{}
		msg.Cmd = proto.Int(1)
		msg.UserId = proto.String(player.User.UserId)
		msg.Content = proto.String("您的账号在其他设备登陆，请确认个人账号的安全。")
		msg.Buttons = append(msg.Buttons, "确定")
		msg.Buttons = append(msg.Buttons, "退出")
		domainUser.GetPlayerManager().SendClientMsg(player.User.UserId, int32(pb.MessageId_SHOW_TIPS), msg)

		domainUser.GetPlayerManager().Kickout(player.User.UserId)

		ok := false
		for i := 0; i < 10; i++ {
			time.Sleep(150 * time.Millisecond)
			if !domainUser.GetPlayerManager().AddItem(player.User.UserId, player.User.IsRobot, sess) {
				continue
			}
			ok = true
			break
		}

		glog.Info("===>继续处理登录userId:", player.User.UserId, " ok:", ok, " sessKey:", player.SessKey, " sess:", sess)

		if !ok {
			glog.Info("登录失败，踢出上一玩家失败userId:", player.User.UserId)
			res.Code = pb.MsgLoginRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
	}

	// 加载用户信息
	if !player.OnLogin() {
		domainUser.GetPlayerManager().DelItem(player.User.UserId, player.User.IsRobot)
		glog.Info("userId:", player.User.UserId, " OnLogin失败!")
		res.Code = pb.MsgLoginRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	glog.Info("==>玩家登录成功userId:", player.User.UserId, " deviceId:", msg.GetDeviceId(), " sessKey:", player.SessKey, " sess:", sess)

	if player.NewPlayer {
		domainUser.GetUserFortuneManager().EarnGold(player.User.UserId, 10000, "新账号") //新手奖励改成10000 modify by yelong
		player.User.UpgradePrizeVersion = msg.GetVersionName()
		newUserTask.GetNewUserTaskManager().InitUserTask(player.User.UserId)
	}

	umToken := msg.GetUmToken()
	if umToken != "" {
		domainUser.SavePushUser(umToken, player.User.UserId)
	}
	sess.LoggedIn = true
	sess.OnLogout = player.OnLogout
	sess.Data = player
	player.LoginIP = sess.IP
	player.LoginDeviceId = msg.GetDeviceId()

	player.SendToClientFunc = func(msgId int32, body proto.Message) {
		sess.SendToClient(server.BuildClientMsg(int32(msgId), body))
	}
	player.OnLogoutFunc = onLogout

	player.SetUserCache()
	//起凡的QQ：3284317919
	//台湾iOS：2750076349

	helloContent := fmt.Sprintf("欢迎%v进入游戏！预祝您游戏愉快！加客服QQ：3284317919 领取新手兑换码！", player.User.Nickname)
	glog.Info("==>玩家登录成功 send hello", helloContent)
	player.SendToClient(int32(pb.MessageId_CHAT), util.BuildSysBugle(helloContent))

	if player.User.IsRobot {
		domainUser.GetFakeRankingList().SetRobot(player.User.UserId)
		domainUser.GetBackgroundUserManager().SetUser(player.User.UserId, true)
	}

	if !player.NewPlayer && !player.User.IsRobot {
		domainUser.GetUserFortuneManager().CheckVipLevel(player.User.UserId)
	}

	go processOfflineMsg(player.User.UserId)

	glog.Info("===>登录结果:LoadUserTask ")
	newUserTask.GetNewUserTaskManager().LoadUserTask(player.User.UserId)

	glog.Info("===>登录结果:成功userId:", player.User.UserId, " sess:", sess)

	// 登录成功
	res.Code = pb.MsgLoginRes_OK.Enum()
	res.ServerTime = proto.Int64(time.Now().Unix())
	return server.BuildClientMsg(m.GetMsgId(), res)
}

func onCreatePlayer(p *domainUser.GamePlayer, robotWinTimes, robotLoseTimes, robotCurDayEarnGold, robotCurWeekEarnGold int32, robotMaxCards []int32) {
	p.NewPlayer = true

	p.UserLog = &domainUser.UserLog{}
	p.UserLog.UserId = p.User.UserId
	p.UserLog.UserName = p.User.UserName
	p.UserLog.CreateTime = p.User.CreateTime
	p.UserLog.Model = p.User.Model
	p.UserLog.Channel = p.User.ChannelId

	matchRecord := &domainUser.MatchRecord{}
	matchRecord.UserId = p.User.UserId
	matchRecord.WinTimes = int(robotWinTimes)
	matchRecord.LoseTimes = int(robotLoseTimes)
	matchRecord.CurDayEarnGold = int(robotCurDayEarnGold)
	matchRecord.CurWeekEarnGold = int(robotCurWeekEarnGold)
	matchRecord.MaxCards = robotMaxCards
	matchRecord.DayEarnGoldResetTime = time.Now()
	matchRecord.WeekEarnGoldResetTime = time.Now()
	domainUser.SaveMatchRecord(matchRecord)
}

func onLogout(userId string) {
	// 登出
	go domainGame.GetDeskManager().OnOffline(userId)
}

func processOfflineMsg(userId string) {
	msgs, err := offlineMsg.FindOfflineMsg(userId)
	if err != nil {
		return
	}

	for _, msg := range msgs {
		serverMsg := &pb.ServerMsg{}
		err = proto.Unmarshal(msg.MsgBody, serverMsg)
		if err != nil {
			continue
		}
		domainUser.GetPlayerManager().SendServerMsg2("", []string{userId}, msg.MsgId, serverMsg.MsgBody)
	}
}
