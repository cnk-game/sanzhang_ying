package active

import (
	"code.google.com/p/goprotobuf/proto"
	domainActive "game/domain/iosActive"
	"game/server"
	"pb"
)

func GetActiveContentHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	res := &pb.MsgGetActiveContentRes{}
	content, beginTime, endTime := domainActive.GetUserIosActiveManager().GetIosActiveContent()
	res.Content = proto.String(content)
	res.BeginTime = proto.String(beginTime)
	res.EndTime = proto.String(endTime)
	//res.Content = proto.String("第一届三刀流赌王选拔赛\n1.时间2016/01/15中午12点开始至2016/02/21晚上12点截止。\n2.期间于牌桌上盈利最多之前五名角色获得以下奖励：\n      第一名：美金500元\n      第二名：美金400元\n      第三名：美金300元\n      第四名：美金200元\n      第五名：美金100元\n3.参与竞赛人员需于2015/02/15晚上12点以前完成个人资料登陆，否则本次竞赛不计入排名。\n4.排名结果会于2/23中午12点正式公布。\n5.活动结束后得奖人员，本公司之客服会于2/29之前联系得奖人员喔。\n6.后续尚有胜场王跟魅力王之竞赛喔。")
	//res.BeginTime = proto.String("2016-01-15 12:00:00")
	//res.EndTime = proto.String("2016-02-22 00:00:00")

	return server.BuildClientMsg(m.GetMsgId(), res)
}
