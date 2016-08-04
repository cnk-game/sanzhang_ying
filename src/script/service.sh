gamePid=game.pid

usage="Usage: service.sh [start|stop]"

mkLogFolders()
{
  mkdir -p logger/game
}

export GODEBUG=gctrace=1
export PORT=55555

start()
{
	nohup ./game --logtostderr=true -db=127.0.0.1 -log_dir=logger/game > logger/game/console.log 2>&1 &
	echo $! > game.pid
}

stop()
{
  if [ -f $gamePid ]; then
    if kill `cat $gamePid` > /dev/null 2>&1; then
      echo "kill game."
    fi
    rm -f $gamePid
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

