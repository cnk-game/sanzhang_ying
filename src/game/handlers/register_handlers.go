package handlers

import (
	"game/handlers/active"
	"game/handlers/admin"
	"game/handlers/cdkey"
	"game/handlers/chat"
	handlecofig "game/handlers/config"
	"game/handlers/exchange"
	"game/handlers/fish"
	"game/handlers/game"
	"game/handlers/littlegame"
	"game/handlers/newusertask"
	"game/handlers/pay"
	"game/handlers/prize"
	"game/handlers/rankingList"
	"game/handlers/report"
	"game/handlers/safeBox"
	"game/handlers/slots"
	"game/handlers/stats"
	"game/handlers/user"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"pb"
)

func registerHandlers(r *mux.Router) {
	glog.V(2).Info("register handlers")
	registerHttpHandlers(r)

	registerUserHandlers()
}

func registerHttpHandlers(r *mux.Router) {
	r.HandleFunc("/lenovo_pay", pay.LenovoPayHandler)
	r.HandleFunc("/iapp_pay", pay.IAppPayHandler)
	r.HandleFunc("/qf_pay", pay.QfPayHandler)
	r.HandleFunc("/ips_pay", pay.IPSPayHandler)
	r.HandleFunc("/caohua_pay", pay.ChaohuaPayHandler)
	r.HandleFunc("/leshi_pay", pay.LeshiPayHandler)
	r.HandleFunc("/xunlei_pay", pay.XunLeshiPayHandler)
	r.HandleFunc("/hm_pay", pay.HmPayHandler)
	r.HandleFunc("/setCardConfig", admin.SetCardConfigDataHandler)
	r.HandleFunc("/queryUserById", admin.QueryUserByIdHandler)
	r.HandleFunc("/queryUserByName", admin.QueryUserByNameHandler)
	r.HandleFunc("/setUserFortune", admin.SetUserFortuneHandler)
	r.HandleFunc("/lockUser", admin.LockUserHandler)
	r.HandleFunc("/sendPrizeMail", admin.SendPrizeMailHandler)
	r.HandleFunc("/sendSysMsg", admin.SendSystemMsgHandler)
	r.HandleFunc("/online", admin.GetOnlineCountHandler)
	r.HandleFunc("/onlineByType", admin.GetOnlineTypeHandler)
	r.HandleFunc("/goldLimitUserCount", admin.GetUserCountGoldLimitHandler)
	r.HandleFunc("/getAllGold", admin.GetAllGoldHandler)
	r.HandleFunc("/setCurVersion", admin.SetCurVersionHandler)
	r.HandleFunc("/jinli_pay", pay.JinliPayHandler)
	r.HandleFunc("/kupai_pay", pay.KupaiPayHandler)
	r.HandleFunc("/lianxiang_pay", pay.LianxiangPayHandler)
	r.HandleFunc("/Sanxing_pay", pay.SanxingPayHandler)
	r.HandleFunc("/meizu_pay", pay.MeizuPayHandler)
	r.HandleFunc("/51_pay", pay.Ios51PayHandler)
	r.HandleFunc("/ligame_change_probably", littlegame.ChangeProbablyHandler)

	// add by wangsq start -- fish

	r.HandleFunc("/fish/shoplist", fish.FishShopListHandler)
	r.HandleFunc("/fish/userinfo", fish.FishUserInfoHandler)
	r.HandleFunc("/fish/buyfish", fish.FishBuyHandler)
	r.HandleFunc("/fish/harvestfish", fish.FishHarvestHandler)
	r.HandleFunc("/fish/giftfish", fish.FishGiftHandler)
	r.HandleFunc("/fish/updategold", fish.FishUpdateGoldHandler)
	r.HandleFunc("/fish/getGiftLog", fish.FishGetGiftLogHandler)
	r.HandleFunc("/log/getSysTipLog", admin.GetSystemTipLogHandler)
	r.HandleFunc("/log/getSlotPoolLog", admin.GetSlotPoolLogHandler)
	r.HandleFunc("/log/getGameFeeLog", admin.GetGameFeeLogHandler)
	r.HandleFunc("/log/getCharmLog", admin.GetCharmLogHandler)
	// add by wangsq end

	//
	r.HandleFunc("/admin/user", admin.ChangeShowUserHandler)
	r.HandleFunc("/admin/ipsChangePayType", admin.ChangeIpsPayTypeHandler)
	r.HandleFunc("/admin/add_coins", admin.PayActiveAddCoinsHandler)
	r.HandleFunc("/admin/get_pay", admin.GetRechargeCountHandler)
}

func registerUserHandlers() {
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_LOGIN))
	GetMsgRegistry().RegisterUnLoginMsg(int32(pb.MessageId_RE_CONNECT))

	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_LOGIN), user.LoginHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_RE_CONNECT), user.ReConnectHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_USER_INFO), user.GetUserInfoHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ENTER_POKER_DESK), game.EnterGameHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_KICK_PLAYER), game.KickPlayerHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_VERIFY_CODE), user.GetVerifyCodeHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_REGISTER), user.RegisterHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_RESET_PWD_SAFEBOX), safeBox.ResetPwdSafeBoxHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_RESET_PWD_USER), user.ResetPwdUserHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_GOLD_SAFEBOX), safeBox.GetSafeBoxHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_LOGIN_SAFEBOX), safeBox.LoginSafeBoxHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_SAVE_GOLD_SAFEBOX), safeBox.SaveSafeBoxHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GIFT_GOLD_SAFEBOX), safeBox.GiftSafeBoxHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_CHANGE_PWD_SAFEBOX), safeBox.ChangePwdSafeBoxHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_PAY_TOKEN), pay.GetPayTokenHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_OPEN_FISH_TANK), game.OpenFishTankHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_LEAVE_POKER_DESK), game.LeaveGameHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_PLAYER_OPERATE_CARDS), game.OpCardHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_UPDATE_USER_INFO), user.UpdateUserInfoHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_RANKING_LIST), rankingList.GetRankingListHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_EXCHANGE_GOLD), user.ExchangeGoldHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_CHAT), chat.ChatMsgHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_SHOP_LOG), user.GetShopLogsHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_EXCHANGE_GAME_GOODS), user.ExchangeGameGoodsHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_PLAY_LUCK_WHEEL), slots.PlayLuckyWheelHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_USE_MAGIC_ITEM), user.UseMagicItemHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GAIN_ONLINE_PRIZE), prize.GainOnlinePrizeHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GAIN_VIP_PRIZE), prize.GainVipPrizeHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GAIN_TASK_PRIZE), prize.GainTaskPrizeHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.ServerMsgId_MQ_MATCH_RESULT), game.MatchResultHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_UPDATE_GOLD), user.UpdateGoldHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_JOIN_WAIT_QUEUE), game.JoinWaitQueueHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_LEAVE_WAIT_QUEUE), game.LeaveWaitQueueHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_LOOKUP_USER_BET), game.LookupBetGoldHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_SLOT_MACHINES_PLAY), slots.PlaySlotMachineHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_SLOT_MACHINES_REPLACE), slots.ReplaceSlotMachineCardHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_SLOT_MACHINES_PRIZE), slots.GainSlotMachinePrizeHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_BIND_PRIZE_ADDRESS), prize.BindPrizeAddressHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_EXCHANGE_SCORE_GOODS), prize.ExchangeGoodsByScoreHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_MATCH_RECORD), user.GetMatchRecordHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_SUBSIDY_PRIZE), prize.SubsidyPrizeHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_REWARD_IN_GAME), game.RewardInGameHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_EXCHANGE_CODE_PRIZE), cdkey.ExchangeCDKeyHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_PRIZE_ADDRESS), user.GetShippingAddressHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_ONLINE_STATUS), stats.GetOnlineStatusHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ROBOT_SET_GOLD), user.RobotSetGoldHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_PRIZE_MAIL_LIST), prize.GetPrizeMailsHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GAIN_MAIL_PRIZE), prize.GainMailPrizeHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.ServerMsgId_MQ_PRIZE_MAIL), prize.ServerPrizeMailHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_APP_ENTER_BACKGROUND), game.AppEnterBackgroundHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_APP_ENTER_FOREGROUND), game.AppEnterForegroundHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_REPORT_USER_HEAD_ILLEGAL), report.ReportUserHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_SIGN_IN), prize.SignInHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_SIGN_IN_RECORD), prize.SignInRecordHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_RECHARGE_INFO), user.GetRechargeInfoHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.ServerMsgId_MQ_LOCK_USER), user.LockUserHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.ServerMsgId_MQ_UPDATE_RECHAGE_DIAMOND), user.UpdateRechargeDiamondHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_BUY_DAILY_GIFT_BAG), prize.BuyDailyGiftBagHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_ACTIVE_STATUS), active.GetActiveStatusHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ACTIVE_CONTENT), active.GetActiveContentHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_CAN_USE_THIRDPARTY_PAY), pay.GetPayChangeHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_EXCHANGE_GOODS), prize.GetExchangeGoodsHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_ONLINE_CONFIG), user.GetIosOnlineConfigHandler)
	//wjs activdata
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_ACTIVE_DATA), active.GetActiveDatasHandler)
	//wjs exchnge goods
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_GOOD_LIST), exchange.GetGoodListHandler)
	//wjs macht_config
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_MATCHCONFIG_LIST), handlecofig.GetMatchConfigListHandler)
	//wjs get_active_token
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_ACTIVE_TOKEN), active.GetActiveTokenHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_NEWBETASK_LIST), newUserTask.GetTaskHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_NEWBETASK_REWARD), newUserTask.GetPrizeHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_EXCHANGE_NEWBETASK_HF), newUserTask.GetHuafeiHandler)
	//wjs
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_GET_LITTLE_GAME_CONFIG), littlegame.GetLigameConfHandler)
	GetMsgRegistry().RegisterMsg(int32(pb.MessageId_LITTLE_GAME_REQ), littlegame.LigameChipHandler)
}
