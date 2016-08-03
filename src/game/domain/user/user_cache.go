package user

import (
	"github.com/golang/glog"
	"pb"
	"sync"
	"util"
)

const (
	cacheMaxEntries = 2000
)

type UserCache struct {
	sync.RWMutex
	cache *util.Cache
}

var cache *UserCache

func init() {
	cache = &UserCache{}
	cache.cache = util.New(cacheMaxEntries)
}

func GetUserCache() *UserCache {
	return cache
}

func (c *UserCache) SetUser(user *pb.UserDef) {
	c.Lock()
	defer c.Unlock()

	glog.V(2).Info("===>设置用户信息缓存user:", user)

	c.cache.Add(user.GetUserId(), user)
}

func (c *UserCache) GetUser(userId string) *pb.UserDef {
	c.Lock()
	user, ok := c.cache.Get(userId)
	c.Unlock()

	if !ok {
		u, err := FindByUserId(userId)
		if err != nil {
			return nil
		}
		r := FindMatchRecord(userId)
		user = u.BuildMessage(r.BuildMessage())

		c.Lock()
		c.cache.Add(userId, user)
		c.Unlock()
	}

	switch v := user.(type) {
	case *pb.UserDef:
		return v
	}

	return nil
}
