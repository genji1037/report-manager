package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"strconv"
	"strings"
	"time"
)

func InitCustomerMetrics(r *gin.Engine) {
	customMetrics := []*ginprometheus.Metric{
		{
			ID:          "reqDurR",
			Name:        "request_duration_seconds_route",
			Description: "The HTTP request latencies in seconds per route.",
			Type:        "summary_vec",
			Args:        []string{"code", "method", "handler", "host", "url"},
		},
	}
	p := ginprometheus.NewPrometheus("gin", customMetrics)
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.String()
		index := strings.Index(url, "?")
		if index >= 0 {
			url = url[:index]
		}
		return url
	}

	r.Use(func(c *gin.Context) {
		if c.Request.URL.String() == p.MetricsPath {
			c.Next()
			return
		}

		start := time.Now()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		url := p.ReqCntURLLabelMappingFn(c)

		for i, metric := range p.MetricsList {
			if metric.ID == "reqDurR" {
				p.MetricsList[i].MetricCollector.(*prometheus.SummaryVec).WithLabelValues(status, c.Request.Method, c.HandlerName(), c.Request.Host, url).Observe(elapsed)
			}
		}
	})

	p.Use(r)
}
