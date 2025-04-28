package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
)

type Balancer struct {
	logger   *slog.Logger
	strategy Strategy
	proxies  map[string]*httputil.ReverseProxy
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
