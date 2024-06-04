package component

import (
	"github.com/go-resty/resty/v2"
)

func NewRestyClient() *resty.Client {
	return resty.New()
}
