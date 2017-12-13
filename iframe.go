package vkapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"regexp"
)

var (
	iframeCheckSignReq *regexp.Regexp
)

func init() {
	iframeCheckSignReq = regexp.MustCompile("(?:^|&)([a-z0-9_]+)=")
}

// Проверка подписи ВК
func IframeCheckSign(r *http.Request, secret string) (ok bool) {
	var sign string
	/*
		Параметры надо отсортировывать в той же последовательности что они переданы
	*/
	for _, v := range iframeCheckSignReq.FindAllStringSubmatch(r.URL.RawQuery, -1) {
		if v[1] == "hash" || v[1] == "sign" || v[1] == "api_result" {
			continue
		}

		sign += r.FormValue(v[1])
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(sign))

	if r.FormValue("sign") != hex.EncodeToString(h.Sum(nil)) {
		return
	}

	ok = true
	return
}
