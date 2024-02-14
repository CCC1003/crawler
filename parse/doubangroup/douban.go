package doubangroup

import (
	"crawler/collect"
	"fmt"
	"regexp"
	"time"
)

const urlListRe = `(https://www.doubangroup.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`
const ContentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div class="aside">`

var cookie = "bid=znKnt-7lWzE; _ga=GA1.1.1024234719.1704808192; _ga_RXNMP372GL=GS1.1.1704808192.1.0.1704808200.52.0.0; viewed=\"1007305\"; ll=\"108303\"; _pk_id.100001.8cb4=75ab89c34ff29dff.1706082828.; doubangroup-fav-remind=1; ct=y; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1706234447%2C%22https%3A%2F%2Fwww.baidu.com%2Flink%3Furl%3DEMf3oP-2srnkNkXMhP79IWiuKe7MPfsgf27b7QMRlJzAuaKR8rHfIpD8P7V1Nybs%26wd%3D%26eqid%3Dd03af640009885e40000000365b31248%22%5D; _pk_ses.100001.8cb4=1; __utma=30149280.1024234719.1704808192.1706144791.1706234448.10; __utmc=30149280; __utmz=30149280.1706234448.10.3.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; ap_v=0,6.0; dbcl2=\"277697100:TT+ck9rXL/M\"; ck=83M2; frodotk_db=\"6acfa8ca95fff80d406efa5fc1ffe827\"; push_noty_num=0; push_doumail_num=0; __utmt=1; __utmv=30149280.27769; __utmb=30149280.17.5.1706235397784"

var DoubangroupTask = &collect.Task{
	Property: collect.Property{
		Name:     "find_douban_sun_room",
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
		Cookie:   cookie,
	},

	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			var roots []*collect.Request
			for i := 0; i < 25; i++ {
				str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
				roots = append(roots, &collect.Request{
					Priority: 1,
					Url:      str,
					Method:   "GET",
					RuleName: "解析网站URL",
				})
			}
			return roots, nil
		},
		Trunk: map[string]*collect.Rule{
			"解析网站URL": &collect.Rule{ParseFunc: ParseURL},
			"解析阳台房":   &collect.Rule{ParseFunc: GetSunRoom},
		},
	},
}

func ParseURL(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(urlListRe)

	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(
			result.Requests, &collect.Request{
				Method:   "GET",
				Url:      u,
				Task:     ctx.Req.Task,
				Depth:    ctx.Req.Depth + 1,
				RuleName: "解析阳台房",
			})
	}
	return result, nil
}

func GetSunRoom(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(ContentRe)

	ok := re.Match(ctx.Body)
	if !ok {
		return collect.ParseResult{
			Items: []interface{}{},
		}, nil
	}

	result := collect.ParseResult{
		Items: []interface{}{ctx.Req.Url},
	}

	return result, nil
}
