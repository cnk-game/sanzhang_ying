var isWaitQueryMain = false
var isWaitQueryOnline = false



function login() {
    var username = $('#username').val()
    var userpwd = $('#userpwd').val()
    if (username.length == 0 || userpwd.length == 0) {
        alert("用户名或密码不能为空")
        return
    }
    $.post("/login", {username:username, userpwd:userpwd}, function (data) {
        if (0 == data) {
            location.href = "/main"
        } else {
            $('#login_tip').text("用户名或密码验证失败")
        }
    })
}

$(document).ready(function () {
    $('#login_btn').click(login)
});