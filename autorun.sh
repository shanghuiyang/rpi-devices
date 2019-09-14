
pushd /home/pi
nohup ./devices.pi > devices.log 2>&1 &
popd

sudo motion
