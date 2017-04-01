/**
维度			分析事项
url     	访问次数 / 访问文件数 /   静态文件 / 带宽 /
ip     		访客数 / 点击总数 / 比例 / 日期
状态码  		总请求次数 / 有效请求次数 / 失败请求次数  / 404 /
日志文件  	大小



log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for" '
		      '$request_time $upstream_response_time ';

 */
package main

import (
	//"bufio"
	//"fmt"
	//"io"
	"os"
	//"regexp"
	"strings"
	"regexp"
	"io"
	"bufio"
	"html/template"
	"net/http"
	//"time"
	"strconv"
	"time"
	//"fmt"
	//"sort"
	"fmt"
)

var urlArr =  make(map[string]urlData)
var ipArr = make(map[string]ipData)
var statusArr = make(map[string]statusData)

var totalRequest int
var successRequest int
var failedRequest int
var uniqueVisitors int
var uniqueFiles int
var ipHits int
var unique404 int
var logSize int64
var bandwidth int
var logFile string

type assignData struct {
	UrlArr	map[string]urlData
	IpArr	map[string]ipData
	StatusArr map[string]statusData
	Date 	string
	TotalRequest int
	SuccessRequest int
	FailedRequest int
	UniqueVisitors int
	UniqueFiles int
	IpHits int
	Unique404 int
	LogSize int64
	Bandwidth int
	LogFile string
	//logFileData  map[string]logFile
}



type ipData struct {
	Ip		string
	Visitors	int
	His		int
	Proportion	float32
	Bandwidth	float32
}

type urlData struct {
	Url		string
	Visitors	int
	His		int
	Proportion	float32
	Bandwidth	float32
	Protocol	string
	Method		string
}

type statusData struct {
	Visitors	int
	His		int
	Proportion	float32
	Bandwidth	float32
	Code 		string
}


func checkErr(err error)  {
	if err != nil {
		panic("open failed!")
	}
}

func updateStatus(arr []string)  {
	_,err := statusArr[arr[3]]
	totalRequest ++
	statusCodd , _ := strconv.Atoi(arr[3])
	if statusCodd >=  200 || statusCodd < 400{
		successRequest ++
	} else {
		failedRequest ++
	}
	if statusCodd == 404 {
		unique404 ++
	}
	if err == false{
		statusArr[arr[3]] = statusData{1,1,0.00,0.11,arr[3]}
	} else {
		temp := statusArr[arr[3]]
		temp.Visitors ++
		temp.His ++
		statusArr[arr[3]] = temp
	}
}

func updateIp(arr []string)  {
	_,err := ipArr[arr[0]]
	if err == false{
		uniqueVisitors ++
		ipArr[arr[0]] = ipData{arr[0],1,1,0.00,0.11}
	} else {
		temp := ipArr[arr[0]]
		temp.Visitors ++
		temp.His ++
		ipArr[arr[0]] = temp
	}
}

func updateUrl(arr []string)  {
	_,err := urlArr[arr[2]]

	//fmt.Println(arr)
	//[205.251.252.110 [05/Apr/2016:02:25:07 +0000] "GET /apiserver/index.php/home/?c=index&a=weekDownload&app_id=556164008&device=iPhone7%2C1&device_id=h14647d39a1c23a70f1c715f1cc0104dd&idfa=F1618DFE-A432-4780-ABA2-9B85E7455C3E&l=en&localVersion=5505&mac=e0%3Ab5%3A2d%3A28%3A12%3Ac8&os_version=8.0.2&platform=1&position=7_0_0&time=1459823107 HTTP/1.1" 200 52 "-" "Amazon CloudFront" "99.247.137.239"]

	tempBandwidth,_ := strconv.Atoi(arr[4])
	bandwidth = bandwidth + tempBandwidth
	//fmt.Println(res)
	if err == false{
		uniqueFiles ++
		request := strings.Fields(strings.Trim(arr[2],string('"')))
		urlArr[arr[2]] = urlData{request[1],1,1,1.00,1.00,request[0],request[2]}
		//fmt.Println(urlArr)


	} else {
		temp := urlArr[arr[2]]
		//temp["Visitors"] = 11
		temp.Visitors ++
		temp.His ++
		//fmt.Println(temp)
		urlArr[arr[2]] = temp
		//fmt.Println(urlArr[arr[2]])
		//os.Exit(0)


		//urlArr[arr[2]].Visitors ++
		//urlArr[arr[2]].His ++
	}
	//fmt.Println(res)
}
func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}


func init()  {
	//urlArr["aa"] = urlData{"url",1,1,1.00,1.00,"http/1.1","post"}
	//fmt.Println(urlArr,ipArr,statusArr,logFileData)

	//f, err := os.Open("log.log")
	logFile = "test.log"
	f, err := os.Open(logFile)
	logSize ,_ = FileSize(logFile)

	checkErr(err)
	defer f.Close()
	b := bufio.NewReader(f)
	line, err := b.ReadString('\n')
	i := 0
	for ; err == nil; line, err = b.ReadString('\n') {
		i ++
		//if i > 3 {
		//	fmt.Println(ipArr);
		//	sort.Interface(ipArr)
		//	fmt.Println(ipArr)
		//	os.Exit(0)
		//}
		//fmt.Print(line)
		//reg := regexp.MustCompile(`^(.*?)\s-\s-\s\[(.*?)\]\s\"(.*?)\s(.*?)\s(.*?)\"\s(.*?)\s(.*?)\s\"(.*?)\"\s\"(.*?)\"\s\"(.*?)\"`)
		//regStr := regexp.Match(reg,line)

		reg3 := regexp.MustCompile(`(\d{0,3}\.\d{0,3}\.\d{0,3}\.\d{0,3}|\[(.*?)\]|\"(.*?)\"|\d+)`) //\"(.*?)\s(.*?)\s(.*?)\"|
		res := reg3.FindAllString(line, -1)
		updateUrl(res)
		updateIp(res)
		updateStatus(res)
		//l := len(res)
		//fmt.Println(l)
		//if l != 8 {
		//	fmt.Println("------",l)
		//	for k,v:= range res{
		//		fmt.Println(k,v)
		//	}
		//	fmt.Println(line)
		//}

		//os.Exit(0)
	}
	if err == io.EOF {
		//fmt.Print(line)
	} else {
		panic("read occur error!")
	}



}

func logInfo(w http.ResponseWriter, r *http.Request)  {

	currentTime := time.Now().Local()
	formatDate := currentTime.Format("2006-01-02 15:04:05.000")


	assignData := assignData{
		/*	"UrlArr": */		urlArr,
		/*	"IpArr":*/		ipArr,
		/*	"StatusArr":*/		statusArr,
		/*	"Date":*/		formatDate,
		/*	"TotalRequest":*/ 	totalRequest,
		/*	"SuccessRequest":*/ 	successRequest,
		/*	"FailedRequest": */	failedRequest,
		/*	"UniqueVisitors": */	uniqueVisitors,
		/*	"UniqueFiles": 	*/	uniqueFiles,
		/*	"IpHits": 	*/	ipHits,
		/*	"Unique404": 	*/	unique404,
		/*	"LogSize": 	*/	logSize,
		/*	"Bandwidth": 	*/	bandwidth,
		/*	"LogFile": 	*/	logFile,
	}

	//}

	t,_ := template.ParseFiles("index.html")
	t.Execute(w, assignData)
}

func main() {



	http.HandleFunc("/index", logInfo)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}


	fmt.Println("end")
}