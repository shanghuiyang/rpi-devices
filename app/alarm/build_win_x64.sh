CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o alarm.exe alarm.go
zip -FS alarm.zip alarm.exe sounds/*.mp3
