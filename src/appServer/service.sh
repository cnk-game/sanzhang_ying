appPid=app.pid

usage="Usage: service.sh [start|stop]"

mkLogFolders()
{
  mkdir -p logger
}

start()
{
	nohup ./appServer http://1251002466.cdn.myqcloud.com/1251002466/assets/poker/PokerGame.apk > logger/console.log 2>&1 &
	echo $! > app.pid
}

stop()
{
  if [ -f $appPid ]; then
    if kill `cat $appPid` > /dev/null 2>&1; then
      echo "kill app."
    fi
    rm -f $appPid
  fi
}

case $1 in
  (start)
    mkLogFolders
    stop
    echo StartService.
    start
    ;;
  (stop)
    echo StopService.
    stop
    ;;
  (*)
    echo $usage
    exit 1
    ;;
esac

