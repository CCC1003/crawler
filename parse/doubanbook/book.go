package doubanbook

import (
	"crawler/collect"
	"go.uber.org/zap"
	"regexp"
	"strconv"
)

var cookie = "bid=1FXa1uNW184; ap_v=0,6.0; __utmc=30149280; __utma=81379588.2066573877.1707911958.1707911958.1707911958.1; __utmc=81379588; __utmz=81379588.1707911958.1.1.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/tag/%E7%A5%9E%E7%BB%8F%E7%BD%91%E7%BB%9C/; _vwo_uuid_v2=DA9FE1E570CA6C406894B512F8E0FFF37|b8a8abc50ae2c0157b2b1f7cbfed5626; _pk_ref.100001.3ac3=%5B%22%22%2C%22%22%2C1707911959%2C%22https%3A%2F%2Fwww.douban.com%2Ftag%2F%E7%A5%9E%E7%BB%8F%E7%BD%91%E7%BB%9C%2F%22%5D; ll=\"108231\"; __utma=30149280.1439395164.1707911953.1707911953.1707912991.2; __utmz=30149280.1707912991.2.2.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; push_noty_num=0; push_doumail_num=0; __utmv=30149280.27769; viewed=\"36659680\"; dbcl2=\"277697100:4QPz0knmZzs\"; ck=tMH1; __utmt=1; _pk_ses.100001.3ac3=*; frodotk_db=\"0e7dd7260b8fc8dc2b79857222af5a94\"; __utmt_douban=1; __utmb=30149280.26.10.1707912991; __utmb=81379588.21.10.1707911958; _pk_id.100001.3ac3=4cc4ab2090f23199.1707911959..1707915062.undefined."

var DoubanBookTask = &collect.Task{
	Property: collect.Property{
		Name:     "douban_book_list",
		WaitTime: 2,
		MaxDepth: 5,
		Cookie:   cookie,
	},
	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			roots := []*collect.Request{
				&collect.Request{
					Priority: 1,
					Url:      "https://book.douban.com",
					Method:   "GET",
					RuleName: "数据tag",
				},
			}
			return roots, nil
		},
		Trunk: map[string]*collect.Rule{
			"数据tag": &collect.Rule{ParseFunc: ParseTag},
			"书籍列表":  &collect.Rule{ParseFunc: ParseBookList},
			"书籍简介": &collect.Rule{
				ItemFields: []string{
					"书名",
					"作者",
					"页数",
					"出版社",
					"得分",
					"价格",
					"简介",
				},
				ParseFunc: ParseBookDetail,
			},
		},
	},
}

const regexpStr = `<a href="([^"]+)" class="tag">([^<]+)</a>`

func ParseTag(ctx *collect.Context) (collect.ParseResult, error) {
	compile := regexp.MustCompile(regexpStr)
	matches := compile.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		result.Requests = append(
			result.Requests, &collect.Request{
				Method:   "GET",
				Task:     ctx.Req.Task,
				Url:      "https://book.douban.com" + string(m[1]),
				Depth:    ctx.Req.Depth + 1,
				RuleName: "书籍列表",
			})
	}

	zap.S().Debugln("parse book tag,count:", len(result.Requests))

	return result, nil
}

const BooklistRe = `<a.*?href="([^"]+)" title="([^"]+)"`

func ParseBookList(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(BooklistRe)
	submatch := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}

	for _, m := range submatch {

		req := &collect.Request{
			Priority: 100,
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍简介",
		}

		req.TmpData = &collect.Temp{}
		req.TmpData.Set("book_name", string(m[2]))
		result.Requests = append(result.Requests, req)
	}

	//result.Requests = result.Requests[:1]
	zap.S().Debugln("parse book list,count:", len(result.Requests))

	return result, nil
}

var autoRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)

// var public = regexp.MustCompile(`<span class="pl">出版社:</span>([^<]+)<br/>`)
var public = regexp.MustCompile(`<span class="pl">出版社:</span>[\d\D]*?<a.*?>([^<]+)</a>`)
var pageRe = regexp.MustCompile(`<span class="pl">页数:</span> ([^<]+)<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
var intoRe = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)

func ParseBookDetail(ctx *collect.Context) (collect.ParseResult, error) {
	bookName := ctx.Req.TmpData.Get("book_name")
	page, _ := strconv.Atoi(ExtraString(ctx.Body, pageRe))

	book := map[string]interface{}{
		"书名":  bookName,
		"作者":  ExtraString(ctx.Body, autoRe),
		"页数":  page,
		"出版社": ExtraString(ctx.Body, public),
		"得分":  ExtraString(ctx.Body, scoreRe),
		"价格":  ExtraString(ctx.Body, priceRe),
		"简介":  ExtraString(ctx.Body, intoRe),
	}
	data := ctx.Output(book)

	result := collect.ParseResult{
		Items: []interface{}{data},
	}

	zap.S().Debugln("parse book detail", data)
	return result, nil
}

func ExtraString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
