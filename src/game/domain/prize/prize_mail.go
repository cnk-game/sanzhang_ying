package prize

import (
	"code.google.com/p/goprotobuf/proto"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"util"
)

type PrizeMail struct {
	UserId    string `bson:"userId"`
	MailId    string `bson:"mailId"`
	Content   string `bson:"content"`
	Gold      int    `bson:"gold"`
	Diamond   int    `bson:"diamond"`
	Exp       int    `bson:"exp"`
	Score     int    `bson:"score"`
	ItemType  int    `bson:"itemType"`
	ItemCount int    `bson:"itemCount"`
	hashCode  *util.HashCode
}

const (
	prizeMailC = "prize_mail"
)

func (mail *PrizeMail) HashCode() *util.HashCode {
	return mail.hashCode
}

func (mail *PrizeMail) SetHashCode(hashCode *util.HashCode) {
	mail.hashCode = hashCode
}

func FindPrizeMails(userId string) ([]*PrizeMail, error) {
	mails := []*PrizeMail{}
	err := util.WithUserCollection(prizeMailC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).All(&mails)
	})

	if err == nil {
		for _, mail := range mails {
			mail.SetHashCode(util.NewHashCode(mail))
		}
	}
	return mails, err
}

func SavePrizeMail(mail *PrizeMail) error {
	hashCode := util.NewHashCode(mail)
	if mail.HashCode() != nil && mail.HashCode().Compare(hashCode) {
		return nil
	}

	return util.WithSafeUserCollection(prizeMailC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": mail.UserId, "mailId": mail.MailId}, mail)
		if err == nil {
			// 保存成功
			mail.SetHashCode(hashCode)
		}
		return err
	})
}

func RemoveMail(mailId string) error {
	return util.WithUserCollection(prizeMailC, func(c *mgo.Collection) error {
		return c.Remove(bson.M{"mailId": mailId})
	})
}

func (mail *PrizeMail) BuildMessage() *pb.PrizeMailDef {
	msg := &pb.PrizeMailDef{}
	msg.MailId = proto.String(mail.MailId)
	msg.Content = proto.String(mail.Content)
	msg.Prize = &pb.PrizeDef{}
	msg.Prize.Gold = proto.Int(mail.Gold)
	msg.Prize.Diamond = proto.Int(mail.Diamond)
	msg.Prize.Exp = proto.Int(mail.Exp)
	msg.Prize.Score = proto.Int(mail.Score)

	if mail.ItemType >= 1 && mail.ItemType <= 3 {
		if mail.ItemType == 1 {
			msg.Prize.ItemType = pb.MagicItemType_FOURFOLD_GOLD.Enum()
		} else if mail.ItemType == 2 {
			msg.Prize.ItemType = pb.MagicItemType_PROHIBIT_COMPARE.Enum()
		} else if mail.ItemType == 3 {
			msg.Prize.ItemType = pb.MagicItemType_REPLACE_CARD.Enum()
		}
		msg.Prize.ItemCount = proto.Int(mail.ItemCount)
	}

	return msg
}
