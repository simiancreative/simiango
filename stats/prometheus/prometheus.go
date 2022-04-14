package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/simiancreative/simiango/server"

	promcore "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Register(collector promcore.Collector) {
	promcore.MustRegister(collector)
}

func Handle() {
	r := server.GetRouter()
	r.GET("/metrics", prometheusHandler())
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// adding a service monitor to k8s example
//
// apiVersion: monitoring.coreos.com/v1
// kind: ServiceMonitor
// metadata:
//   name: {monitor name}
//   namespace: metrics
// spec:
//   endpoints:
//   - port: metrics
//   jobLabel: jobLabel
//   namespaceSelector:
//     matchNames:
//     - metrics
//   selector:
//     matchLabels: # the service must have a mathcing label
//       app: {service name}
