var select_user_id = ""
var loader = null

function initMenu() {
    $("#game_mgr_menu").addClass("active")
    $("#game_mgr_sub_menu").show()
    $("#player_mgr").addClass("active")

}

function initTable(tableTitle, tableHeadArr) {
    $("#data_table_title").text(tableTitle)

    var dom_thead_tr = $('<tr></tr>');
    var dom_tbody_tr = $('<tr></tr>');
    for (var i = 0; i < tableHeadArr.length; i++) {
        dom_thead_tr.append($("<th>" + tableHeadArr[i] + "</th>"));
        dom_tbody_tr.append($("<td></td>"));
    }
    var dom_thead = $('<thead></thead>');
    dom_thead.append(dom_thead_tr);
    var dom_tbody = $('<tbody></tbody>');
    dom_tbody.append(dom_tbody_tr);
    $('#data_table').append(dom_thead).append(dom_tbody)
}

function clear_table() {
    var rownum=$("#data_table tr").length - 2;
    for (var i = 0; i < rownum; i++) {
        $("#data_table tr:eq(2)").remove();
    }
}

// query
function query_user(userId, nickname) {
    loader = layer.load('查询中…')
    $.getJSON("/game_mgr/query_user", {userId:userId, nickname:nickname}, function (data) {
        layer.close(loader)
        if (data == null) return
        clear_table()
        $("#data_table tr:gt(0):eq(0)").hide()

        $.each(data, function(i, item) {
            var row = "<tr>"
            row += "<td>" + item.userId + "</td>"
            row += "<td>" + item.username + "</td>"
            row += "<td>" + item.nickname + "</td>"
            row += "<td>" + item.level + "</td>"
            row += "<td>" + item.vipLevel + "</td>"
            row += "<td>" + item.gold + "</td>"
            row += "<td>" + item.diamond + "</td>"
            row += "<td>" + getLocalTime(item.createTime) + "</td>"
            row += "<td>" + "<button class=\"btn btn-primary btn-xs\" onclick=\"select_player('"+item.userId+"')\"><i class=\"fa fa-pencil\"></i></button>" + "</td>"
            row += "</tr>"
            $("#data_table tr:last").after(row)
        })
        $('#player_list_panel').show()
    })
}

// 按帐号查询
function query_by_userId() {
    console.log("query_by_username")
    query_user($('#userId').val(), "")
    $("#userId").val("")
}

// 按昵称查询
function query_by_nickname() {
    console.log("query_by_nickname")
    query_user("", $('#nickname').val())
    $("#nickname").val("")
}

// 选择一个查询的玩家
function select_player(select_btn) {
    select_user_id = select_btn;
    console.log("select edit player => " + select_user_id)
    $('#player_info_editor_title').text("玩家-" + select_btn + "-个人信息")
    $('#player_info_editor_panel').show()
}

// 玩家增加金币
function inc_gold() {
    console.log("inc_gold")
    gold = $('#incGoldCount').val()
    $("#incGoldCount").val("")
    loader = layer.load('请稍候…')
    $.getJSON("/game_mgr/set_user_fortune", {userId:select_user_id, diamond:0, gold:gold}, function (data) {
        layer.close(loader)
        if (data == null) return
        if (data == "succeed") {
            console.log("succeed")
        }
    })
}

// 玩家增加钻石
function inc_diamond() {
    console.log("inc_diamond")
    diamond = $('#incDiamondCount').val()
    $("#incDiamondCount").val("")
    loader = layer.load('请稍候…')
    $.getJSON("/game_mgr/set_user_fortune", {userId:select_user_id, diamond:diamond, gold:0}, function (data) {
        layer.close(loader)
        if (data == null) return
        if (data == "succeed") {
            console.log("succeed")
        }
    })
}

// 玩家邮件通知
function send_mail() {
    console.log("send_mail")
    content = $('#mailContent').val()
    goldCount = $('#goldCount').val()
    diamondCount = $('#diamondCount').val()
    $("#mailContent").val("")
    $("#goldCount").val("")
    $("#diamondCount").val("")
    loader = layer.load('请稍候…')
    var date = {userId:select_user_id, content:content, gold:goldCount, diamond:diamondCount, itemType:0, itemCount:0}
    $.getJSON("/game_mgr/send_user_prize_mail", date, function (data) {
        layer.close(loader)
        if (data == null) {
            return
        }
        if (data == "succeed") {
            console.log("succeed")
        }
    });
}

// 冻结帐号
function freeze_player() {
    loader = layer.load('请稍候…')
    $.getJSON("/game_mgr/lock_user", {userId:select_user_id, isLock:"true"}, function (data) {
        layer.close(loader)
        if (data == null) return
        if (data == "succeed") {
            console.log("succeed")
        }
    })
}

// 解冻帐号
function un_freeze_player() {
    loader = layer.load('请稍候…')
    $.getJSON("/game_mgr/lock_user", {userId:select_user_id, isLock:"false"}, function (data) {
        layer.close(loader)
        if (data == null) return
        if (data == "succeed") {
            console.log("succeed")
        }
    })
}

$(document).ready(function () {
    initMenu()

    initTable("查询结果", ["用户ID", "帐号", "昵称", "等级", "VIP", "金币", "钻石", "创建时间", "选择玩家"])

    $('#player_info_editor_panel').hide()
    $('#player_list_panel').hide()

    $('#username_query').click(function () {
        query_by_userId()
    })

    $('#nickname_query').click(function () {
        query_by_nickname()
    })

    $('#inc_gold').click(inc_gold)
    $('#inc_diamond').click(inc_diamond)
    $('#send_mail').click(send_mail)
    $('#freeze_player').click(freeze_player)
    $('#un_freeze_player').click(un_freeze_player)
});
function getLocalTime(nS) {     
    return new Date(parseInt(nS) * 1000).toLocaleString().substr(0,10)}