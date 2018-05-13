package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type WeatherInfoJson struct {
	Weatherinfo   WeatherinfoObject
	Performesinfo PerformESObject
}

type PerformESObject struct {
	PageNo     string
	SearchCode string
}

type WeatherinfoObject struct {
	StaffID      string
	EsSearchCode string
	PageNo       int
	WD           string
	WS           string
	SD           string
	WSE          string
	Time         string
	IsRadar      string
	Radar        string
}

func (wiObject *WeatherInfoJson) ReadFile(filename string) *WeatherInfoJson {

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
	}

	if err := json.Unmarshal(bytes, &wiObject); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
	}
	return wiObject
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func main() {
	var wi *WeatherInfoJson = new(WeatherInfoJson)
	//var result interface{}
	wioj := wi.ReadFile("weather.json")
	//设置url请求的参数
	v := url.Values{}
	v.Set("StaffID", wioj.Weatherinfo.StaffID)
	v.Set("LogStaffID", "3322")
	v.Set("Name", "Joke")
	v.Set("Status", "1")
	v.Set("Staffs", "test")

	var apiUrl string = "http://localhost/vSearch.GenericAPI/api/newStaffInfo/"
	var s string
	s = strings.Join([]string{apiUrl, wioj.Weatherinfo.EsSearchCode, "/", strconv.Itoa(wioj.Weatherinfo.PageNo)}, "")

	//body
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	request, err := http.NewRequest("POST", s, body)
	if err != nil {
		fmt.Println("Fatal Error", err.Error())
	}

	request.Header.Set("Content-Type", "application/json;param=value")

	resp, err := client.Do(request)
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Fatal error")
	}
	fmt.Println(string(content))

	if Exists("result.txt") {
		del := os.Remove("./result.txt")
		if del != nil {
			fmt.Println(del)
		}
	}

	file, err := os.Create("result.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	fmt.Fprintf(file, string(content))
	// fmt.Println(wioj.Weatherinfo.City)
	// fmt.Println(wioj.Weatherinfo.WD)
	// fmt.Println(wioj.Performesinfo.PageNo)
}
