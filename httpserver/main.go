package main

import (
	"flag"
	"fmt"
	"httpserver/controller"
	"httpserver/metrics"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	http.HandleFunc("/", rootHandler)
	c, python, java := true, false, "no!"
	//定义healthz用于返回200
	http.HandleFunc("/healthz", healthz)
	metrics.Register()
	//定义延迟操作 并定义promethes监控方式
	controller.RegisterRoutes()
	fmt.Println(c, python, java)
	err := http.ListenAndServe(":80", nil)
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", rootHandler)
	// mux.HandleFunc("/healthz", healthz)
	// mux.HandleFunc("/debug/pprof/", pprof.Index)
	// mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}

}

func healthz(w http.ResponseWriter, r *http.Request) {
	//问题4当访问 localhost/healthz 时，应返回 200
	//调用healthz函数，并返回200返回值
	io.WriteString(w, "200\n")
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")

	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}

	_, ok := os.LookupEnv("VERSION")
	if !ok {
		os.Setenv("VERSION", "golandv1.18.1")
	}
	version := os.Getenv("VERSION")
	//写入respose Header
	io.WriteString(w, version)

	addr, statuscode := GetipStatus(w, r)
	//打印输出addr status并记录进日志
	fmt.Printf("addr%s : status is %d", addr, statuscode)
	glog.V(2).Infof("addr is %s: status is %d", addr, statuscode)

	delay := randInt(10, 20)
	time.Sleep(time.Microsecond * time.Duration(delay))
	io.WriteString(w, "==========Details of the http request header: ============")
	req, err := http.NewRequest("GET", "http://service1", nil)
	if err != nil {
		fmt.Println("%s", err)
	}
	lowCaseHeader := make(http.Header)
	for key, value := range r.Header {
		lowCaseHeader[strings.ToLower(key)] = value
	}
	glog.Info("headers:", lowCaseHeader)
	req.Header = lowCaseHeader
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Info("HTTP get failed with error:", "error", err)
	} else {
		glog.Info("HTTP get successed")
	}
	if resp != nil {
		resp.Write(w)
	}

	glog.V(4).Infof("Respond in %d ms", delay)
}

func GetipStatus(w http.ResponseWriter, r *http.Request) (addr string, statuscode int) {
	//为了防止中间存在代理模式所以要定义RemoteHeader,用于获取客户端真实IP
	const remoteAddrHeader = "X-AppEngine-Remote-Addr"
	if addr = r.Header.Get(remoteAddrHeader); addr != "" {
		//获取客户端ip,并赋予变量
		addr = r.RemoteAddr
		//删除header中代理相关配置，防止攻击
		r.Header.Del(remoteAddrHeader)
	} else {
		r.RemoteAddr = "127.0.0.1"
		addr = r.RemoteAddr
	}

	statuscode = http.StatusOK
	w.WriteHeader(statuscode)
	return addr, statuscode

}
