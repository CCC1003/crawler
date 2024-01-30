package main

import (
	"crawler/collect"
	"crawler/engine"
	"crawler/log"
	"crawler/proxy"
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
	//Get
	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 30000 * time.Millisecond,
		Logger:  logger,
		Proxy:   p,
	}

	seeds := make([]*collect.Task, 0, 1000)
	seeds = append(seeds, &collect.Task{
		Property: collect.Property{
			Name: "js_find_douban_sun_room",
		},
		Fetcher: f,
	})

	s := engine.NewEngine(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
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
