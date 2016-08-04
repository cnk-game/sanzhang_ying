var isWaitQuery = false
var curPageIdx = 0
var totalPageCount = 0

function initMenu() {
    $("#statistics_menu").addClass("active")
    $("#statistics_sub_menu").show()
    $("#new_player_info").addClass("active")

    $("#new_player_info").hide()
    $("#pay_info").hide()
    $("#game_mgr_menu").hide()
    $("#manager_menu").hide()

    $.getJSON("/query_admin", function (data) {
        if (data == null) {
            return
        }
        if (data == 1) {
            $("#new_player_info").show()
            $("#pay_info").show()
            $("#game_mgr_menu").show()
            $("#manager_menu").show()
            console.log("is admin user")
        }
    });
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

function setInputTime(input_name, year, month, day) {
    if ( ("" + month).length == 1) {
        month = "0" + month
    }
    if ( ("" + day).length == 1) {
        day = "0" + day
    }
    $('#' + input_name).val(year + "-" + month + "-" + day)
}

function parseDataTime(time) {
    var obj = new Object()

    var idx = time.indexOf('-')
    obj.year = time.substring(0, idx)
    time = time.substring(idx+1)

    idx = time.indexOf('-')
    obj.month = time.substring(0, idx)
    obj.day = time.substring(idx+1)

    return obj
}

function new_player_query(pageIdx, bYear, bMonth, bDay, eYear, eMonth, eDay) {
    if (isWaitQuery) return;
    isWaitQuery = true;

    var date = {pageIdx:pageIdx, bYear:bYear, bMonth:bMonth, bDay:bDay, eYear:eYear, eMonth:eMonth, eDay:eDay}

    $.getJSON("/new_player_query", date, function (data) {
        if (data == null) return
        clear_table()
        $("#data_table tr:gt(0):eq(0)").hide()
        curPageIdx = data.CurPage
        totalPageCount = data.TotalPage

        $.each(data.UserList, function(i, item) {
            var row = "<tr>"
            row += "<td>" + item.UserId + "</td>"
            row += "<td>" + item.UserName + "</td>"
            row += "<td>" + (item.TotalOnlineSeconds/60).toFixed(2) + "</td>"
            row += "<td>" + item.MatchTimes + "</td>"
            if (item.Channel == "") {
                row += "<td>官方渠道</td>"
            } else {
                row += "<td>" + item.Channel + "</td>"
            }
            row += "<td>" + item.CreateTime + "</td>"
            row += "</tr>"
            $("#data_table tr:last").after(row)
        })
        $("#data_table_title").text("新增用户-" + data.TotalCount + "人")
        $('#table_page_num').text((curPageIdx+1) + "/" + totalPageCount)
        if (0 == curPageIdx) {
            $('#table_prve_btn').attr("disabled", true)
        } else {
            $('#table_prve_btn').attr("disabled", false)
        }
        if (curPageIdx == totalPageCount) {
            $('#table_next_btn').attr("disabled", true)
        } else {
            $('#table_next_btn').attr("disabled", false)
        }
        isWaitQuery = false
    });
}

function query(pageIdx) {
    var beginTime = $('#begin_time_input').val()
    var endTime = $('#end_time_input').val()
    var begin = parseDataTime(beginTime)
    var end = parseDataTime(endTime)
    new_player_query(pageIdx, begin.year, begin.month, begin.day, end.year, end.month, end.day)
}

function query_time() {
    curPageIdx = 0
    totalPageCount = 0
    query(0)
}

function query_prve_page() {
    if (curPageIdx > 0) {
        query(curPageIdx-1)    
    }
}

function query_next_page() {
    if ((curPageIdx+1) < totalPageCount) {
        query(curPageIdx+1)
    }
}

$(document).ready(function () {
    initMenu()

    initTable("新增用户", ["用户ID", "昵称", "在线时长(分钟)", "比赛次数", "来源渠道", "时间"])

    var beginTime = new Date()
    beginTime.setHours(0)
    beginTime.setMinutes(0)
    beginTime.setSeconds(0)
    var endTime = new Date()
    endTime.setHours(0)
    endTime.setMinutes(0)
    endTime.setSeconds(0)
    endTime.setMilliseconds(0)
    endTime.setTime(endTime.getTime() + 24 * 60 * 60 * 1000)

    $('#begin_time').datetimepicker({
        language:  'zh-CN',
        weekStart: 1,
        todayBtn:  1,
        autoclose: 1,
        todayHighlight: 1,
        startView: 2,
        minView: 2,
        forceParse: 0,
        showMeridian: 1,
        initialDate: beginTime
    });
    setInputTime('begin_time_input', beginTime.getFullYear(), beginTime.getMonth()+1, beginTime.getDate())

    $('#end_time').datetimepicker({
        language:  'zh-CN',
        weekStart: 1,
        todayBtn:  1,
        autoclose: 1,
        todayHighlight: 1,
        startView: 2,
        minView: 2,
        forceParse: 0,
        showMeridian: 1,
        initialDate: endTime
    });
    setInputTime('end_time_input', endTime.getFullYear(), endTime.getMonth()+1, endTime.getDate())

    $('#query').click(query_time)
    $('#table_prve_btn').click(query_prve_page)
    $('#table_next_btn').click(query_next_page)

    $('#table_prve_btn').attr("disabled", true)
    $('#table_page_num').text(curPageIdx + "/" + totalPageCount)
    $('#table_next_btn').attr("disabled", true)
});