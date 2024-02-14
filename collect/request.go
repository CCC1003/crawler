package collect

import (
	"crawler/collector"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"regexp"
	"sync"
	"time"
)

type Property struct {
	Name     string        `json:"name"`
	Url      string        `json:"url"`
	Cookie   string        `json:"cookie"`
	WaitTime time.Duration `json:"wait_time"`
	Reload   bool          `json:"reload"`
	MaxDepth int64         `json:"max_depth"`
}

type Task struct {
	Property
	Visited     map[string]bool
	VisitedLock sync.Mutex
	Fetcher     Fetcher
	Storage     collector.Storage
	Rule        RuleTree
	Logger      *zap.Logger
}

type Context struct {
	Body []byte
	Req  *Request
}

func (c *Context) GetRule(ruleName string) *Rule {
	return c.Req.Task.Rule.Trunk[ruleName]
}

func (c *Context) Output(data interface{}) *collector.DataCell {
	res := &collector.DataCell{}
	res.Data = make(map[string]interface{})
	res.Data["Task"] = c.Req.Task.Name
	res.Data["Rule"] = c.Req.RuleName
	res.Data["Data"] = data
	res.Data["Url"] = c.Req.Url
	res.Data["Time"] = time.Now().Format("2006-01-02 15:04:05")
	return res
}

func (c *Context) ParseJSReg(name string, reg string) ParseResult {
	re := regexp.MustCompile(reg)

	matches := re.FindAllSubmatch(c.Body, -1)
	result := ParseResult{}
	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(
			result.Requests, &Request{
				Method:   "GET",
				Task:     c.Req.Task,
				Url:      u,
				Depth:    c.Req.Depth + 1,
				RuleName: name,
			})
	}
	return result
}

func (c *Context) OutputJS(reg string) ParseResult {
	re := regexp.MustCompile(reg)
	ok := re.Match(c.Body)
	if !ok {
		return ParseResult{
			Items: []interface{}{},
		}
	}
	result := ParseResult{
		Items: []interface{}{c.Req.Url},
	}
	return result
}

type Request struct {
	unique   string
	Task     *Task
	Url      string
	Method   string
	Depth    int64
	Priority int64
	RuleName string
	TmpData  *Temp
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}

func (r *Request) Check() error {
	if r.Depth > r.Task.MaxDepth {
		return errors.New("Max depth limit reached")
	}
	return nil
}
func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.Url + r.Method))
	return hex.EncodeToString(block[:])
}
