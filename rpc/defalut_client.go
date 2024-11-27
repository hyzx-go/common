package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyzx-go/common-b2c/log"
	"github.com/hyzx-go/common-b2c/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type httpClient struct {
	cli *http.Client
}

// NewHttpClient Create a Http, Cannot support golang init()
func NewHttpClient(cli *http.Client) Http {
	return &httpClient{cli: cli}
}

// GetClient get origin http client
func (h *httpClient) GetClient() *http.Client {
	return h.cli
}

// Sync call api
func (h *httpClient) Sync(ctx context.Context, reqDTO HttpReqDTO, timeout time.Duration) (data string, err error) {
	if reqDTO.IsPrintLog {
		log.Ctx(ctx).Info(utils.Format("httpClient pending invoke params -> %v", reqDTO))
	}
	statusCode, data, err := h.Call(ctx, reqDTO, timeout)
	if reqDTO.IsPrintLog {
		log.Ctx(ctx).Info(utils.Format("httpClient invoke success -> statusCode ->%v data-> %v error->%v", statusCode, data, err))
	}
	return data, err
}

// Call if u have any question please call me - xiaomin.tai@carsome.com
func (h *httpClient) Call(ctx context.Context, reqDTO HttpReqDTO, timeout time.Duration) (statusCode int, data string, err error) {

	if reqDTO.BaseUrl != "" {
		reqDTO.Url = reqDTO.BaseUrl + reqDTO.Url
	}

	if reqDTO.Url == "" || reqDTO.RequestType == "" {
		return 0, data, errors.New("HttpRequest params cannot be empty ")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var httpReqDTO *http.Request
	switch reqDTO.RequestType {
	case Get:

		httpReqDTO, err = http.NewRequest("GET", reqDTO.Url, nil)
		if err != nil {
			break
		}
		if len(reqDTO.Params) > 0 {
			query := httpReqDTO.URL.Query()
			for k, v := range reqDTO.Params {
				query.Set(k, v)
			}
			httpReqDTO.URL.RawQuery = query.Encode()
		}
		break

	case PostJson:
		var body []byte
		if val, ok := reqDTO.Data.(*bytes.Buffer); ok {
			body = val.Bytes()
		} else {
			body, err = json.Marshal(reqDTO.Data)
			if err != nil {
				break
			}
		}

		httpReqDTO, err = http.NewRequest("POST", reqDTO.Url, bytes.NewReader(body))
		if err != nil {
			break
		}
		httpReqDTO.Header.Add("Content-Type", string(JsonContentType))
		if len(reqDTO.Params) > 0 {
			query := httpReqDTO.URL.Query()
			for k, v := range reqDTO.Params {
				query.Set(k, v)
			}
			httpReqDTO.URL.RawQuery = query.Encode()
		}
		break

	case PostForm:
		var body []byte
		if val, ok := reqDTO.Data.(url.Values); ok {
			body = []byte(val.Encode())
		}
		httpReqDTO, err = http.NewRequest("POST", reqDTO.Url, bytes.NewReader(body))
		if err != nil {
			break
		}
		httpReqDTO.Header.Add("Content-Type", string(FormContentType))
		httpReqDTO.PostForm = buildFormParams(reqDTO.Params)
		if len(reqDTO.Params) > 0 {
			query := httpReqDTO.URL.Query()
			for k, v := range reqDTO.Params {
				query.Set(k, v)
			}
			httpReqDTO.URL.RawQuery = query.Encode()
		}
		break

	case Put:

		var body []byte
		if reqDTO.Data != nil {
			body, err = json.Marshal(reqDTO.Data)
			if err != nil {
				break
			}
		}

		httpReqDTO, err = http.NewRequest("PUT", reqDTO.Url, bytes.NewReader(body))
		if err != nil {
			break
		}
		httpReqDTO.Header.Add("Content-Type", string(JsonContentType))
		break

	case Patch:
		var body []byte
		body, err = json.Marshal(reqDTO.Data)
		if err != nil {
			break
		}

		httpReqDTO, err = http.NewRequest("PATCH", reqDTO.Url, bytes.NewReader(body))
		if err != nil {
			break
		}
		httpReqDTO.Header.Add("Content-Type", string(JsonContentType))

		break
	default:
		return 0, data, errors.New("HttpRequest type cannot support ")
	}

	if err != nil {
		return 0, data, fmt.Errorf("HttpRequest create params error, %w", err)
	}

	buildHeaders(httpReqDTO, reqDTO.Headers)

	httpRes, err := h.cli.Do(httpReqDTO.WithContext(ctx))
	if err != nil {
		return http.StatusInternalServerError, data, fmt.Errorf("sync request error, %w", err)
	}

	defer func() {
		if err := httpRes.Body.Close(); err != nil {
			log.Ctx(ctx).Error(fmt.Sprintf("close http body failed, %v", err))
		}
	}()

	if !(httpRes.StatusCode >= http.StatusOK && httpRes.StatusCode <= http.StatusIMUsed) {
		return httpRes.StatusCode, data, fmt.Errorf("sync request error ,url -> %s response code -> %s",
			reqDTO.Url, strconv.Itoa(httpRes.StatusCode))
	}

	dataByte, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return httpRes.StatusCode, data, fmt.Errorf("cannot get future result , url -> %s, %w", reqDTO.Url, err)
	}

	return httpRes.StatusCode, utils.ToString(dataByte), nil
}
