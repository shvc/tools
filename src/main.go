package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type tokenStruct struct {
	Token string
}

func token(serverURL string, entryKey string, entryOp string) (string, error) {
	// http://host:port/token?entryKey=$key&entryOp=put
	tokenURL := fmt.Sprintf("%s/token?entryKey=%s&entryOp=%s", serverURL, entryKey, entryOp)
	log.Println(tokenURL)
	tokenResp, err := http.Post(tokenURL, "application/json", strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	defer tokenResp.Body.Close()
	tokenBody, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	log.Println(string(tokenBody))
	token := tokenStruct{}
	err = json.Unmarshal(tokenBody, &token)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	log.Printf("put token: %s", token.Token)
	return token.Token, nil
}

func upload(serverURL, entryKey, token, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	// http://host:port/pblocks/$key?token=$ptoken Content-Type: application/octet-stream
	postURL := fmt.Sprintf("%s/pblocks/%s?token=%s", serverURL, entryKey, token)
	postResp, err := http.Post(postURL, "application/octet-stream", file)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer postResp.Body.Close()
	postBody, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(postBody))
	return
}

func download(serverURL, entryKey, token, filename string) {
	// http://host:port/pblocks/$key?token=$token
	getURL := fmt.Sprintf("%s/pblocks/%s?token=%s", serverURL, entryKey, token)
	getResp, err := http.Get(getURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer getResp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	io.Copy(file, getResp.Body)

	log.Println("Download finish ", filename)
	return
}

func md5sum(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	md5h := md5.New()
	io.Copy(md5h, file)
	fmd5 := fmt.Sprintf("%x", md5h.Sum([]byte("")))
	return fmd5, nil
}

func main() {
	host := flag.String("host", "localhost", "hostname")
	port := flag.Int("port", 3000, "port")
	key := flag.String("key", "key", "key")
	ufile := flag.String("ufile", "/tmp/ufile", "upload filename")
	sfile := flag.String("save", "/tmp/downfile", "save download filename")
	num := flag.Int("num", 1024, "upload/download times")
	debug := flag.Bool("debug", false, "enable/disable debug mode")
	ignore := flag.Bool("i", false, "ignore failed validation")

	flag.Parse()
	rand.Seed(time.Now().Unix())
	randValue := rand.Intn(8192)

	serverURL := fmt.Sprintf("http://%s:%d", *host, *port)
	fmt.Printf("key:%s, ufile:%s, dfile:%s debug: %v, num:%d, ignore:%v\n", *key, *ufile, *sfile, *debug, *num, *ignore)
	log.Printf("key:%s, ufile:%s, dfile:%s debug: %v, num:%d, ignore:%v\n", *key, *ufile, *sfile, *debug, *num, *ignore)

	for i := 0; i < *num; i++ {
		randKey := fmt.Sprintf("%s-%04d-%05d", *key, randValue, i)
		ptoken, err := token(serverURL, randKey, "put")
		if err != nil {
			log.Fatal(err)
		}
		upload(serverURL, randKey, ptoken, *ufile)
		umd5, err := md5sum(*ufile)
		if err != nil {
			log.Println("calc ufile md5 error: ", err)
			break
		}

		gtoken, err := token(serverURL, randKey, "get")
		if err != nil {
			log.Fatal(err)
		}
		download(serverURL, randKey, gtoken, *sfile)
		dmd5, err := md5sum(*sfile)
		if err != nil {
			log.Println("calc sfile md5 error: ", err)
			break
		}
		if umd5 != dmd5 {
			log.Printf("checkmd5 failed  %s != %s", umd5, dmd5)
		} else {
			log.Printf("checkmd5 success %s == %s", umd5, dmd5)
		}

	}

}
