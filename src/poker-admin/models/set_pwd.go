package models

import (
    "poker-admin/util"
    "fmt"
    "labix.org/v2/mgo/bson"
)



func SetUserPwd(userid, password string) {
    session := util.GetSession()
    c := session.DB(util.WebsiteDBName).C(admin_table_nameC)
    defer session.Close()

    err := c.Update(bson.M{"UserId": userid}, bson.M{"$set": bson.M{"UserPwd": password}})
    if err != nil {
        fmt.Println("SetUserPwd => success")
    }else{
        fmt.Println("SetUserPwd => error")
    }
}
