package controller

import (
	"fmt"
	"httpserver/metrics"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func registerImagesRoutes() {
	http.HandleFunc("/images", handleImages)
}

func registerMetricsRoutes() {
	http.Handle("/metrics", promhttp.Handler())
}

// 添加 0-2 秒的随机延时
func handleImages(w http.ResponseWriter, r *http.Request) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	randInt := rand.Intn(2000)
	time.Sleep(time.Millisecond * time.Duration(randInt))
	w.Write([]byte(fmt.Sprintf("<h1>%d<h1>", randInt)))
}

func RegisterRoutes() {
	registerImagesRoutes()
	registerMetricsRoutes()

}
