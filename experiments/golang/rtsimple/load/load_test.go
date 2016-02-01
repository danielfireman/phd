package load

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()
	g := Generator{
		MaxQPS:           3,
		Duration:         5 * time.Second,
		URL:              srv.URL,
		ConcurrencyLevel: 2,
	}
	r := g.Run()
	if r.QPS <= float64(g.ConcurrencyLevel) {
		t.Fatal("qps got:%d want:>=%d", r.QPS, g.ConcurrencyLevel)
	}
}
