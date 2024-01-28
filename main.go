package main

import (
	"crawler/collect"
	"crawler/engine"
	"crawler/log"
	"crawler/parse/douban"
	"crawler/proxy"
	"fmt"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	//log
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("log init")

	//proxy
	proxyURLs := []string{"http://127.0.0.1:7890"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
	}

	cookie := "bid=znKnt-7lWzE; _ga=GA1.1.1024234719.1704808192; _ga_RXNMP372GL=GS1.1.1704808192.1.0.1704808200.52.0.0; viewed=\"1007305\"; ll=\"108303\"; _pk_id.100001.8cb4=75ab89c34ff29dff.1706082828.; douban-fav-remind=1; ct=y; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1706234447%2C%22https%3A%2F%2Fwww.baidu.com%2Flink%3Furl%3DEMf3oP-2srnkNkXMhP79IWiuKe7MPfsgf27b7QMRlJzAuaKR8rHfIpD8P7V1Nybs%26wd%3D%26eqid%3Dd03af640009885e40000000365b31248%22%5D; _pk_ses.100001.8cb4=1; __utma=30149280.1024234719.1704808192.1706144791.1706234448.10; __utmc=30149280; __utmz=30149280.1706234448.10.3.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; ap_v=0,6.0; dbcl2=\"277697100:TT+ck9rXL/M\"; ck=83M2; frodotk_db=\"6acfa8ca95fff80d406efa5fc1ffe827\"; push_noty_num=0; push_doumail_num=0; __utmt=1; __utmv=30149280.27769; __utmb=30149280.17.5.1706235397784"

	var workList = make([]*collect.Request, 0, 1000)
	for i := 0; i <= 0; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		workList = append(workList, &collect.Request{
			Url:       str,
			WaitTime:  1 * time.Second,
			Cookie:    cookie,
			ParseFunc: douban.ParseURL,
		})
	}

	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 30000 * time.Millisecond,
		Logger:  logger,
		Proxy:   p,
	}
	for _, r := range workList {
		fmt.Println(r.Url)
	}

	s := engine.NewSchedule(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(workList),
	)

	s.Run()

	//for len(workList) > 0 {
	//	items := workList
	//	workList = nil
	//	for _, item := range items {
	//		body, err := f.Get(item)
	//		time.Sleep(1 * time.Second)
	//		if err != nil {
	//			logger.Error("read content failed", zap.Error(err))
	//			continue
	//		}
	//		res := item.ParseFunc(body, item)
	//		for _, item := range res.Items {
	//			logger.Info("result", zap.String("get url:", item.(string)))
	//		}
	//		workList = append(workList, res.Requests...)
	//	}
	//
	//}

}
