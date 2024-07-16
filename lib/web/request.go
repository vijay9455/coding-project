package web

import (
	"calendly/lib/validator"
	"context"
	"encoding/json"
	"net/http"
)

type Request struct {
	*http.Request

	pathParams map[string]string
}

type Handle func(request *Request) (*JSONResponse, ErrorInterface)

func NewRequest(r *http.Request) Request {
	return Request{
		Request: r.WithContext(SetRequestID(r.Context(), GetRequestIDFromRequestHeader(r))),
	}
}

func (req *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}
	newHttpRequest := req.Request.WithContext(ctx)
	newReq := new(Request)
	*newReq = *req
	newReq.Request = newHttpRequest

	return newReq
}

func (req *Request) SetPathParam(key, value string) {
	req.pathParams[key] = value
}

func (r *Request) QueryParam(key string) string {
	return r.URL.Query().Get(key)
}

func (req *Request) PathParams() map[string]string {
	return req.pathParams
}

func (req *Request) ValidateBodyToStruct(s any) error {
	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&s)
	if err != nil {
		return validator.HandleValidationErrors(err)
	}

	return validator.ValidateStruct(s)
}

func (req *Request) ValidateQueryParamsToStruct(s any) error {
	jsonString, _ := json.Marshal(req.queryMap())
	err := json.Unmarshal(jsonString, &s)
	if err != nil {
		return validator.HandleValidationErrors(err)
	}
	return validator.ValidateStruct(s)
}

func (req *Request) queryMap() map[string]string {
	queryMap := make(map[string]string)
	for key, values := range req.URL.Query() {
		queryMap[key] = values[0]
	}
	return queryMap
}
