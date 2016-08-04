var isWaitQuery = false
var isWaitQueryOrder = false
var isWaitQueryUser = false

var queryUserId = ""
var curPageIdx = 0
var totalPageCount = 0

function initMenu() {
    $("#statistics_menu").addClass("active")
    $("#statistics_sub_menu").show()
    $("#pay_info").addClass("active")

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

function setInputTime(input_name, year, month, day, hour, minute) {
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

function new_pay_query(pageIdx, bYear, bMonth, bDay, eYear, eMonth, eDay) {
    if (isWaitQuery) return;
    isWaitQuery = true;

    var date = {pageIdx:pageIdx, bYear:bYear, bMonth:bMonth, bDay:bDay, eYear:eYear, eMonth:eMonth, eDay:eDay}

    $.getJSON("/pay_query", date, function (data) {
        if (data == null) 
            return
        clear_table()
        $("#data_table tr:gt(0):eq(0)").hide()

        curPageIdx = data.CurPage
        totalPageCount = data.TotalPage

        $.each(data.PayLogList, function(i, item) {
            var row = "<tr>"
            row += "<td>" + item.OrderId + "</td>"
            row += "<td>" + item.UserId + "</td>"
            row += "<td>" + item.Amount + "</td>"
            row += "<td>" + item.Channel + "</td>"
            row += "<td>" + item.Time + "</td>"
            row += "</tr>"
            $("#data_table tr:last").after(row)
        })
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

function orderid_query() {
    var orderId = $('#orderid').val()

    if (isWaitQueryOrder) return;
    isWaitQueryOrder = true;
    $.getJSON("/pay_query", {orderId:orderId}, function (data) {
        if (data == null) return
        $("#data_table tr:gt(0):eq(0)").hide()
        clear_table()
        var row = "<tr>"
        row += "<td>" + data.OrderId + "</td>"
        row += "<td>" + data.UserId + "</td>"
        row += "<td>" + data.Amount + "</td>"
        row += "<td>" + data.Channel + "</td>"
        row += "<td>" + data.Time + "</td>"
        row += "</tr>"
        $("#data_table tr:last").after(row)
        $('#table_page_num').text("0/0")
        $('#table_prve_btn').attr("disabled", true)
        $('#table_next_btn').attr("disabled", true)
        isWaitQueryOrder = false
    });
}

function query_by_user(userId, pageIdx) {
    if (isWaitQueryUser) return;
    isWaitQueryUser = true;
    $.getJSON("/pay_query", {userId:userId, pageIdx:pageIdx}, function (data) {
        if (data == null) return
        $("#data_table tr:gt(0):eq(0)").hide()
        clear_table()

        curPageIdx = data.CurPage
        totalPageCount = data.TotalPage

        $.each(data.PayLogList, function(i, item) {
            var row = "<tr>"
            row += "<td>" + item.OrderId + "</td>"
            row += "<td>" + item.UserId + "</td>"
            row += "<td>" + item.Amount + "</td>"
            row += "<td>" + item.Channel + "</td>"
            row += "<td>" + item.Time + "</td>"
            row += "</tr>"
            $("#data_table tr:last").after(row)
        })
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
        isWaitQueryUser = false
    });
}

function userid_query() {
    queryUserId = $('#userid').val()
    query_by_user(queryUserId, 0)
    
}

function query_by_time(pageIdx) {
    var beginTime = $('#begin_time_input').val()
    var endTime = $('#end_time_input').val()
    var begin = parseDataTime(beginTime)
    var end = parseDataTime(endTime)
    new_pay_query(pageIdx, begin.year, begin.month, begin.day, end.year, end.month, end.day)
}

function time_query() {
    queryUserId = ""
    curPageIdx = 0
    totalPageCount = 0
    query_by_time(curPageIdx)
}

function prve_page_query() {
    if (queryUserId == "") {
        // 按时间查询
        if (curPageIdx > 0) {
            query_by_time(curPageIdx-1)    
        }
    } else {
        // 按用户查询
        if (curPageIdx > 0) {
            query_by_user(curPageIdx-1)    
        }
    }
}

function next_page_query() {
    if (queryUserId == "") {
        // 按时间查询
        if ((curPageIdx+1) < totalPageCount) {
            query_by_time(curPageIdx+1)
        }
    } else {
        // 按用户查询
        if ((curPageIdx+1) < totalPageCount) {
            query_by_user(curPageIdx+1)
        }
    }
}

$(document).ready(function () {
    initMenu()

    initTable("充值记录", ["订单号", "用户ID", "金额", "支付渠道", "时间"])

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

    $('#time_query').click(time_query)
    $('#orderid_query').click(orderid_query)
    $('#userid_query').click(userid_query)
    $('#table_prve_btn').click(prve_page_query)
    $('#table_next_btn').click(next_page_query)

    $('#table_prve_btn').attr("disabled", true)
    $('#table_page_num').text(curPageIdx + "/" + totalPageCount)
    $('#table_next_btn').attr("disabled", true)
});