pushd /home/pi

sudo ifconfig wlan0 down
sudo pppd call dial &
sudo ifconfig wlan0 up
sleep 1m
sudo route add -net 0.0.0.0 ppp0

nohup python sunny_video.py --clientid=xxxxxxxxx > sunny_video.log 2>&1 &
nohup python sunny_car.py --clientid=xxxxxxxxx > sunny_car.log 2>&1 &
nohup sudo ./car > car.log 2>&1 &

popd

sudo motion
