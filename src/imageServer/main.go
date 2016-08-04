package main

import (
    "fmt"
    "log"
    "net/http"
    ."github.com/qiniu/api/conf"
    "github.com/qiniu/api/rs"
)

var bucketName string = "dry-poker"

func uptoken() string {
    putPolicy := rs.PutPolicy {
        Scope:         bucketName,
        //CallbackUrl: callbackUrl,   
        //CallbackBody:callbackBody,    
        //ReturnUrl:   returnUrl,  
        //ReturnBody:  returnBody,    
        //AsyncOps:    asyncOps,    
        //EndUser:     endUser,    
        //Expires:     expires,   
    }
    return putPolicy.Token(nil)
}

func deleteKey(key string) {
    client := rs.New(nil)
    client.Delete(nil, bucketName, key)
}

func RequestUploadTokenHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("http => ", r)
    r.ParseForm()
    fmt.Println("RemoteAddr => ", r.RemoteAddr)
    userId := r.Form.Get("userId")
    userToken := r.Form.Get("userToken")
    fileName := r.Form.Get("fileName")
    deleteKey(fileName)
    token := uptoken()

    fmt.Println("userId => ", userId)
    fmt.Println("userToken => ", userToken)
    fmt.Println("fileName => ", fileName)
    fmt.Println("token => ", token)
    
    fmt.Fprintf(w, `{result:"success", token:"` + string(token) + `"}`)
}

func main() {
	// 七牛云存储分配的私钥
    ACCESS_KEY = "SlTOK6JxHcUSdiMte5N42ISWqRXJ6rKgfLPecVdt"
    SECRET_KEY = "ZcdPWd5roFILjIet2VyMKCl-iFx17TDYWHjsDC9L"
    http.HandleFunc("/", RequestUploadTokenHandler)
    fmt.Println("启动成功")
    log.Fatal(http.ListenAndServe(":7002", nil))
}
