package prize

import (
	"code.google.com/p/goprotobuf/proto"
	domainPrize "game/domain/prize"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	//"gopkg.in/mgo.v2/bson"
	domainIosActive "game/domain/iosActive"
	"pb"
	"time"
	//"util"
	"bytes"
	"encoding/json"
	"net/http"
)

type SEND_MAIL struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func SendMail(subject, content string) bool {
	req := SEND_MAIL{}
	req.Subject = subject
	req.Content = content

	b, err := json.Marshal(req)
	if err != nil {
		glog.Info("SendMail err:", err)
		return false
	}

	body := bytes.NewBuffer([]byte(b))
	_, err1 := http.Post("http://127.0.0.1:9001/mail/simple", "application/json;charset=utf-8", body)
	if err1 != nil {
		glog.Info("SendMail err:", err1)
		return false
	} else {
		return true
	}
}

func ExchangeGoodsByScoreHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgExchangeScoreGoodsReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.MsgExchangeScoreGoodsRes{}

	itemId := int(msg.GetItemId())
	channelId := int(msg.GetChannelId())
	glog.Info("魅力兑换 itemId:", itemId)
	glog.Info("魅力兑换 channelId:", channelId)
	goodsName := ""
	Goods:=domainPrize.ExchangeGoods{};
	if channelId == 183 {
		goods, ok := domainIosActive.GetIosActiveManager().GetItemInfo(itemId)
		if !ok {
			glog.Error("找不到兑换配置itemId:", itemId)
			res.Code = pb.MsgExchangeScoreGoodsRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
		goodsName = goods.Name

		if goods.Remains == 0 {
			glog.Error("兑换的商品已经没有了userId:", player.User.UserId, "|count=", goods)
			res.Code = pb.MsgExchangeScoreGoodsRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		if !domainUser.GetUserFortuneManager().ConsumeCharm(player.User.UserId, goods.Price) {
			glog.Error("消耗魅力失败userId:", player.User.UserId, " needCharm:", goods.Price)
			res.Code = pb.MsgExchangeScoreGoodsRes_LACK_SCORE.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		domainIosActive.GetIosActiveManager().UpdateItemCount(itemId)

	} else {
		goods, ok := domainPrize.GetExchangeGoodsManager().GetExchangeGoods(itemId)
		Goods = goods;
		glog.Infof("goods.TotalCount=%d",Goods.TotalCount);
		glog.Infof("goods.RemainderCount=%d",Goods.RemainderCount);
		
		if !ok {
			glog.Error("找不到兑换配置itemId:", itemId)
			res.Code = pb.MsgExchangeScoreGoodsRes_FAILED.Enum()
			res.Reason = proto.String("找不到兑换配置的物品");
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		info := domainUser.GetUserFortuneManager().GetCharmExchangeInfo(player.User.UserId, itemId)
		if info >= goods.MaxCount && goods.MaxCount != -1 {
			glog.Error("兑换次数到达上限userId:", player.User.UserId, "|count=", info)
			res.Code = pb.MsgExchangeScoreGoodsRes_FAILED.Enum()
			res.Reason = proto.String("兑换次数到达上限");
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		if !domainUser.GetUserFortuneManager().ConsumeCharm(player.User.UserId, goods.NeedScore) {
			glog.Error("消耗魅力失败userId:", player.User.UserId, " needCharm:", goods.NeedScore)
			res.Code = pb.MsgExchangeScoreGoodsRes_LACK_SCORE.Enum()
			res.Reason = proto.String("消耗魅力失败");
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
		
		if(goods.RemainderCount <= 0){
			glog.Error("商品剩余个数为0:", player.User.UserId, " |count=", goods.RemainderCount)
			res.Code = pb.MsgExchangeScoreGoodsRes_FAILED.Enum()
			res.Reason = proto.String("商品剩余个数为0");
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
	}

	domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)
	//util.MongoLog_CharmExchange(goods.NeedScore)
	//util.MongoLog_CharmPool(-goods.NeedScore)

	/*if msg.GetItemId() == 9 {
		// 四倍卡
		prizeMail := &domainPrize.PrizeMail{}
		prizeMail.UserId = player.User.UserId
		prizeMail.MailId = bson.NewObjectId().Hex()
		prizeMail.Content = "魅力兑换"
		prizeMail.ItemType = int(pb.MagicItemType_FOURFOLD_GOLD)
		prizeMail.ItemCount = 1

		player.PrizeMails.AddPrizeMail(prizeMail)

		mailMsg := &pb.MsgGetPrizeMailListRes{}
		mailMsg.Mails = append(mailMsg.Mails, prizeMail.BuildMessage())

		player.SendToClient(int32(pb.MessageId_GET_PRIZE_MAIL_LIST), mailMsg)

		res.Code = pb.MsgExchangeScoreGoodsRes_OK.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}
	if msg.GetItemId() == 10 {
		// 禁比卡
		prizeMail := &domainPrize.PrizeMail{}
		prizeMail.UserId = player.User.UserId
		prizeMail.MailId = bson.NewObjectId().Hex()
		prizeMail.Content = "魅力兑换"
		prizeMail.ItemType = int(pb.MagicItemType_PROHIBIT_COMPARE)
		prizeMail.ItemCount = 1

		player.PrizeMails.AddPrizeMail(prizeMail)

		mailMsg := &pb.MsgGetPrizeMailListRes{}
		mailMsg.Mails = append(mailMsg.Mails, prizeMail.BuildMessage())

		player.SendToClient(int32(pb.MessageId_GET_PRIZE_MAIL_LIST), mailMsg)

		res.Code = pb.MsgExchangeScoreGoodsRes_OK.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}
	if msg.GetItemId() == 11 {
		// 换牌卡
		prizeMail := &domainPrize.PrizeMail{}
		prizeMail.UserId = player.User.UserId
		prizeMail.MailId = bson.NewObjectId().Hex()
		prizeMail.Content = "魅力兑换"
		prizeMail.ItemType = int(pb.MagicItemType_REPLACE_CARD)
		prizeMail.ItemCount = 1

		player.PrizeMails.AddPrizeMail(prizeMail)

		mailMsg := &pb.MsgGetPrizeMailListRes{}
		mailMsg.Mails = append(mailMsg.Mails, prizeMail.BuildMessage())

		player.SendToClient(int32(pb.MessageId_GET_PRIZE_MAIL_LIST), mailMsg)

		res.Code = pb.MsgExchangeScoreGoodsRes_OK.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}*/

	l := &domainPrize.ExchangeGoodsLog{}
	l.UserId = player.User.UserId
	l.UserName = player.User.UserName
	l.ItemId = itemId
	l.IsShipped = false
	l.ShippingAddressName = player.User.ShippingAddressName
	l.ShippingAddressPhone = player.User.ShippingAddressPhone
	l.ShippingAddressAddress = player.User.ShippingAddressAddress
	l.ShippingAddressZipCode = player.User.ShippingAddressZipCode
	l.ShippingAddressQQ = player.User.ShippingAddressQQ

	l.Time = time.Now()
	err = domainPrize.SaveExchangeGoodsLog(l)
	if err != nil {
		res.Code = pb.MsgExchangeScoreGoodsRes_FAILED.Enum()
		glog.Error("保存物品兑换记录失败l:", l)
		res.Reason = proto.String("保存物品兑换记录失败");
		return server.BuildClientMsg(m.GetMsgId(), res)
	} else {
		res.Code = pb.MsgExchangeScoreGoodsRes_OK.Enum()
		res.GoodsId = proto.Int(Goods.GoodsId)
		res.RemainderCount = proto.Int(Goods.RemainderCount-1)
		
		//更新剩余商品个数到数据库
		glog.Infof("goods.goodsid=%d",Goods.GoodsId);
		glog.Infof("goods.RemainderCount=%d",Goods.RemainderCount);
		ok := domainPrize.GetExchangeGoodsManager().SetRemainderToDB(Goods.GoodsId,Goods.RemainderCount-1)
		if(ok != nil){
			res.Code = pb.MsgExchangeScoreGoodsRes_FAILED.Enum()
			glog.Error("更新物品剩余数量到数据库失败:", ok)
			res.Reason = proto.String("更新物品剩余数量到数据库失败");
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
		//更新剩余商品个数到内存
		Goods.RemainderCount = Goods.RemainderCount-1;
		domainPrize.GetExchangeGoodsManager().SetRemainderToMem(Goods.GoodsId,Goods);
	}

	domainUser.GetUserFortuneManager().UpdateCharmExchangeInfo(player.User.UserId, itemId, 1)

	if channelId == 183 {
		content := "玩家Id：" + l.UserId + "， 玩家登陆名：" + l.UserName + "， 商品名称：" +
			goodsName + "， 收件人名人：" + l.ShippingAddressName + "，收件人电话：" +
			l.ShippingAddressPhone + "，收件人地址：" + l.ShippingAddressAddress + "，QQ:" + l.ShippingAddressQQ
		bSend := SendMail("扑克三刀流玩家兑换商品", content)
		if !bSend {
			glog.Info("发送兑换物品邮件失败:", content)
		} else {
			glog.Info("发送兑换物品邮件成功:", content)
		}
	}

	return server.BuildClientMsg(m.GetMsgId(), res)
}
