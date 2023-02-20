package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func (r *Response) Success(pay Payments) (int, *Response) {
	r.Code = http.StatusOK
	r.Message = "payment success"
	r.Payload = pay
	return r.Code, r
}

func (r *Response) BadRequest(err error) (int, *Response) {
	r.Code = http.StatusBadRequest
	r.Message = err.Error()
	return r.Code, r
}

func (r *Response) NotFound(err error) (int, *Response) {
	r.Code = http.StatusNotFound
	r.Message = err.Error()
	return http.StatusOK, r
}

func (r *Response) Inactive(err error) (int, *Response) {
	r.Code = 0
	r.Message = err.Error()
	return http.StatusOK, r
}

func (r *Response) Found(acc_id int, user uuid.UUID) (int, *Response) {
	r.Code = http.StatusFound
	r.Message = "account exist"
	r.Payload = gin.H{"account": acc_id, "user": user}
	return http.StatusOK, r
}

func (r *Response) NotAllowed(err error) (int, *Response) {
	r.Code = 0
	r.Message = err.Error()
	return http.StatusConflict, r
}

func (r *Response) TransactionError(err error) (int, *Response) {
	r.Code = http.StatusInternalServerError
	r.Message = err.Error()
	return r.Code, r
}
