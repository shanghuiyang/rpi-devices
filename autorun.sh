pushd /home/pi

cp car.log car.log.bak
nohup sudo ./car > car.log 2>&1 &

popd

sudo motion

