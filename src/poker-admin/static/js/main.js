var isWaitQueryMain = false
var isWaitQueryOnline = false

function initMenu() {
    $("#tab_console").addClass("active")

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

function main_query() {
    if (isWaitQueryMain) return;
    isWaitQueryMain = true;
    $.getJSON("/main_today_query", function (data) {
        if (data == null) return;
        $("#data_table tr:gt(0):eq(0)").hide();

        $("#today").text("今日综合信息")
        $("#total_player_count").text("总累计用户:" + data.TotalPlayerCount)
        $("#new_player_count").text(data.TodayNewPlayerCount)
        $("#login_player_count").text(data.TodayLoginPlayerCount)
        $("#pay_player_count").text(data.TodayPayPlayerCount)
        $("#pay_count").text(data.TodayPayCount)
        $("#total_pay").text(data.TodayTotalPay)

        if (data.ChannelGameInfoList == null) {
            $("#data_panel").hide()
        } else {
            $.each(data.ChannelGameInfoList, function(i, item) {
                var row = "<tr>"
                row += "<td>" + item.Channel + "</td>"
                row += "<td>" + item.TodayNewPlayerCount + "</td>"
                row += "<td>" + item.TodayLoginPlayerCount + "</td>"
                row += "<td>" + item.TodayPayPlayerCount + "</td>"
                row += "<td>" + item.TodayPayCount + "</td>"
                row += "<td>" + item.TodayTotalPay + "</td>"
                if (0 == item.TodayLoginPlayerCount) {
                    row += "<td>" + (0).toFixed(2) + "%</td>"
                } else {
                    row += "<td>" + ((item.TodayPayPlayerCount / item.TodayLoginPlayerCount) * 100).toFixed(2) + "%</td>"
                }
                
                row += "</tr>"
                $("#data_table tr:last").after(row);
            });
        }
        //console.log(data)
        isWaitQueryMain = false;
    });

    $.getJSON("/main_now_query", function (data) {
        if (data == null) return;
        $.each(data.infos, function(i, item) {
            var index = item.gameType;

            $("#room_" + index).text(item.count)
        });
    });
}

function online_query() {
    if (isWaitQueryOnline) return;
    isWaitQueryOnline = true;
    $.getJSON("/today_online_status", function (data) {
        if (data == null) {
            $("#online_status").hide()
            return
        }
        //console.log(data)
        Morris.Line({
            element: 'hero-graph',
            data: data,
            xkey: 'Period',
            ykeys: ['OnlinePlayerCount'],
            labels: ['在线人数'],
            lineColors:['#4ECDC4']
        });
        isWaitQueryOnline = false;
    });
}


function clear_table() {
    var rownum=$("#data_table tr").length - 2;
    for (var i = 0; i < rownum; i++) {
        $("#data_table tr:eq(2)").remove();
    }
}




$(document).ready(function () {
    initMenu()
    
    initTable("渠道信息", ["渠道ID", "新增用户", "活跃用户", "付费人数", "付费次数", "付费金额", "付费率"])

    main_query()

    online_query()
    
});
