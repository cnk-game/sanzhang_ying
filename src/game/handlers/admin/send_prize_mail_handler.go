package admin

import (
	"code.google.com/p/goprotobuf/proto"
	"config"
	"encoding/json"
	domainOfflineMsg "game/domain/offlineMsg"
	"game/domain/user"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"pb"
)

type PrizeMail struct {
	UserId    string `json:"userId"`
	Gold      int    `json:"gold"`
	Diamond   int    `json:"diamond"`
	Exp       int    `json:"exp"`
	Score     int    `json:"score"`
	ItemType  int    `json:"itemType"`
	ItemCount int    `json:"itemCount"`
	Content   string `json:"content"`
}

func SendPrizeMailHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		return
	}

	mail := &PrizeMail{}
	err = json.Unmarshal(b, mail)
	if err != nil || mail.UserId == "" {
		glog.Error(err)
		w.Write([]byte(`0`))
		return
	}

	glog.Info("===>SendPrizeMailHandler:", mail)

	_, err = user.FindByUserId(mail.UserId)
	if err != nil {
		glog.Error(err)
		w.Write([]byte(`0`))
		return
	}

	prizeMail := &pb.PrizeMailDef{}
	prizeMail.MailId = proto.String(bson.NewObjectId().Hex())
	prizeMail.Content = proto.String(mail.Content)
	prizeMail.Prize = &pb.PrizeDef{}
	prizeMail.Prize.Gold = proto.Int(mail.Gold)
	prizeMail.Prize.Diamond = proto.Int(mail.Diamond)
	prizeMail.Prize.Exp = proto.Int(mail.Exp)
	prizeMail.Prize.Score = proto.Int(mail.Score)
	if mail.ItemType > 0 && mail.ItemCount > 0 {
		if mail.ItemType == 1 {
			prizeMail.Prize.ItemType = pb.MagicItemType_FOURFOLD_GOLD.Enum()
		} else if mail.ItemType == 2 {
			prizeMail.Prize.ItemType = pb.MagicItemType_PROHIBIT_COMPARE.Enum()
		} else if mail.ItemType == 3 {
			prizeMail.Prize.ItemType = pb.MagicItemType_REPLACE_CARD.Enum()
		}
		prizeMail.Prize.ItemCount = proto.Int(mail.ItemCount)
	}
	domainOfflineMsg.SaveOfflineMsgLog(mail.UserId, mail.Gold, mail.Diamond)

	user.GetPlayerManager().SendServerMsg("", []string{mail.UserId}, int32(pb.ServerMsgId_MQ_PRIZE_MAIL), prizeMail)

	w.Write([]byte(`1`))
}
