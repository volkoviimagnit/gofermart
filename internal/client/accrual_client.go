package client

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/volkoviimagnit/gofermart/internal/client/response"
)

type AccrualHttpClient struct {
	host string
}

func NewAccrualHttpClient(host string) IAccrualClient {
	return &AccrualHttpClient{host: host}
}

func (a *AccrualHttpClient) GetDefaultRetryAfterSeconds() time.Duration {
	return 60
}

func (a *AccrualHttpClient) GetOrderStatus(orderId string) (IAccrualOrderStatus, IError) {
	client := resty.New()
	client.SetRetryCount(3).
		SetRetryWaitTime(10*time.Second).
		SetRetryMaxWaitTime(20*time.Second).
		SetHeader("Accept-Encoding", "gzip,deflate")

	var body []byte
	var url = a.host + "/api/orders/" + orderId
	restyResponse, respError := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Get(url)

	if respError != nil {
		return nil, &UndefinedError{err: respError}
	}

	logrus.Warningf("%s %d", orderId, restyResponse.StatusCode())

	switch restyResponse.StatusCode() {
	case http.StatusOK:
		responseBody := restyResponse.Body()

		orderDTO := response.OrderDTO{}
		errDecode := json.Unmarshal(responseBody, &orderDTO)
		if errDecode != nil {
			return &orderDTO, &UndefinedError{err: errDecode}
		}
		return &orderDTO, nil
	case http.StatusNoContent:
		return nil, &StatusNoContentError{}
	case http.StatusTooManyRequests:

		headerRetryAfter, err := strconv.Atoi(restyResponse.Header().Get("Retry-After"))
		var retryAfter time.Duration
		if err != nil {
			retryAfter = a.GetDefaultRetryAfterSeconds()
		} else {
			retryAfter = time.Duration(headerRetryAfter)
		}
		return nil, &StatusTooManyRequestsError{
			retryAfterSeconds: retryAfter,
		}
	default:
		return nil, &InternalServerError{
			retryAfterSeconds: a.GetDefaultRetryAfterSeconds(),
		}
	}

}

type InternalServerError struct {
	retryAfterSeconds time.Duration
}

func (e *InternalServerError) NeedRetry() bool {
	return e.retryAfterSeconds > 0
}

func (e *InternalServerError) RetryAfterSeconds() time.Duration {
	return e.retryAfterSeconds
}

func (e *InternalServerError) Error() string {
	return "сервер недоступен"
}

type StatusNoContentError struct{}

func (e *StatusNoContentError) Error() string {
	return "заказ не существует в системе"
}

func (e *StatusNoContentError) NeedRetry() bool {
	return false
}

func (e *StatusNoContentError) RetryAfterSeconds() time.Duration {
	return 0
}

type StatusTooManyRequestsError struct {
	retryAfterSeconds time.Duration
}

func (e *StatusTooManyRequestsError) NeedRetry() bool {
	return e.retryAfterSeconds > 0
}

func (e *StatusTooManyRequestsError) RetryAfterSeconds() time.Duration {
	return e.retryAfterSeconds
}

func (e *StatusTooManyRequestsError) Error() string {
	return "no more than N requests per minute allowed"
}

type UndefinedError struct {
	err error
}

func (e *UndefinedError) NeedRetry() bool {
	return false
}

func (e *UndefinedError) RetryAfterSeconds() time.Duration {
	return 0
}

func (e *UndefinedError) Error() string {
	return "undefined error - " + e.err.Error()
}
