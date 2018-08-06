package vkapi

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

var (
	// GroupAccessTokenReg - регуларка для получения id группы для которой получен токен
	GroupAccessTokenReg *regexp.Regexp
	// ObjLinkReg - регулярка для получения владельца и ид по ссылке
	ObjLinkReg *regexp.Regexp

	linkScreenNameReg *regexp.Regexp
)

func init() {
	GroupAccessTokenReg = regexp.MustCompile("^access_token_([0-9]+)$")
	linkScreenNameReg = regexp.MustCompile("vk.com/(.+)")
	ObjLinkReg = regexp.MustCompile("(?:wall|page|topic|photo|album|video|product)-?([0-9]+)_([0-9]+)")
}

// Разбиваем массив строк на несколько
func chunkSliceString(arr []string, size int) (ans [][]string) {
	msize := len(arr) / size
	if len(arr)%size != 0 {
		msize++
	}
	ans = make([][]string, msize)

	l := len(arr)
	now := 0
	i := 0
	for {
		next := now + size

		if now+size > l {
			next = l
		}

		ans[i] = arr[now:next]
		i++

		if next == l {
			break
		}
		now = next
	}

	return
}

// Проверяем надо ли пропустить ошибу
func (vk *API) checkErrorSkip(str string) bool {
	for _, e := range vk.ErrorToSkip {
		if strings.Contains(str, e) {
			return true
		}
	}

	return false
}

// StopAllQuery - Останавливаем все запросы
func StopAllQuery() {
	exited = true
	contMap.Lock()
	for _, f := range contMap.h {
		f()
	}
	contMap.Unlock()
}

// GetCriteriaJSON - формируем json для критериев тергетинга
func GetCriteriaJSON(d AdsGetTargetingStatsCriteria) (b []byte) {
	b, err := json.Marshal(d)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	if d.GeoPointType == "" {
		b = []byte(strings.Replace(string(b), `"geo_point_type":"",`, "", -1))
	}
	return
}

// GetScreenNameFromLink - получаем screen_name из ссылке
func GetScreenNameFromLink(link string) (screenName string) {
	f := linkScreenNameReg.FindStringSubmatch(link)
	if len(f) == 0 {
		return
	}

	screenName = f[1]
	return
}
