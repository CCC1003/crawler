package main

import (
	"crawler/collect"
	"crawler/collector"
	"crawler/collector/sqlstorage"
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

	var storage collector.Storage
	storage, err = sqlstorage.New(
		sqlstorage.WithSqlUrl("root:123456@tcp(47.92.241.189:3306)/crawler?charset=utf8"),
		sqlstorage.WithLogger(logger.Named("sqlDB")),
		sqlstorage.WithBatchCount(2),
	)
	if err != nil {
		logger.Error("create sqlStorage failed")
		return
	}
	seeds := make([]*collect.Task, 0, 1000)
	seeds = append(seeds, &collect.Task{
		Property: collect.Property{
			Name: "douban_book_list",
		},
		Fetcher: f,
		Storage: storage,
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
