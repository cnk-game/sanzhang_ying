package models

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"poker-admin/util"
)

const (
	UserType_Admin   = 1
	UserType_Channel = 2

	admin_table_nameC = "website_admin"

	SessionKey_UserId   = "userId"
	SessionKey_UserName = "userName"
	SessionKey_UserType = "userType"
	SessionKey_Channel  = "channel"
)

type AdminInfo struct {
	Channel  string `bson:"Channel"`
	UserId   string `bson:"UserId"`
	UserName string `bson:"UserName"`
	UserPwd  string `bson:"UserPwd"`
	Mark  string `bson:"Mark"`
	UserType int    `bson:"UserType"`
}

var AdminMgr *AdminManager

type AdminManager struct {
	AdminMap map[string]*AdminInfo
}

func init() {
	AdminMgr = &AdminManager{}
	AdminMgr.AdminMap = make(map[string]*AdminInfo)
}

func (m *AdminManager) Login(username, userpwd string) (int, *AdminInfo) {
	var adminInfo *AdminInfo
	adminInfo = m.AdminMap[username]
	if nil != adminInfo {
		if adminInfo.UserPwd != userpwd {
			return 2, nil
		}
		fmt.Println("login==================>1")
		return 0, adminInfo
	}
	adminInfo = loadAdmin(username)
	if nil == adminInfo {
		return 1, nil
	}
	if adminInfo.UserPwd != userpwd {
		return 2, nil
	}
	m.AdminMap[adminInfo.UserName] = adminInfo
	fmt.Println("login==================>2")
	return 0, adminInfo
}

func (m *AdminManager) Logout(username string) {
	delete(m.AdminMap, username)
}

func (m *AdminManager) GetInfo(username string) *AdminInfo {
	var adminInfo *AdminInfo
	adminInfo = m.AdminMap[username]
	if nil != adminInfo {
		return adminInfo
	}
	return loadAdmin(username)
}

func (m *AdminManager) RemoveUser(username string) {
	m.Logout(username)
	info := m.GetInfo(username)
	if nil != info && len(info.Channel) > 0 {
		delAdmin(username)
	}
}

func (m *AdminManager) AddUser(username, userpwd, channel,mark string) *AdminInfo {
	adminInfo := &AdminInfo{}
	adminInfo.Channel = channel
	adminInfo.Mark = mark
	adminInfo.UserName = username
	adminInfo.UserPwd = userpwd
	adminInfo.UserType = UserType_Channel
	adminInfo.UserId = bson.NewObjectId().Hex()
	fmt.Println("AdminManager.AddUser = ", adminInfo)

	addAdmin(adminInfo)
	return adminInfo
}

func (m *AdminManager) GetAllUser(userlist interface{}) {
	session := util.GetLogSession()
	c := session.DB(util.WebsiteDBName).C(admin_table_nameC)
	defer session.Close()

	err := c.Find(nil).All(userlist)
	if err != nil {
		fmt.Println("GetAllUser => error")
	}
}


////////////////////////////////////////////////////////////////////////////////
func IsAdmin(usertype int) bool {
	return usertype == UserType_Admin
}

func IsAdminChannel(channel string) bool {
	return len(channel) == 0
}

////////////////////////////////////////////////////////////////////////////////
func addAdmin(info *AdminInfo) {
	session := util.GetLogSession()
	c := session.DB(util.WebsiteDBName).C(admin_table_nameC)
	defer session.Close()

    fmt.Println("==============addAdmin: ", info)
	err := c.Insert(info)
	if err != nil {
		fmt.Println("AddAdmin => error")
	}
}

func delAdmin(UserName string) {
	session := util.GetLogSession()
	c := session.DB(util.WebsiteDBName).C(admin_table_nameC)
	defer session.Close()

	err := c.Remove(bson.M{"UserName": UserName})
	if err != nil {
		fmt.Println("DelAdmin => error")
	}
}

func loadAdmin(UserName string) *AdminInfo {
	session := util.GetLogSession()
	c := session.DB(util.WebsiteDBName).C(admin_table_nameC)
	defer session.Close()

	info := &AdminInfo{}
	err := c.Find(bson.M{"UserName": UserName}).One(info)
	if err != nil {
		fmt.Println("LoadAdmin => error")
		return nil
	}
	return info
}
