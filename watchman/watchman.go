package watchman

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" //nolint:gosec

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InitWatchman
/**
 * @Description:
 * @param addr
 */
func InitWatchman(addr string) {
	go func() {
		if addr == "" {
			addr = ":6060"
		}
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			fmt.Println(err)
		}
	}()
}
