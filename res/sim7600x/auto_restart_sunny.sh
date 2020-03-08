while :
do
    if !(pgrep -x "sunny" >/dev/null)
    then
        echo "restart sunny"
        setsid ./sunny clientid xxxxxxxx &
        sleep 1s
    fi
    sleep 2s
done
