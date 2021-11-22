package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	ip, err := getOutboundIP()
	if err != nil {
		log.Printf("[uploadsvr]failed to ip address, error: %v", err)
	}
	if ip != nil {
		log.Printf("http://%v:8083/upload", ip)
	}

	http.HandleFunc("/upload", upload)
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Printf("[uploadsvr]failed to ListenAndServe, error: %v", err)
		return
	}

}
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		page := fmt.Sprintf(tpl, "")
		w.Write([]byte(page))
		return
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		page := fmt.Sprintf(tpl, "upload "+handler.Filename+" failed")
		w.Write([]byte(page))
		log.Printf("[uploadsvr]failed to create file, error: %v", err)
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		page := fmt.Sprintf(tpl, "upload "+handler.Filename+" failed")
		w.Write([]byte(page))
		log.Printf("[uploadsvr]failed to copy file, error: %v", err)
		return
	}
	page := fmt.Sprintf(tpl, "uploaded "+handler.Filename+" successfully")
	w.Write([]byte(page))

	log.Printf("[uploadsvr]received %s in success", handler.Filename)
}

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

var tpl = `
<html>
	<head>
		<title>upload file to pi</title>
	</head>
	<body>
		<form enctype="multipart/form-data" action="/upload" method="post">
			<input type="file" name="uploadfile">
            <input type="hidden" name="token" value="{...{.}...}">
            <br><br>
			<input type="submit" value="upload" style="color:white;background-color:steelblue;font-size:15px;">
		</form>
	</body>
	<p style="color:red;font-size:15px;">
	<br>
	%v
	</p>
</html>
`
