package main

import (
	"flag"
	"fmt"
	"strings"

	//"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	// "github.com/pkg/profile"
	"log"
	"net/http"
	"net/http/httputil"

	b64 "encoding/base64"

	"github.com/hollychen503/htpasswd"
)

var port string
var filePath string

func init() {
	flag.StringVar(&port, "port", "80", "give me a port number")
	flag.StringVar(&filePath, "htpasswd", "./htpasswd", "htpasswd file path")
}

func main() {
	// defer profile.Start().Stop()
	flag.Parse()

	log.Println("htpasswd file:", filePath)

	http.HandleFunc("/", whoamI)

	fmt.Println("Starting up on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func whoamI(w http.ResponseWriter, req *http.Request) {
	//u, _ := url.Parse(req.URL.String())
	log.Println("+++++++++++ request ++++++++++++++++++++++++++++")
	dump, _ := httputil.DumpRequest(req, true)
	log.Println(string(dump))
	log.Println("------------------------------------------------")

	uri := req.Header.Get("X-Forwarded-Uri")
	if uri == "/v2/" {
		log.Println("Access /v2/. Ignore")
		return
	}

	// 获取用户名，密码
	// Authorization: Basic dGVzdHVzZXI6dGVzdHBhc3N3b3Jk
	usrpw := req.Header.Get("Authorization")
	if len(usrpw) == 0 {
		log.Println("Without Authorization header. ignore.")
		return
	}
	// 取出账号密码 b64
	upslice := strings.Fields(usrpw)
	if len(upslice) < 2 {
		log.Println("Invalid basic Authorization info ")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// 解码
	sDec, err := b64.StdEncoding.DecodeString(upslice[1])
	if err != nil {
		log.Println("Can not decode basic auth info")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	log.Println("username:password =", string(sDec))

	decSli := strings.Split(string(sDec), ":")
	if len(decSli) < 2 {
		log.Println("Malformed auth info")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	realName := decSli[0]
	realPw := decSli[1]

	//tmpPwHash := htpasswd.HashedPasswords(map[string]string{})

	//err = tmpPwHash.SetPassword(tmpName, tmpPw, htpasswd.HashBCrypt)
	//if err != nil {
	//	fmt.Println("failed to gen password")
	//	return
	//}

	///
	passwords, err := htpasswd.ParseHtpasswdFile(filePath)
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to parse htpasswd file on", filePath)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for k, v := range passwords {
		fmt.Println(k, ":", v)
	}
	/*
		if tmpPwHash[tmpName] != passwords[tmpName] {
			log.Println(" invalid password")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	*/
	//  CompareHashAndPassword(hashedPassword, password []byte)
	fmt.Println("hashedPw:", passwords[realName])
	//fmt.Println("pw:", tmpPw)
	err = bcrypt.CompareHashAndPassword([]byte(passwords[realName]), []byte(realPw))
	if err != nil {
		log.Println(err)
		log.Println("Invalid user name or password")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	//  /v2/holly/hello/manifests/latest
	//  0 1 2     3     4
	segs := strings.Split(uri, "/")
	fmt.Println(segs)
	if len(segs) < 2 {
		log.Println("Unknown URL, ignore")
		return
	}
	// segs[0] 为空 ！
	//if segs[2] == "common" { // 暂时不开放 common
	//	log.Println("Use common namespace.")
	//	return
	//}

	if segs[2] != realName {
		log.Println("namespace name does not matched user name!", segs[2], realName)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	log.Println("Authorized!")

	return
	/*
		queryParams := u.Query()
		wait := queryParams.Get("wait")
		if len(wait) > 0 {
			duration, err := time.ParseDuration(wait)
			if err == nil {
				time.Sleep(duration)
			}
		}
		hostname, _ := os.Hostname()
		fmt.Fprintln(w, "Hostname:", hostname)
		ifaces, _ := net.Interfaces()
		for _, i := range ifaces {
			addrs, _ := i.Addrs()
			// handle err
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				fmt.Fprintln(w, "IP:", ip)
			}
		}
		req.Write(w)
	*/
}
