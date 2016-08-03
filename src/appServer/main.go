package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"os"
)

var apkDownloadURL string = ""

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	/*
	fmt.Println("======================================================")
	fmt.Println("RemoteAddr => ", r.RemoteAddr)
	fmt.Println("UserAgent => ", r.UserAgent())
	fmt.Println("Proto => ", r.Proto)
	fmt.Println("URL => ", r.URL)
	fmt.Println("Host => ", r.Host)
	fmt.Println("RemoteAddr => ", r.RemoteAddr)
	fmt.Println(strings.Contains(r.UserAgent(), "MicroMessenger"))
	fmt.Println("======================================================")
	*/
	if strings.Contains(r.UserAgent(), "MicroMessenger") {
		fmt.Fprintf(w, `
<html>
<body>
	<br>
	<br>
	<div class="text" style="text-align:center;font-size:60px">KO三国 - Android版下载</div>
	<br>
	<hr/>
	<br>
	<br>
	<div class="text" style="text-align:center;font-size:36px">微信二维码扫描不支持直接下载</div>
	<br>
	<div class="text" style="text-align:center;font-size:36px">请选择右上角菜单，选择"在浏览器中打开"进行下载</div>
</body>
</html>
		`)
	} else {
		//http.Redirect(w, r, `http://1251002466.cdn.myqcloud.com/1251002466/assets/mobile_test/DryGame20141113.apk`, http.StatusFound)
		http.Redirect(w, r, apkDownloadURL, http.StatusFound)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("启动失败，请提供安装包下载地址")
		return
	}
	apkDownloadURL = os.Args[1]
	fmt.Println("启动成功，APK下载地址：", apkDownloadURL)
	http.HandleFunc("/", RequestHandler)
	log.Fatal(http.ListenAndServe(":10002", nil))	
}
