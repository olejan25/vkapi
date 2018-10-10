package vkapi

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/fe0b6/tools"
)

var (
	// GroupAccessTokenReg - регуларка для получения id группы для которой получен токен
	GroupAccessTokenReg *regexp.Regexp
	// ObjLinkReg - регулярка для получения владельца и ид по ссылке
	ObjLinkReg *regexp.Regexp
	// InternalLinkReg - регулярка для получения данных из внутренних ссылок
	InternalLinkReg *regexp.Regexp
	// LinkDomainReg - регулярка для получения домена
	LinkDomainReg *regexp.Regexp

	linkScreenNameReg *regexp.Regexp
)

func init() {
	GroupAccessTokenReg = regexp.MustCompile("^access_token_([0-9]+)$")
	linkScreenNameReg = regexp.MustCompile("vk.com/(.+)")
	LinkDomainReg = regexp.MustCompile("vk.com/([^? \n]+)\\S*")
	ObjLinkReg = regexp.MustCompile("(?:wall|page|topic|photo|album|video|product|market)(-)?([0-9]+)(?:_([0-9]+))?")
	InternalLinkReg = regexp.MustCompile("\\[(id|club|public|group)([0-9]+)\\|[^\\]]+\\]")
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

// EncryptToken - шифруем токен
func EncryptToken(key string, token string) (encToken string, err error) {
	b, err := tools.AESEncrypt([]byte(key), []byte(token))
	if err != nil {
		log.Println("[error]", err)
		return
	}

	encToken = hex.EncodeToString(b)
	return
}

// DecryptToken - расшифровываем токен
func DecryptToken(key string, encToken string) (token string, err error) {
	bToken, err := hex.DecodeString(encToken)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	b, err := tools.AESDecrypt([]byte(key), bToken)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	token = string(b)
	return
}
