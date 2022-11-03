package filter

import (
	"github.com/fagongzi/log"
	"net/http"
	"time"

	"github.com/fagongzi/gateway/pkg/pb/metapb"
	"github.com/fagongzi/gateway/pkg/util"
	"github.com/valyala/fasthttp"
)

// Context filter context
type Context interface {
	StartAt() time.Time
	EndAt() time.Time

	OriginRequest() *fasthttp.RequestCtx
	ForwardRequest() *fasthttp.Request
	Response() *fasthttp.Response

	API() *metapb.API
	DispatchNode() *metapb.DispatchNode
	Server() *metapb.Server
	Analysis() *util.Analysis

	SetAttr(key string, value interface{})
	GetAttr(key string) interface{}
}

// Filter filter interface
type Filter interface {
	Name() string
	Init(cfg string) error

	Pre(c Context) (statusCode int, err error)
	Post(c Context) (statusCode int, err error)
	PostErr(c Context, code int, err error)
}

// BaseFilter base filter support default implementation
type BaseFilter struct{}

// Init init filter
func (f BaseFilter) Init(cfg string) error {
	log.Infof("execute baseFilter Init function %s", cfg)
	return nil
}

// Pre execute before proxy
func (f BaseFilter) Pre(c Context) (statusCode int, err error) {
	log.Info("execute baseFilter Pre function")
	return http.StatusOK, nil
}

// Post execute after proxy
func (f BaseFilter) Post(c Context) (statusCode int, err error) {
	log.Info("execute baseFilter Post function")
	return http.StatusOK, nil
}

// PostErr execute proxy has errors
func (f BaseFilter) PostErr(c Context, code int, err error) {
	log.Infof("execute baseFilter PostErr function,code: %d error: %v", code, err)
}
