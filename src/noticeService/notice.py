
import jpush as jpush
from conf import app_key, master_secret
from time import sleep, localtime

def push_broadcast(_jpush, title, content):
	push = _jpush.create_push()
	push.audience = jpush.all_
	push.notification = jpush.notification(alert=content)
	push.notification['android'] = {'title':title}
	push.platform = jpush.all_
	push.options = {'time_to_live':3600}
	push.send()

if __name__ == '__main__':
	_jpush = jpush.JPush(app_key, master_secret)
	while True:
		now = localtime()
		print "check time", now.tm_year, now.tm_mon, now.tm_mday, now.tm_hour, now.tm_min
		if now.tm_hour == 12 and now.tm_min == 30 :
			print "send 12:30 notice"
			push_broadcast(_jpush, "【免费的午餐】时间到！", "现在登录游戏马上领取金币！")
		if now.tm_hour == 18 and now.tm_min == 30 :
			print "send 18:30 notice"
			push_broadcast(_jpush, "【天上掉馅饼】时间到！", "现在登录游戏马上领取金币！")
		sleep(1*60)
