var loader = null

function initMenu() {
    $("#game_mgr_menu").addClass("active")
    $("#game_mgr_sub_menu").show()
    $("#notify_mgr").addClass("active")

}

$(document).ready(function () {
    
    initMenu()

    $('#send_broadcast').click(function () {
        content = $('#broadcastContent').val()
        var date = {content:content}
        loader = layer.load('请稍候…')
        $.getJSON("/game_mgr/send_system_message", date, function (data) {
            layer.close(loader)
            if (data == null) {
                return
            }
            if (data == "succeed") {
                console.log("succeed")
            }
        });
        $("#broadcastContent").val("")
    })

    $('#send_prize').click(function () {
        userId = $('#userId').val()
        goldCount = $('#goldCount').val()
        diamondCount = $('#diamondCount').val()
        propId = $('#propId').val()
        propCount = $('#propCount').val()
        prizeDesc = $('#prizeDesc').val()

        var date = {userId:userId, content:prizeDesc, gold:goldCount, diamond:diamondCount, itemType:propId, itemCount:propCount}
        loader = layer.load('请稍候…')
        $.getJSON("/game_mgr/send_user_prize_mail", date, function (data) {
            layer.close(loader)
            if (data == null) {
                return
            }
            if (data == "succeed") {
                console.log("succeed")
            }
        });

        $("#userId").val("")
        $("#goldCount").val("")
        $("#diamondCount").val("")
        $("#propId").val("")
        $("#propCount").val("")
        $("#prizeDesc").val("")
    })
});