package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	fmt.Println("start test api")

	//loginFirst()
	testConfig("list")
}

func testConfig(tp string) {
	token := "kv3dj00EAAteePVwbcZ5JyRMGNDYJFHb"
	httpUrl := "http://localhost:8080/admin/"
	switch tp {
	case "list":
		httpUrl += "config"
		data := url.Values{"length": {"20"}, "start": {"0"}, "keyword": {"13"}}
		str := httpGet(httpUrl, data, token)
		fmt.Println(str)
	case "create":
		httpUrl += "config/create"
		data := url.Values{"keystr": {"fromtest"}, "valuestr": {"test12306"}}
		str := httpPostForm(httpUrl, strings.NewReader(data.Encode()), token)
		fmt.Println(str)
	case "update":
		httpUrl += "config/update"
		data := url.Values{"id": {"8"}, "keystr": {"fromtest"}, "valuestr": {"i am update wawwaw"}}
		str := httpPostForm(httpUrl, strings.NewReader(data.Encode()), token)
		fmt.Println(str)
	case "del":
		httpUrl += "config/del"
		data := url.Values{"id": {"8"}, "keystr": {"fromtest"}, "valuestr": {"i am update wawwaw"}}
		str := httpPostForm(httpUrl, strings.NewReader(data.Encode()), token)
		fmt.Println(str)
	case "one":
		httpUrl += "config/one"
		data := url.Values{"id": {"8"}, "keystr": {"fromtest"}, "valuestr": {"i am update wawwaw"}}
		str := httpPostForm(httpUrl, strings.NewReader(data.Encode()), token)
		fmt.Println(str)
	}
	fmt.Println("finish test config")

}

func getUser() {
	httpUrl := "http://localhost:8080/admin/user"
	data := url.Values{"key": {"value"}, "name": {"test111"}, "username": {"test222"}, "userpass": {"123456"}}
	str := httpGet(httpUrl, data, "3hnri00cfk7ejx45icm200dwo0wsnv81")
	fmt.Println(str)
}

func loginFirst() {
	//{"code":0,"message":"","token":"3hnri00cfk7ejx45icm200dwo0wsnv81"}
	httpUrl := "http://localhost:8080/auth/login"
	data := url.Values{"key": {"value"}, "name": {"test111"}, "username": {"test222"}, "userpass": {"123456"}}
	str := httpPostForm(httpUrl, strings.NewReader(data.Encode()), "12306")
	fmt.Println(str)

}

func testJson() {
	httpUrl := "http://localhost:8080/ping"

	data := make(map[string]interface{})
	data["json"] = "are you ok"
	data["username"] = "test111"
	data["name"] = "666"
	bytesData, _ := json.Marshal(data)
	str := httpPostJson(httpUrl, bytes.NewReader(bytesData))
	fmt.Println(str)
}

func testPost() {
	httpUrl := "http://localhost:8080/ping"
	data := url.Values{"key": {"value"}, "name": {"test111"}, "username": {"comeon"}}
	str := httpPostForm(httpUrl, strings.NewReader(data.Encode()), "12306")
	fmt.Println(str)
}

func testGet() {
	httpUrl := "http://localhost:8080/ping"
	data := url.Values{"key": {"value"}, "name": {"test111"}, "username": {"comeon"}}
	str := httpGet(httpUrl, data, "12306")
	fmt.Println(str)
}

func httpClient() *http.Client {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	return client
}

func httpGet(httpUrl string, data url.Values, token string) string {

	u, _ := url.ParseRequestURI(httpUrl)
	u.RawQuery = data.Encode()

	client := httpClient()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Token", token)

	rsps, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return ""
	}
	defer rsps.Body.Close()

	body, err := ioutil.ReadAll(rsps.Body)
	if err != nil {
		fmt.Println("Read body failed:", err)
		return ""
	}

	fmt.Println(string(body))

	return string(body)
}

func httpPostForm(url string, uri io.Reader, token string) string {
	client := httpClient()

	req, _ := http.NewRequest("POST", url, uri)
	req.Header.Add("Token", token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rsps, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return ""
	}
	defer rsps.Body.Close()

	body, err := ioutil.ReadAll(rsps.Body)
	if err != nil {
		fmt.Println("Read body failed:", err)
		return ""
	}

	fmt.Println(string(body))

	return string(body)
}

func httpPostJson(url string, uri io.Reader) string {
	client := httpClient()

	req, _ := http.NewRequest("POST", url, uri)
	req.Header.Add("Token", "10086")
	req.Header.Add("Content-Type", "application/json")

	rsps, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return ""
	}
	defer rsps.Body.Close()

	body, err := ioutil.ReadAll(rsps.Body)
	if err != nil {
		fmt.Println("Read body failed:", err)
		return ""
	}

	fmt.Println(string(body))

	return string(body)
}
