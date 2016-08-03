package prize

import (
	"github.com/golang/glog"
	"pb"
)

type PrizeMails struct {
	UserId string
	mails  map[string]*PrizeMail
}

func NewPrizeMails(userId string) *PrizeMails {
	mails := &PrizeMails{}
	mails.UserId = userId
	mails.mails = make(map[string]*PrizeMail)

	prizeMails, err := FindPrizeMails(userId)
	if err != nil {
		glog.V(2).Info("FindPrizeMails err:", err)
		return mails
	}

	for _, mail := range prizeMails {
		mails.mails[mail.MailId] = mail
	}

	return mails
}

func (mails *PrizeMails) SaveMails() {
	for _, mail := range mails.mails {
		mail.UserId = mails.UserId
		SavePrizeMail(mail)
	}
}

func (mails *PrizeMails) BuildMessage() *pb.MsgGetPrizeMailListRes {
	msg := &pb.MsgGetPrizeMailListRes{}

	for _, mail := range mails.mails {
		msg.Mails = append(msg.Mails, mail.BuildMessage())
	}

	return msg
}

func (mails *PrizeMails) GetPrizeMail(mailId string) *PrizeMail {
	return mails.mails[mailId]
}

func (mails *PrizeMails) AddPrizeMail(mail *PrizeMail) {
	mail.UserId = mails.UserId
	mails.mails[mail.MailId] = mail
}

func (mails *PrizeMails) RemoveMail(mailId string) {
	delete(mails.mails, mailId)

	RemoveMail(mailId)
}
