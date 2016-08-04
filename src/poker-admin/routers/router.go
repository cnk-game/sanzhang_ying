package routers

import (
	"poker-admin/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.LoginController{})
    beego.Router("/login", &controllers.LoginController{})
    beego.Router("/logout", &controllers.LogoutController{})
    beego.Router("/setpwd", &controllers.PwdMgrController{})
    beego.Router("/savepwd", &controllers.SavePwdController{})

    beego.Router("/main", &controllers.MainController{})
    beego.Router("/main_today_query", &controllers.MainTodayQueryController{})
    beego.Router("/main_now_query", &controllers.MainNowQueryController{})
    beego.Router("/today_online_status", &controllers.TodayOnlineStatusController{})

    beego.Router("/notify_mgr", &controllers.NotifyMgrController{})
    beego.Router("/player_mgr", &controllers.PlayerMgrController{})

    beego.Router("/host_mgr",&controllers.HostMgrController{})
    beego.Router("/slot_mgr",&controllers.SlotMgrController{})
    beego.Router("/fee_mgr",&controllers.FeeMgrController{})
    beego.Router("/charm_mgr",&controllers.CharmMgrController{})

    beego.Router("/operate_host_query",&controllers.HostQueryController{})
    beego.Router("/operate_slot_query",&controllers.SlotQueryController{})
    beego.Router("/operate_fee_query",&controllers.FeeQueryController{})
    beego.Router("/operate_charm_query",&controllers.CharmQueryController{})

    beego.Router("/synthesize_info", &controllers.SynthesizeInfoController{})
    beego.Router("/synthesize_query", &controllers.SynthesizeQueryController{})

    beego.Router("/new_player_info", &controllers.NewPlayerInfoController{})
    beego.Router("/new_player_query", &controllers.NewPlayerQueryController{})

    beego.Router("/pay_info", &controllers.PayInfoController{})
    beego.Router("/pay_query", &controllers.PayQueryController{})

    beego.Router("/user_mgr", &controllers.UserMgrController{})
    beego.Router("/add_user", &controllers.AddUserController{})
    beego.Router("/remove_user", &controllers.RemoveUserController{})
    beego.Router("/query_user_list", &controllers.QueryUserListController{})
    
    beego.Router("/config_mgr", &controllers.ConfigMgrController{})
    beego.Router("/query_config", &controllers.QueryConfigController{})
    beego.Router("/update_config", &controllers.UpdateConfigController{})
    beego.Router("/query_admin", &controllers.QueryAdminController{})
    
    beego.Router("/game_mgr/get_match_config", &controllers.GetMatchConfigController{})
    beego.Router("/game_mgr/save_match_config", &controllers.SaveMatchConfigController{})

    beego.Router("/game_mgr/get_prize_version", &controllers.GetPrizeVersionController{})
    beego.Router("/game_mgr/save_prize_version", &controllers.SavePrizeVersionController{})

    beego.Router("/game_mgr/set_user_fortune", &controllers.SetUserFortuneController{})
    beego.Router("/game_mgr/lock_user", &controllers.LockUserController{})
    beego.Router("/game_mgr/send_system_message", &controllers.SendSystemMessageController{})
    beego.Router("/game_mgr/send_user_prize_mail", &controllers.SendUserPrizeMailController{})
    beego.Router("/game_mgr/query_user", &controllers.QueryUserController{})
}