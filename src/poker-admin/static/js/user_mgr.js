var isWaitQuery = false
var loader = null

function initMenu() {
    $("#manager_menu").addClass("active")
    $("#manager_sub_menu").show()
    $("#user_mgr").addClass("active")


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

function query_user_list() {
    if (isWaitQuery) return;
    isWaitQuery = true;

    $.getJSON("/query_user_list", function (data) {
        if (data == null) return
        $("#main-content").show()
        clear_table()
        $("#data_table tr:gt(0):eq(0)").hide()

        $.each(data, function(i, item) {
            var row = "<tr>"
            row += "<td>" + item.Channel + "</td>"
            row += "<td>" + item.Mark + "</td>"
            row += "<td>" + item.UserName + "</td>"
            row += "<td><button class=\"btn btn-danger btn-xs\" onclick=\"removeUser(\'"+item.UserName+"\')\"><i class=\"fa fa-trash-o \"></i></button></td>"
            row += "</tr>"
            $("#data_table tr:last").after(row)
        })
        isWaitQuery = false
    });
}

function addUser() {
    console.log("======================>add user")
    var channel = $('#channel').val()
    var mark = $('#mark').val()
    var username = $('#username').val()
    var userpwd = $('#userpwd').val()
    $("#channel").val("")
    $("#mark").val("")
    $("#username").val("")
    $("#userpwd").val("")
    loader = layer.load('请稍候…')
    $.post("/add_user", {channel:channel, username:username, userpwd:userpwd, mark:mark}, function (data) {
        layer.close(loader)
        if (0 == data.Result) {
            var row = "<tr>"
            row += "<td>" + data.Channel + "</td>"
            row += "<td>" + data.Mark + "</td>"
            row += "<td>" + data.UserName + "</td>"
            row += "<td><button class=\"btn btn-danger btn-xs\" onclick=\"removeUser(\'"+data.UserName+"\')\"><i class=\"fa fa-trash-o \"></i></button></td>"
            // \'
            row += "</tr>"
            $("#data_table tr:last").after(row)
        } else {
            if (3 == data.Result) {
                alert("该用户已经存在")
            } else {
                alert("添加失败")
            }
        }
    })
}

function removeUser(select_username) {
    var username = select_username
    console.log(username)
    $.post("/remove_user", {username:username}, function (data) {
        if (0 == data) {
            query_user_list()
        } else {
            alert("删除失败")
        }
    })
}


$(document).ready(function () {
    
    initMenu()

    initTable("后台用户列表", ["渠道ID","渠道说明", "用户名", "删除用户"])

    $('#add_user_button').click(addUser)

    $("#main-content").hide()

    query_user_list()
});