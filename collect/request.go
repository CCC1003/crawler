package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"
)

type Task struct {
	Url      string
	Cookie   string
	WaitTime time.Duration
	Reload   bool
	MaxDepth int
	RootReq  *Request
	Fetcher  Fetcher
}

type Request struct {
	unique    string
	Task      *Task
	Url       string
	Method    string
	Depth     int
	Priority  int
	ParseFunc func([]byte, *Request) ParseResult
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
