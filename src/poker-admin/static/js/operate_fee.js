var isWaitQuery = false
var curPageIdx = 0
var totalPageCount = 0
var beginYear = 0
var beginMonth = 0
var beginDay = 0
var endYear = 0
var endMonth = 0
var endDay = 0

function calcIntervalDate(now, intervalDay) {
    var newDate = new Date()
    newDate.setFullYear(now.getFullYear(), now.getMonth(), now.getDate())
    newDate.setHours(0)
    newDate.setMinutes(0)
    newDate.setSeconds(0)
    newDate.setMilliseconds(0)
    newDate.setTime(newDate.getTime() - 24 * 60 * 60 * 1000 * intervalDay)
    return newDate
}

function initMenu() {
    $("#opreate_mgr_menu").addClass("active")
    $("#opreate_mgr_sub_menu").show()
    $("#fee_mgr").addClass("active")

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

function operate_fee_query(beginY, beginM, beginD, endY, endM, endD) {
    if (isWaitQuery) return;
    isWaitQuery = true;

    var date = { bYear:beginY, bMonth:beginM, bDay:beginD, eYear:endY, eMonth:endM, eDay:endD}

    $.getJSON("/operate_fee_query", date, function (data) {
        if (data == null) return;
        $("#data_table tr:gt(0):eq(0)").hide();
        clear_table()

        $.each(data, function(i, item) {
            var row = "<tr>"
            row += "<td>" + item.Date + "</td>"
            row += "<td>" + item.FeeGold + "</td>"
            
            row += "</tr>"
            $("#data_table tr:last").after(row);
        });
        isWaitQuery = false;
    });
}

function parseDate(begin, end) {
    var obj = new Object()

    var idx = begin.indexOf('-')
    obj.beginY = begin.substring(0, idx)

    begin = begin.substring(idx+1)
    idx = begin.indexOf('-')
    obj.beginM = begin.substring(0, idx)
    obj.beginD = begin.substring(idx+1)

    idx = end.indexOf('-')
    obj.endY = end.substring(0, idx)
    end = end.substring(idx+1)
    idx = end.indexOf('-')
    obj.endM = end.substring(0, idx)
    obj.endD = end.substring(idx+1)

    return obj
}

function setInputTime(input_name, year, month, day) {
    if ( ("" + month).length == 1) {
        month = "0" + month
    }
    if ( ("" + day).length == 1) {
        day = "0" + day
    }
    //console.log(input_name + "============" + year + "-" + month + "-" + day)
    $('#' + input_name).val(year + "-" + month + "-" + day)
}

function query_list(pageIdx) {
    operate_fee_query( beginYear, beginMonth, beginDay, endYear, endMonth, endDay)
}

function query_lastday(lastday) {
    curPageIdx = 0
    totalPageCount = 0
    var end = new Date()
    var begin = calcIntervalDate(end, lastday)

    beginYear = begin.getFullYear()
    beginMonth = begin.getMonth() + 1
    beginDay = begin.getDate()
    endYear = end.getFullYear()
    endMonth = end.getMonth() + 1
    endDay = end.getDate()
    setInputTime('begin_date_input', beginYear, beginMonth, beginDay)
    setInputTime('end_date_input', endYear, endMonth, endDay)
    query_list(curPageIdx)
}

function query_time() {
    curPageIdx = 0
    totalPageCount = 0
    var begin = $('#begin_date_input').val()
    var end = $('#end_date_input').val()
    var obj = parseDate(begin, end)

    beginYear = obj.beginY
    beginMonth = obj.beginM
    beginDay = obj.beginD
    endYear = obj.endY
    endMonth = obj.endM
    endDay = obj.endD
    query_list(curPageIdx)
}

function query_last15day() {
    query_lastday(15)
}

function query_last30day() {
    query_lastday(30)
}

function query_prve_page() {
    if (curPageIdx > 0) {
        query_list(curPageIdx-1)    
    }
}

function query_next_page() {
    if ((curPageIdx+1) < totalPageCount) {
        query_list(curPageIdx+1)
    }
}

$(document).ready(function () {
    initMenu()

    initTable("综合信息", ["日期", "金币数"])

    var endDate = new Date()
    var beginDate = calcIntervalDate(endDate, 7)
    $('#begin_date').datetimepicker({
        language:  'zh-CN',
        weekStart: 1,
        todayBtn:  1,
        autoclose: 1,
        todayHighlight: 0,
        startView: 2,
        minView: 2,
        forceParse: 0,
        initialDate: beginDate
    });
    setInputTime('begin_date_input', beginDate.getFullYear(), beginDate.getMonth()+1, beginDate.getDate())
    $('#end_date').datetimepicker({
        language:  'zh-CN',
        weekStart: 1,
        todayBtn:  1,
        autoclose: 1,
        todayHighlight: 1,
        startView: 2,
        minView: 2,
        forceParse: 0,
        initialDate: endDate
    });
    setInputTime('end_date_input', endDate.getFullYear(), endDate.getMonth()+1, endDate.getDate())

    $('#query_time').click(query_time)
    $('#query_last15day').click(query_last15day)
    $('#query_last30day').click(query_last30day)
    $('#table_prve_btn').click(query_prve_page)
    $('#table_next_btn').click(query_next_page)

    $('#table_prve_btn').attr("disabled", true)
    $('#table_page_num').text(curPageIdx + "/" + totalPageCount)
    $('#table_next_btn').attr("disabled", true)
    
    query_lastday(7)
});