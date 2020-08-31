package proxy

import (
	"net/http"
	"time"
)

var httpClient = http.Client{
	Timeout: time.Minute,
}
