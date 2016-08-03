package pay

import (
	"crypto/md5"
	"fmt"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"io"
	"strconv"
	"sync"
	"time"
	"util"
)

type PAY_TOKEN_INFO struct {
	UserId    string
	Token     string
	ProductId string
}

type TokenManager struct {
	sync.RWMutex
	tokens map[string]*PAY_TOKEN_INFO
}

var tokenManager *TokenManager

func Init() {
	tokenManager = &TokenManager{}
	tokenManager.tokens = make(map[string]*PAY_TOKEN_INFO)
}

func GetTokenManager() *TokenManager {
	return tokenManager
}

const (
	orderLogC = "order_log"
)

type OrderLog struct {
	UserId    string    `bson:"userId"`
	ProductId string    `bson:"pId"`
	Order     string    `bson:"order"`
	Time      time.Time `bson:"time"`
}

func SaveGetOrderLog(userId string, productId string, order string) error {
	now := time.Now()
	cur_C := orderLogC + "_" + strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day()))
	log := OrderLog{userId, productId, order, now}

	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(&log)
	})
}

func (this *TokenManager) GetToken(userId string, productId string) string {
	this.Lock()
	defer this.Unlock()

	tm := time.Now()
	sign := "QF" + userId + fmt.Sprintf("%d-%d-%d %02d:%02d:%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	token := Md5func(sign)

	token = string(token[0:30])

	item := &PAY_TOKEN_INFO{}
	item.UserId = userId
	item.ProductId = productId
	item.Token = token

	glog.Info("GetToken token=", token)

	this.tokens[token] = item

	SaveGetOrderLog(userId, productId, token)

	return token
}

func (this *TokenManager) CheckToken(token string) (bool, string) {
	this.Lock()
	defer this.Unlock()

	tokenInfo, ok := this.tokens[token]
	if !ok {
		glog.Info("CheckToken filed, token=", token)
		return false, ""
	}

	uId := tokenInfo.UserId
	delete(this.tokens, token)

	return true, uId
}

func (this *TokenManager) CommonCheckToken(token string) (bool, string, string) {
	this.Lock()
	defer this.Unlock()

	tokenInfo, ok := this.tokens[token]
	if !ok {
		glog.Info("CheckToken filed, token=", token)
		return false, "", ""
	}

	uId := tokenInfo.UserId
	productId := tokenInfo.ProductId

	return true, uId, productId
}

func (this *TokenManager) GetProductId(token string) (bool, string, string) {
	this.Lock()
	defer this.Unlock()

	tokenInfo, ok := this.tokens[token]
	if !ok {
		glog.Info("CheckToken filed, token=", token)
		return false, "", ""
	}

	uId := tokenInfo.UserId
	productId := tokenInfo.ProductId

	return true, uId, productId
}

func Md5func(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
