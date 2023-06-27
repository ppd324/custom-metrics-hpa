package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

/*
// 自定义指标
var (

	//Desc是每个普罗米修斯度量使用的指标。它本质上是指标的不可变元数据。此包中包含的普通指标实现在后台管理其 Desc。
	desc prometheus.Desc
	//只增计数器(接口类型)，可通过Inc() 或Add(float64)进行计数器值的增加，一般可用于统计连接数量、请求个数，计数器通常用于计算服务的请求、完成的任务、发生的错误等
	counter prometheus.Counter
	//CounterVec 是一个收集器，它捆绑了一组计数器，这些计数器共享相同的Desc，但其变量标签的值不同。例如可统计访问不同API的请求数量，不同状态码的数量
	counterVec prometheus.CounterVec
	//gauge表示一个数值，表示可以任意上下移动的单个数值。gauge通常用于测量值，如温度或当前内存使用情况，但也用于可以上升和下降的“计数”，如正在运行的goroutine的数量。
	gauge prometheus.Gauge
	//GaugeVec 是一个收集器，它捆绑了一组gauge，这些gauge共享相同的Desc，但其变量标签具有不同的值。可以按不同维度划分的同一事物
	gaugeVec prometheus.GaugeVec
	//Histogram 直方图对事件或样本流的单个观测值进行计数。
	histogram prometheus.Histogram
	//HistogramVec是一个收集器，它捆绑了一组直方图，这些直方图共享相同的Desc，但变量标签具有不同的值。可按不同维度划分的同一事物
	histogramVec prometheus.HistogramVec

)
*/
type Metrics struct {
	qps     *prometheus.CounterVec //统计api请求数量，用于计算qps
	cpuTemp prometheus.Gauge       //模拟当前CPU温度
}

func NewCustomMetrics() *Metrics {
	m := &Metrics{
		qps: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "total_req_num",
				Help: "Number of total http requests.",
			},
			[]string{"status", "api"},
		),
		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU.",
		}),
	}
	return m
}
func main() {
	m := NewCustomMetrics()
	//将自定义指标添加到到默认prometheus监控指标上
	prometheus.MustRegister(m.qps)
	prometheus.MustRegister(m.cpuTemp)
	//注册新的prometheus监控指标，只监控注册的自定义指标
	/*
		reg := prometheus.NewRegistry()
		reg.MustRegister(m.qps)
		reg.MustRegister(m.cpuTemp)
	*/
	m.cpuTemp.Set(50.0)
	beginTemp := 50.0
	reqNum := 0
	//http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		statusCode := 200
		m.qps.WithLabelValues("200", "/").Inc()
		reqNum++
		log.Println("total request num is:", reqNum)
		switch path {
		case "/cpuTempUp":
			m.cpuTemp.Add(10.0)
			m.qps.WithLabelValues("200", "cpuTempUp").Inc()
			beginTemp += 10
			log.Println("cpu temp up,current is:", beginTemp)
		case "/cpuTempDown":
			m.cpuTemp.Sub(10.0)
			beginTemp -= 10
			m.qps.WithLabelValues("200", "cpuTempDown").Inc()
			log.Println("cpu temp down,current is:", beginTemp)
		case "/metrics":
			//调用prometheus默认监控器，将默认监控指标发送至请求端
			promhttp.Handler().ServeHTTP(w, r)
			//可自定义配置注册器，将监控器器所注册的指标发送至请求端
			//promhttp.HandlerFor(reg,promhttp.HandlerOpts{Registry: reg}).ServeHTTP(w,r)
		default:
			w.WriteHeader(statusCode)
			_, err := w.Write([]byte("custom metrics"))
			if err != nil {
				return
			}
		}
	})

	http.ListenAndServe(":3000", nil)
}
