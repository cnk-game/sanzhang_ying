var loader = null

function initMenu() {
    $("#manager_menu").addClass("active")
    $("#manager_sub_menu").show()
    $("#config_mgr").addClass("active")


}

function query_config() {
    loader = layer.load('请稍候…')
    $.getJSON("/query_config", function (data) {
        layer.close(loader)
        if (data == null) return
        $("#main-content").show()
        $("#onlineSecond").val(data.VaildNewPlayerOnlineSecond)
        $("#matchCount").val(data.VaildNewPlayerMatchCount)
        console.log(data.VaildNewPlayerOnlineSecond, data.VaildNewPlayerMatchCount)
        query_match_config()
    });
}

function update_config() {
    var vaild_new_player_online_second = $('#onlineSecond').val()
    var vaild_new_player_match_count = $('#matchCount').val()
    loader = layer.load('请稍候…')
    $.getJSON("/update_config", {vaild_new_player_online_second:vaild_new_player_online_second, vaild_new_player_match_count:vaild_new_player_match_count}, function (data) {
        layer.close(loader)
        if (data == null)
            return
        if (0 == data) {
            alert("修改成功")
        } else {
            alert("修改失败")
        }
    });
}

function query_match_config() {
    loader = layer.load('请稍候…')
    $.getJSON("/game_mgr/get_match_config", function (data) {
        layer.close(loader)
        if (data == null) 
            return
        $("#common-level1-Single").val(data.CommonLevel1Single)
        $("#common-level1-Double").val(data.CommonLevel1Double)
        $("#common-level1-ShunZi").val(data.CommonLevel1ShunZi)
        $("#common-level1-JinHua").val(data.CommonLevel1JinHua)
        $("#common-level1-ShunJin").val(data.CommonLevel1ShunJin)
        $("#common-level1-BaoZi").val(data.CommonLevel1BaoZi)
        $("#common-level1-WinGold").val(data.CommonLevel1WinGold)
        $("#common-level1-WinRateHigh").val(data.CommonLevel1WinRateHigh)
        $("#common-level1-LoseGold").val(data.CommonLevel1LoseGold)
        $("#common-level1-WinRateLow").val(data.CommonLevel1WinRateLow)

        $("#common-level2-Single").val(data.CommonLevel2Single)
        $("#common-level2-Double").val(data.CommonLevel2Double)
        $("#common-level2-ShunZi").val(data.CommonLevel2ShunZi)
        $("#common-level2-JinHua").val(data.CommonLevel2JinHua)
        $("#common-level2-ShunJin").val(data.CommonLevel2ShunJin)
        $("#common-level2-BaoZi").val(data.CommonLevel2BaoZi)
        $("#common-level2-WinGold").val(data.CommonLevel2WinGold)
        $("#common-level2-WinRateHigh").val(data.CommonLevel2WinRateHigh)
        $("#common-level2-LoseGold").val(data.CommonLevel2LoseGold)
        $("#common-level2-WinRateLow").val(data.CommonLevel2WinRateLow)

        $("#common-level3-Single").val(data.CommonLevel3Single)
        $("#common-level3-Double").val(data.CommonLevel3Double)
        $("#common-level3-ShunZi").val(data.CommonLevel3ShunZi)
        $("#common-level3-JinHua").val(data.CommonLevel3JinHua)
        $("#common-level3-ShunJin").val(data.CommonLevel3ShunJin)
        $("#common-level3-BaoZi").val(data.CommonLevel3BaoZi)
        $("#common-level3-WinGold").val(data.CommonLevel3WinGold)
        $("#common-level3-WinRateHigh").val(data.CommonLevel3WinRateHigh)
        $("#common-level3-LoseGold").val(data.CommonLevel3LoseGold)
        $("#common-level3-WinRateLow").val(data.CommonLevel3WinRateLow)

        $("#item-level1-Single").val(data.ItemLevel1Single)
        $("#item-level1-Double").val(data.ItemLevel1Double)
        $("#item-level1-ShunZi").val(data.ItemLevel1ShunZi)
        $("#item-level1-JinHua").val(data.ItemLevel1JinHua)
        $("#item-level1-ShunJin").val(data.ItemLevel1ShunJin)
        $("#item-level1-BaoZi").val(data.ItemLevel1BaoZi)
        $("#item-level1-WinGold").val(data.ItemLevel1WinGold)
        $("#item-level1-WinRateHigh").val(data.ItemLevel1WinRateHigh)
        $("#item-level1-LoseGold").val(data.ItemLevel1LoseGold)
        $("#item-level1-WinRateLow").val(data.ItemLevel1WinRateLow)

        $("#item-level2-Single").val(data.ItemLevel2Single)
        $("#item-level2-Double").val(data.ItemLevel2Double)
        $("#item-level2-ShunZi").val(data.ItemLevel2ShunZi)
        $("#item-level2-JinHua").val(data.ItemLevel2JinHua)
        $("#item-level2-ShunJin").val(data.ItemLevel2ShunJin)
        $("#item-level2-BaoZi").val(data.ItemLevel2BaoZi)
        $("#item-level2-WinGold").val(data.ItemLevel2WinGold)
        $("#item-level2-WinRateHigh").val(data.ItemLevel2WinRateHigh)
        $("#item-level2-LoseGold").val(data.ItemLevel2LoseGold)
        $("#item-level2-WinRateLow").val(data.ItemLevel2WinRateLow)

        $("#item-level3-Single").val(data.ItemLevel3Single)
        $("#item-level3-Double").val(data.ItemLevel3Double)
        $("#item-level3-ShunZi").val(data.ItemLevel3ShunZi)
        $("#item-level3-JinHua").val(data.ItemLevel3JinHua)
        $("#item-level3-ShunJin").val(data.ItemLevel3ShunJin)
        $("#item-level3-BaoZi").val(data.ItemLevel3BaoZi)
        $("#item-level3-WinGold").val(data.ItemLevel3WinGold)
        $("#item-level3-WinRateHigh").val(data.ItemLevel3WinRateHigh)
        $("#item-level3-LoseGold").val(data.ItemLevel3LoseGold)
        $("#item-level3-WinRateLow").val(data.ItemLevel3WinRateLow)

        $("#sng-level1-Single").val(data.SngLevel1Single)
        $("#sng-level1-Double").val(data.SngLevel1Double)
        $("#sng-level1-ShunZi").val(data.SngLevel1ShunZi)
        $("#sng-level1-JinHua").val(data.SngLevel1JinHua)
        $("#sng-level1-ShunJin").val(data.SngLevel1ShunJin)
        $("#sng-level1-BaoZi").val(data.SngLevel1BaoZi)
        $("#sng-level1-WinGold").val(data.SngLevel1WinGold)
        $("#sng-level1-WinRateHigh").val(data.SngLevel1WinRateHigh)
        $("#sng-level1-LoseGold").val(data.SngLevel1LoseGold)
        $("#sng-level1-WinRateLow").val(data.SngLevel1WinRateLow)

        $("#sng-level2-Single").val(data.SngLevel2Single)
        $("#sng-level2-Double").val(data.SngLevel2Double)
        $("#sng-level2-ShunZi").val(data.SngLevel2ShunZi)
        $("#sng-level2-JinHua").val(data.SngLevel2JinHua)
        $("#sng-level2-ShunJin").val(data.SngLevel2ShunJin)
        $("#sng-level2-BaoZi").val(data.SngLevel2BaoZi)
        $("#sng-level2-WinGold").val(data.SngLevel2WinGold)
        $("#sng-level2-WinRateHigh").val(data.SngLevel2WinRateHigh)
        $("#sng-level2-LoseGold").val(data.SngLevel2LoseGold)
        $("#sng-level2-WinRateLow").val(data.SngLevel2WinRateLow)

        $("#sng-level3-Single").val(data.SngLevel3Single)
        $("#sng-level3-Double").val(data.SngLevel3Double)
        $("#sng-level3-ShunZi").val(data.SngLevel3ShunZi)
        $("#sng-level3-JinHua").val(data.SngLevel3JinHua)
        $("#sng-level3-ShunJin").val(data.SngLevel3ShunJin)
        $("#sng-level3-BaoZi").val(data.SngLevel3BaoZi)
        $("#sng-level3-WinGold").val(data.SngLevel3WinGold)
        $("#sng-level3-WinRateHigh").val(data.SngLevel3WinRateHigh)
        $("#sng-level3-LoseGold").val(data.SngLevel3LoseGold)
        $("#sng-level3-WinRateLow").val(data.SngLevel3WinRateLow)

        $("#wan-Single").val(data.WanSingle)
        $("#wan-Double").val(data.WanDouble)
        $("#wan-ShunZi").val(data.WanShunZi)
        $("#wan-JinHua").val(data.WanJinHua)
        $("#wan-ShunJin").val(data.WanShunJin)
        $("#wan-BaoZi").val(data.WanBaoZi)
        $("#wan-WinGold").val(data.WanWinGold)
        $("#wan-WinRateHigh").val(data.WanWinRateHigh)
        $("#wan-LoseGold").val(data.WanLoseGold)
        $("#wan-WinRateLow").val(data.WanWinRateLow)
        query_prize_version()
    });
}

function query_prize_version() {
    loader = layer.load('请稍候…')
    $.getJSON("/game_mgr/get_prize_version", function (data) {
        layer.close(loader)
        if (data == null) 
            return
        $("#prize-version").val(data)
    });
}

function save_match_config(GameType, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow) {
    if (Single < 0 || Double < 0 || ShunZi < 0 || JinHua < 0 || ShunJin < 0 || BaoZi < 0 || WinGold < 0 || WinRateHigh < 0 || LoseGold < 0 || WinRateLow < 0) {
        alert("概率值不可以为负数")
        return
    }
    if ((Single + Double + ShunZi + JinHua + ShunJin + BaoZi) <= 0) {
        alert("总概率不可为0")
        return;
    }
    loader = layer.load('请稍候…')
    $.getJSON("/game_mgr/save_match_config", 
        {GameType:GameType, Single:Single, Double:Double, ShunZi:ShunZi, JinHua:JinHua, ShunJin:ShunJin, BaoZi:BaoZi, WinGold:WinGold, WinRateHigh:WinRateHigh, LoseGold:LoseGold, WinRateLow:WinRateLow}, 
        function (data) {
            layer.close(loader)
            if (data == null)
                return
            if ("succeed" == data) {
                alert("修改成功")
            } else {
                alert("修改失败")
            }
        });
}

function save_common_level1_match_config() {
    var Single = $('#common-level1-Single').val()
    var Double = $('#common-level1-Double').val()
    var ShunZi = $('#common-level1-ShunZi').val()
    var JinHua = $('#common-level1-JinHua').val()
    var ShunJin = $('#common-level1-ShunJin').val()
    var BaoZi = $('#common-level1-BaoZi').val()
    var WinGold = $('#common-level1-WinGold').val()
    var WinRateHigh = $('#common-level1-WinRateHigh').val()
    var LoseGold = $('#common-level1-LoseGold').val()
    var WinRateLow = $('#common-level1-WinRateLow').val()
    save_match_config(1, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
}
function save_common_level2_match_config() {
    var Single = $('#common-level2-Single').val()
    var Double = $('#common-level2-Double').val()
    var ShunZi = $('#common-level2-ShunZi').val()
    var JinHua = $('#common-level2-JinHua').val()
    var ShunJin = $('#common-level2-ShunJin').val()
    var BaoZi = $('#common-level2-BaoZi').val()
    var WinGold = $('#common-level2-WinGold').val()
    var WinRateHigh = $('#common-level2-WinRateHigh').val()
    var LoseGold = $('#common-level2-LoseGold').val()
    var WinRateLow = $('#common-level2-WinRateLow').val()
    save_match_config(2, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
}
function save_common_level3_match_config() {
    var Single = $('#common-level3-Single').val()
    var Double = $('#common-level3-Double').val()
    var ShunZi = $('#common-level3-ShunZi').val()
    var JinHua = $('#common-level3-JinHua').val()
    var ShunJin = $('#common-level3-ShunJin').val()
    var BaoZi = $('#common-level3-BaoZi').val()
    var WinGold = $('#common-level3-WinGold').val()
    var WinRateHigh = $('#common-level3-WinRateHigh').val()
    var LoseGold = $('#common-level3-LoseGold').val()
    var WinRateLow = $('#common-level3-WinRateLow').val()
    save_match_config(3, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
}

function save_item_level1_match_config() {
    var Single = $('#item-level1-Single').val()
    var Double = $('#item-level1-Double').val()
    var ShunZi = $('#item-level1-ShunZi').val()
    var JinHua = $('#item-level1-JinHua').val()
    var ShunJin = $('#item-level1-ShunJin').val()
    var BaoZi = $('#item-level1-BaoZi').val()
    var WinGold = $('#item-level1-WinGold').val()
    var WinRateHigh = $('#item-level1-WinRateHigh').val()
    var LoseGold = $('#item-level1-LoseGold').val()
    var WinRateLow = $('#item-level1-WinRateLow').val()
    save_match_config(11, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
}
function save_item_level2_match_config() {
    var Single = $('#item-level2-Single').val()
    var Double = $('#item-level2-Double').val()
    var ShunZi = $('#item-level2-ShunZi').val()
    var JinHua = $('#item-level2-JinHua').val()
    var ShunJin = $('#item-level2-ShunJin').val()
    var BaoZi = $('#item-level2-BaoZi').val()
    var WinGold = $('#item-level2-WinGold').val()
    var WinRateHigh = $('#item-level2-WinRateHigh').val()
    var LoseGold = $('#item-level2-LoseGold').val()
    var WinRateLow = $('#item-level2-WinRateLow').val()
    save_match_config(12, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
}
function save_item_level3_match_config() {
    var Single = $('#item-level3-Single').val()
    var Double = $('#item-level3-Double').val()
    var ShunZi = $('#item-level3-ShunZi').val()
    var JinHua = $('#item-level3-JinHua').val()
    var ShunJin = $('#item-level3-ShunJin').val()
    var BaoZi = $('#item-level3-BaoZi').val()
    var WinGold = $('#item-level3-WinGold').val()
    var WinRateHigh = $('#item-level3-WinRateHigh').val()
    var LoseGold = $('#item-level3-LoseGold').val()
    var WinRateLow = $('#item-level3-WinRateLow').val()
    save_match_config(13, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
}

function save_sng_level1_match_config() {
    var Single = $('#sng-level1-Single').val()
    var Double = $('#sng-level1-Double').val()
    var ShunZi = $('#sng-level1-ShunZi').val()
    var JinHua = $('#sng-level1-JinHua').val()
    var ShunJin = $('#sng-level1-ShunJin').val()
    var BaoZi = $('#sng-level1-BaoZi').val()
    var WinGold = $('#sng-level1-WinGold').val()
    var WinRateHigh = $('#sng-level1-WinRateHigh').val()
    var LoseGold = $('#sng-level1-LoseGold').val()
    var WinRateLow = $('#sng-level1-WinRateLow').val()
    save_match_config(21, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
}
function save_sng_level2_match_config() {
    var Single = $('#sng-level2-Single').val()
    var Double = $('#sng-level2-Double').val()
    var ShunZi = $('#sng-level2-ShunZi').val()
    var JinHua = $('#sng-level2-JinHua').val()
    var ShunJin = $('#sng-level2-ShunJin').val()
    var BaoZi = $('#sng-level2-BaoZi').val()
    var WinGold = $('#sng-level2-WinGold').val()
    var WinRateHigh = $('#sng-level2-WinRateHigh').val()
    var LoseGold = $('#sng-level2-LoseGold').val()
    var WinRateLow = $('#sng-level2-WinRateLow').val()
    save_match_config(22, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
}
function save_sng_level3_match_config() {
    var Single = $('#sng-level3-Single').val()
    var Double = $('#sng-level3-Double').val()
    var ShunZi = $('#sng-level3-ShunZi').val()
    var JinHua = $('#sng-level3-JinHua').val()
    var ShunJin = $('#sng-level3-ShunJin').val()
    var BaoZi = $('#sng-level3-BaoZi').val()
    var WinGold = $('#sng-level3-WinGold').val()
    var WinRateHigh = $('#sng-level3-WinRateHigh').val()
    var LoseGold = $('#sng-level3-LoseGold').val()
    var WinRateLow = $('#sng-level3-WinRateLow').val()
    save_match_config(23, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
}

function save_wan_match_config() {
    var Single = $('#wan-Single').val()
    var Double = $('#wan-Double').val()
    var ShunZi = $('#wan-ShunZi').val()
    var JinHua = $('#wan-JinHua').val()
    var ShunJin = $('#wan-ShunJin').val()
    var BaoZi = $('#wan-BaoZi').val()
    var WinGold = $('#wan-WinGold').val()
    var WinRateHigh = $('#wan-WinRateHigh').val()
    var LoseGold = $('#wan-LoseGold').val()
    var WinRateLow = $('#wan-WinRateLow').val()
    save_match_config(30, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
}

function set_prize_version() {
    var version = $('#prize-version').val()
    loader = layer.load('请稍候…')
    $.getJSON("/game_mgr/save_prize_version", {version:version}, function (data) {
        layer.close(loader)
        if (data == null)
            return
        if ("succeed" == data) {
            alert("修改成功")
        } else {
            alert("修改失败")
        }
    });
}

$(document).ready(function () {
    
    initMenu()

    query_config()

    $("#main-content").hide()

    $("#onlineSecond").val("0")
    $("#matchCount").val("0")

    $('#update_channel_config').click(update_config)

    $('#save_common_level1_match_config').click(save_common_level1_match_config)
    $('#save_common_level2_match_config').click(save_common_level2_match_config)
    $('#save_common_level3_match_config').click(save_common_level3_match_config)

    $('#save_item_level1_match_config').click(save_item_level1_match_config)
    $('#save_item_level2_match_config').click(save_item_level2_match_config)
    $('#save_item_level3_match_config').click(save_item_level3_match_config)

    $('#save_sng_level1_match_config').click(save_sng_level1_match_config)
    $('#save_sng_level2_match_config').click(save_sng_level2_match_config)
    $('#save_sng_level3_match_config').click(save_sng_level3_match_config)

    $('#save_wan_match_config').click(save_wan_match_config)

    $('#set_prize_version').click(set_prize_version)
});