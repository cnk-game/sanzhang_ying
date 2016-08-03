package user

import (
	"code.google.com/p/goprotobuf/proto"
	"flag"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"github.com/larspensjo/config"
	"pb"
	"util"
)

var (
	configFile = flag.String("configfile", "./config/config.ini", "General configuration file")
)

const (
	CHANNEL_IOS_APPSTORE = "179"
	CHANNEL_IOS_XY       = "178"
)

func readConfig() string {
	var TOPIC = make(map[string]string)
	//set config file std
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		glog.Info("Fail to find", *configFile, err)
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfg.HasSection("topicArr") {
		section, err := cfg.SectionOptions("topicArr")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("topicArr", v)
				if err == nil {
					TOPIC[v] = options
				}
			}
		}
	}
	//Initialized topic from the configuration END

	//fmt.Println(TOPIC)
	//fmt.Println(TOPIC["userId"])
	return TOPIC["userId"]
}

func GetUserInfoHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	//userConfig := readConfig()
	//glog.Info("userConfig ", userConfig)
	player := domainUser.GetPlayer(sess.Data) //game_player.go

	f, _ := domainUser.GetUserFortuneManager().GetUserFortune(player.User.UserId)
	player.UserTasks.AccomplishTask(util.TaskAccomplishType_GOLD, f.Gold, player.SendToClientFunc)

	if len(f.VipTaskStates) == 0 {
		domainUser.GetUserFortuneManager().InitVipTaskState(player.User.UserId)
	}

	res := &pb.MsgGetUserInfoRes{}
	res.UserInfo = player.User.BuildMessage(player.MatchRecord.BuildMessage())
	res.UserInfo.IsRobot = proto.Bool(false)

	if f.DoubleCardCount > 0 {
		itemMsg := &pb.UserMagicItemDef{}
		itemMsg.ItemType = pb.MagicItemType_FOURFOLD_GOLD.Enum()
		itemMsg.Count = proto.Int(f.DoubleCardCount)
		res.UserInfo.ItemList = append(res.UserInfo.ItemList, itemMsg)
	}

	if f.ForbidCardCount > 0 {
		itemMsg := &pb.UserMagicItemDef{}
		itemMsg.ItemType = pb.MagicItemType_PROHIBIT_COMPARE.Enum()
		itemMsg.Count = proto.Int(f.ForbidCardCount)
		res.UserInfo.ItemList = append(res.UserInfo.ItemList, itemMsg)
	}

	if f.ChangeCardCount > 0 {
		itemMsg := &pb.UserMagicItemDef{}
		itemMsg.ItemType = pb.MagicItemType_REPLACE_CARD.Enum()
		itemMsg.Count = proto.Int(f.ChangeCardCount)
		res.UserInfo.ItemList = append(res.UserInfo.ItemList, itemMsg)
	}

	res.UserInfo.PrizeTaskList = player.UserTasks.BuildMessage2()
	res.UserInfo.OnlinePrizeList = player.OnlinePrizeGainRecords.BuildMessage()

	msg := &pb.ShareInfoDef{}

	msg.Platforms = []int32{}
	//xunlei
	if player.User.ChannelId == "186" {

		msg.Platforms = append(msg.Platforms, 1)
		msg.Platforms = append(msg.Platforms, 2)
		msg.Platforms = append(msg.Platforms, 3)
		msg.Platforms = append(msg.Platforms, 4)
		msg.WxAppid = proto.String("wx6fcd4b1f935474f2")
		msg.WxAppSceret = proto.String("72c7cfb105ef450130c36ea166b6bcaa")
	} else if player.User.ChannelId == "173" {
		//leishi
		msg.Platforms = append(msg.Platforms, 1)
		msg.Platforms = append(msg.Platforms, 2)
		msg.Platforms = append(msg.Platforms, 3)
		msg.Platforms = append(msg.Platforms, 4)
		msg.WxAppid = proto.String("wx3f1caad21f70140c")
		msg.WxAppSceret = proto.String("1e68a49f9b4f192c2dbcdc99599dc297")
	} else if player.User.ChannelId == "212" {
		//jinli

		msg.Platforms = append(msg.Platforms, 3)
		msg.Platforms = append(msg.Platforms, 4)
		msg.WxAppid = proto.String("")
		msg.WxAppSceret = proto.String("")
	} else {
		msg.Platforms = append(msg.Platforms, 1)
		msg.Platforms = append(msg.Platforms, 2)
		msg.Platforms = append(msg.Platforms, 3)
		msg.Platforms = append(msg.Platforms, 4)
		msg.WxAppid = proto.String("wx05e5e26a69f5d78e")
		msg.WxAppSceret = proto.String("d4624c36b6795d1d99dcf0547af5443d")
	}
	msg.QqAppid = proto.String("1104978566")
	msg.QqAppKey = proto.String("M5QD1vF2HJprOyeq")
	msg.ShareIconResUrl = proto.String("http://res.dapai1.com/Uploads/Picture/app/three.png")

	if player.User.ChannelId == "178" ||
		player.User.ChannelId == "183" {

		msg.ShareURL = proto.String("http://www.puke111.com/ios/xy/download.html")
	} else if player.User.ChannelId == JINLI_CHANNEL {
		msg.ShareURL = proto.String("http://www.puke111.com/jinli/download.html")
	} else if player.User.ChannelId == KUPAI_CHANNEL {
		msg.ShareURL = proto.String("http://www.puke111.com/kupai/download.html")
	} else if player.User.ChannelId == IOS51_CHANNEL {
		msg.ShareURL = proto.String("http://www.puke111.com/ios/51/download.html")
	} else if player.User.ChannelId == CHANNEL_IOS_APPSTORE {
		msg.ShareURL = proto.String("http://www.puke111.com/ios/appstore/download.html")
	} else if player.User.ChannelId == HAIMA_IOS_CHANNEL {
		msg.ShareURL = proto.String("http://www.puke111.com/ios/haima/download.html")
	} else {
		msg.ShareURL = proto.String("http://www.puke111.com/download.html")
	}

	res.UserInfo.SharePlatform = msg

	return server.BuildClientMsg(m.GetMsgId(), res)
}
