while :
do
    if !(pgrep -x "sunny.py" >/dev/null)
    then
        echo "restart sunny.py"
        setsid ./sunny.py --clientid=xxxx > sunny.log 2>&1 &
        sleep 1s
    fi

    sleep 10s
done
