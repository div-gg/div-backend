package models

type (
  Response struct {
    Data interface{} `json:"data,omitempty"`
    Message string `json:"message,omitempty"`
    Status int `json:"status,omitempty"`
  }
)
