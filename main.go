package main

import (
	"crawler/collect"
	"crawler/log"
	"crawler/parse/douban"
	"crawler/proxy"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	//log
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	//proxy
	proxyURLs := []string{"http://127.0.0.1:7890"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
	}

	cookie := "bid=znKnt-7lWzE; _ga=GA1.1.1024234719.1704808192; _ga_RXNMP372GL=GS1.1.1704808192.1.0.1704808200.52.0.0; viewed=\"1007305\"; ll=\"108303\"; _pk_id.100001.8cb4=75ab89c34ff29dff.1706082828.; __utmz=30149280.1706082830.6.2.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; douban-fav-remind=1; ct=y; __utmc=30149280; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1706096086%2C%22https%3A%2F%2Fwww.baidu.com%2Flink%3Furl%3D7ngR135WAkMgSzVT0WFSuatUc6GdBgdQRcOcwTARgNGZuMl8gk5rvArQ97dvoDNA%26wd%3D%26eqid%3Df6e2384f0013e7790000000365b0c205%22%5D; _pk_ses.100001.8cb4=1; ap_v=0,6.0; __utma=30149280.1024234719.1704808192.1706085239.1706096089.8; __utmt=1; __utmb=30149280.200.3.1706100091986"

	var workList []*collect.Request
	for i := 0; i <= 0; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		workList = append(workList, &collect.Request{
			Url:       str,
			Cookie:    cookie,
			ParseFunc: douban.ParseURL,
		})
	}

	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 300000 * time.Millisecond,
		Proxy:   p,
	}

	for len(workList) > 0 {
		items := workList
		workList = nil
		for _, item := range items {
			body, err := f.Get(item)
			time.Sleep(1 * time.Second)
			if err != nil {
				logger.Error("read content failed", zap.Error(err))
				continue
			}
			res := item.ParseFunc(body, item)
			for _, item := range res.Items {
				logger.Info("result", zap.String("get url:", item.(string)))
			}
			workList = append(workList, res.Requests...)
		}

	}

}
