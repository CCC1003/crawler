package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type ProxyFunc func(r *http.Request) (*url.URL, error)

func RoundRobinProxySwitcher(ProxyUrls ...string) (ProxyFunc, error) {
	if len(ProxyUrls) < 1 {
		return nil, errors.New("proxy URL list is empty")
	}
	urls := make([]*url.URL, len(ProxyUrls))

	for i, u := range ProxyUrls {
		parsedU, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parsedU
	}

	return (&roundRobinSwitcher{urls, 0}).GetProxy, nil
}

type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

// GetProxy 取余算法实现轮询调度
func (r *roundRobinSwitcher) GetProxy(pr *http.Request) (*url.URL, error) {

	index := atomic.AddUint32(&r.index, 1) - 1
	u := r.proxyURLs[index%uint32(len(r.proxyURLs))]

	return u, nil
}
