
pushd /home/pi
nohup ./car > car.log 2>&1 &
popd

sudo motion
