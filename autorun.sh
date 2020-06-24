
pushd /home/pi

#sudo ifconfig wlan0 down
#sudo pppd call dial &
#sudo ifconfig wlan0 up
#sleep 1m
#sudo route add -net 0.0.0.0 ppp0

#nohup python sunny_video.py --clientid=xxxxx > sunny_video.log 2>&1 &
#nohup python sunny_car.py --clientid=xxxxx > sunny_car.log 2>&1 &

# backup car.log
cp car.log car.log.bak

# set env var
export BAIDU_SPEECH_APP_KEY="xxxx"
export BAIDU_SPEECH_SECRET_KEY="xxxx"

# -E means: preserve current user's env to su
nohup sudo -E ./car > car.log 2>&1 &

#sleep 5s
#nohup ./auto_restart_sunny.sh > auto_restart_sunny.log 2>&1 &

#nohup ./ip > ip.log 2>&1 &
hostname -I | mutt -s "ip address" yangsh@sina.com

popd

sudo motion

