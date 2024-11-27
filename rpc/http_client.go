package rpc

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type (
	RequestType string
	ContentType string
	Params      map[string]string
	Headers     map[string]string
)

const (
	Authorization               = "Authorization"
	FormContentType ContentType = "application/x-www-form-urlencoded; charset=utf-8"
	JsonContentType ContentType = "application/json; charset=utf-8"
)

const (
	Get      RequestType = "GET"
	PostForm RequestType = "POST_FORM"
	PostJson RequestType = "POST_JSON"
	Put      RequestType = "PUT"
	Patch    RequestType = "PATCH"
)

type HttpReqDTO struct {
	*Builder
}

type Builder struct {
	RequestType RequestType
	Headers     Headers
	BaseUrl     string
	Url         string
	Params      Params
	Data        interface{}
	IsPrintLog  bool
}

func NewHttpClientBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) SetRequestType(requestType RequestType) *Builder {
	b.RequestType = requestType
	return b
}

func (b *Builder) SetHeaders(headers Headers) *Builder {
	b.Headers = headers
	return b
}

func (b *Builder) SetBaseUrl(url string) *Builder {
	b.BaseUrl = url
	return b
}

func (b *Builder) SetUrl(url string) *Builder {
	b.Url = url
	return b
}

func (b *Builder) SetParams(params Params) *Builder {
	b.Params = params
	return b
}

func (b *Builder) SetPrintLog(isPrintLog bool) *Builder {
	b.IsPrintLog = isPrintLog
	return b
}

func (b *Builder) SetData(data interface{}) *Builder {
	b.Data = data
	return b
}

func (b *Builder) Build() HttpReqDTO {
	return HttpReqDTO{Builder: b}
}

type Http interface {
	GetClient() *http.Client
	Sync(ctx context.Context, reqDTO HttpReqDTO, timeout time.Duration) (data string, err error)
	Call(ctx context.Context, reqDTO HttpReqDTO, timeout time.Duration) (statusCode int, data string, err error)
}

func buildFormParams(params Params) url.Values {
	var values = url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values
}

func buildHeaders(httpReq *http.Request, headers Headers) {
	for k, v := range headers {
		httpReq.Header.Add(k, v)
	}
}
