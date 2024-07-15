package middleware

import (
	"calendly/lib/logger"
	"calendly/lib/web"
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/julienschmidt/httprouter"
)

type ApiVersion string

const (
	V1 ApiVersion = "1.0.0"
)

func ServeEndpoint(apiVersion ApiVersion, nextHandler web.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, pathParams httprouter.Params) {
		startTime := time.Now()
		var webReq web.Request

		defer func() {
			if recvr := recover(); recvr != nil {
				logger.Error(webReq.Context(), "Request Failed", map[string]any{
					"error": fmt.Sprintf("%v", recvr), "status": http.StatusInternalServerError,
					"path": r.URL.Path, "request_params": r.URL.Query(), "method": r.Method,
					"stack": string(debug.Stack()), "duration_ms": float64(time.Since(startTime).Nanoseconds()) / 1e6,
				})
				writeResponse(webReq.Context(), w, apiVersion, nil, web.ErrInternalServerError(fmt.Sprintf("%v", recvr)))
			}
		}()

		webReq = web.NewRequest(r)
		for i := range pathParams {
			webReq.SetPathParam(pathParams[i].Key, pathParams[i].Value)
		}
		data, responseErr := nextHandler(&webReq)
		writeResponse(webReq.Context(), w, apiVersion, data, responseErr)
		logRequest(&webReq, startTime, data, responseErr)
	}
}

func getResponseCode(webErr web.ErrorInterface) int {
	if webErr == nil {
		return http.StatusOK
	}
	return webErr.HttpStatusCode()
}

func writeResponse(ctx context.Context, w http.ResponseWriter, apiVersion ApiVersion, responseData *web.JSONResponse, responseErr web.ErrorInterface) {
	w.WriteHeader(getResponseCode(responseErr))
	_, err := w.Write(buildResponseBody(apiVersion, responseData, responseErr).ByteArray())
	if err != nil {
		logger.Error(ctx, "error in writing response", map[string]any{"error": err})
	}
}

func buildResponseBody(apiVersion ApiVersion, data *web.JSONResponse, responseErr web.ErrorInterface) *web.JSONResponse {
	if responseErr == nil {
		return &web.JSONResponse{
			"api_version": apiVersion,
			"success":     true,
			"data":        data,
		}
	}

	return &web.JSONResponse{
		"api_version": apiVersion,
		"success":     false,
		"error": &web.JSONResponse{
			"code":    responseErr.Code(),
			"message": responseErr.Description(),
		},
	}
}

func logRequest(req *web.Request, startTime time.Time, _ *web.JSONResponse, responseErr web.ErrorInterface) {
	logMap := map[string]any{
		"status":         getResponseCode(responseErr),
		"path":           req.URL.Path,
		"request_params": req.URL.Query(),
		"method":         req.Method,
		"duration_ms":    float64(time.Since(startTime).Nanoseconds()) / 1e6,
	}
	if responseErr != nil {
		logMap["error"] = map[string]any{"code": responseErr.Code(), "description": responseErr.Description()}
	}
	logger.Info(req.Context(), "Request processed", logMap)
}
