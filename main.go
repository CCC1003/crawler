package main

import (
	"crawler/collect"
	"crawler/engine"
	"crawler/limiter"
	"crawler/log"
	"crawler/proxy"
	"crawler/storage"
	"crawler/storage/sqlstorage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	//log
	plugin := log.NewStdoutPlugin(zapcore.DebugLevel)
	logger := log.NewLogger(plugin)
	logger.Info("log init")

	zap.ReplaceGlobals(logger)

	//proxy
	//proxyURLs := []string{"http://221.231.13.198:1080", "http://47.106.120.76:8080", "http://122.116.150.2:9000"}
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

	var storage storage.Storage
	storage, err = sqlstorage.New(
		sqlstorage.WithSqlUrl("root:123456@tcp(47.92.241.189:3306)/crawler?charset=utf8"),
		sqlstorage.WithLogger(logger.Named("sqlDB")),
		sqlstorage.WithBatchCount(2),
	)
	if err != nil {
		logger.Error("create sqlStorage failed")
		return
	}

	//2秒1个
	secondLimit := rate.NewLimiter(limiter.Per(1, 2*time.Second), 1)
	//60秒20个
	minuteLimit := rate.NewLimiter(limiter.Per(20, 1*time.Minute), 20)

	multiLimiter := limiter.MultiLimiter(secondLimit, minuteLimit)

	seeds := make([]*collect.Task, 0, 1000)
	seeds = append(seeds, &collect.Task{
		Property: collect.Property{
			Name: "douban_book_list",
		},
		Fetcher: f,
		Storage: storage,
		Limit:   multiLimiter,
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
