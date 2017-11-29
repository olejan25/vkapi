package vkapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	API_VERSION    = "5.69"
	API_METHOD_URL = "https://api.vk.com/method/"
	API_TOKEN_URL  = "https://oauth.vk.com/access_token"
)

var (
	httpTr *http.Transport
)

func init() {
	httpTr = &http.Transport{
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     10 * time.Minute,
	}
}

/*
	Получение токена
*/

// Получение токена
func (vk *Api) GetToken(d TokenData) (ans GetTokenAns, err error) {
	q := url.Values{}
	q.Add("code", d.Code)
	q.Add("client_id", strconv.Itoa(d.ClientId))
	q.Add("client_secret", d.ClientSecret)
	q.Add("redirect_uri", d.RedirectUri)
	q.Add("v", API_VERSION)

	// Формируем запрос
	req, err := http.NewRequest("POST", API_TOKEN_URL, strings.NewReader(q.Encode()))
	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Если статус ответа не правильный
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		log.Println("[error]", resp.Status, resp.StatusCode)
		return
	}

	// Читаем ответ
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Парсим ответ
	err = json.Unmarshal(content, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

/*
	Users
*/

// Получаем информацию о пользователях
func (vk *Api) Users_get(params map[string]string) (ans []UsersGetAns, err error) {

	// Отправляем запрос
	r, err := vk.request("users.get", params)
	if err != nil {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

/*
	Groups
*/

// Получаем информацию о пользователях
func (vk *Api) Groups_getById(params map[string]string) (ans []GroupsGetAns, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.getById", params)
	if err != nil {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

/*
	Wall
*/

// Возвращает список записей со стен пользователей или сообществ по их идентификаторам.
func (vk *Api) Wall_get(params map[string]string) (ans WallGetAns, err error) {

	// Отправляем запрос
	r, err := vk.request("wall.get", params)
	if err != nil {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

// Возвращает список записей со стен пользователей или сообществ по их идентификаторам.
func (vk *Api) Wall_getById(params map[string]string) (ans []WallGetByIdAns, err error) {

	// Отправляем запрос
	r, err := vk.request("wall.getById", params)
	if err != nil {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

/*
	Utils
*/

// Получаем сокращенную ссылку
func (vk *Api) Utils_getShortLink(params map[string]string) (ans UtilsGetShortLinkAns, err error) {

	// Отправляем запрос
	r, err := vk.request("utils.getShortLink", params)
	if err != nil {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

// Получаем статистику по ссылке
func (vk *Api) Utils_getLinkStats(params map[string]string) (ans UtilsGetLinkStatsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("utils.getLinkStats", params)
	if err != nil {
		return
	}

	// Если ответ пустой
	if string(r.Response) == "[]" {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

/*
	Ads
*/

// Получаем список кампаний
func (vk *Api) Ads_getCampaigns(params map[string]string) (ans []AdsGetCampaignsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.getCampaigns", params)
	if err != nil {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

// Получаем список список объявлений
func (vk *Api) Ads_getAdsLayout(params map[string]string) (ans []AdsGetAdsLayoutAns, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.getAdsLayout", params)
	if err != nil {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

// Получаем статистику объявлений
func (vk *Api) Ads_getStatistics(params map[string]string) (ans []AdsGetStatisticsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.getStatistics", params)
	if err != nil {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err, string(r.Response))
		return
	}

	for i, d := range ans {
		// Создаем массив для норм статы
		ans[i].Stats = make([]AdsGetStatisticsAnsStats, len(d.StatsBug))

		for k, s := range d.StatsBug {
			// Пробуем норм разобрать
			var t AdsGetStatisticsAnsStats
			err = json.Unmarshal(s, &t)
			// Если ошибка пробуем разобрать кривой json
			if err != nil {
				var t2 AdsGetStatisticsAnsStatsBug
				err = json.Unmarshal(s, &t2)
				if err != nil {
					log.Println("[error]", err, string(r.Response))
					return
				}

				impr, _ := strconv.ParseInt(t2.Impressions, 10, 32)

				t = AdsGetStatisticsAnsStats{
					Day:         t2.Day,
					Spent:       t2.Spent,
					Clicks:      t2.Clicks,
					Reach:       t2.Reach,
					Impressions: int(impr),
				}
			}

			ans[i].Stats[k] = t
		}
	}

	return
}

/*
	Execute
*/

// Execute
func (vk *Api) Execute(code string) (r Response, err error) {

	// Отправляем запрос
	r, err = vk.request("execute", map[string]string{"code": code})
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			log.Println("[error]", err)
			log.Println(code)
		}
		return
	}

	if len(r.ExecuteErrors) > 0 {
		log.Println(code)
		log.Println(r.ExecuteErrors)
	}

	return
}

/*
	Запрос к ВК
*/

// Обертка для запроса к ВК
func (vk *Api) request(method string, params map[string]string) (ans Response, err error) {
	if vk.AccessToken == "" {
		err = errors.New("no access token")
		log.Println("[error]", err)
		return
	}

	for {
		ans, err = vk.fullRequest(method, params)
		if err != nil {
			if strings.Contains(err.Error(), "server sent GOAWAY") ||
				strings.Contains(err.Error(), "unexpected EOF") {
				if vk.httpErrorWait(method) {
					continue
				}
			}
			return
		}

		// Проверяем ответ
		if ans.Error.ErrorCode != 0 {
			if ans.Error.ErrorMsg == "Too many requests per second" {
				// Ждем между запросами
				if vk.floodWait(method) {
					continue
				}
			} else if ans.Error.ErrorMsg == "Runtime error occurred during code invocation: Comparing values of different or unsupported types" {
				log.Println("[error]", params["code"])
			}

			err = errors.New(ans.Error.ErrorMsg)
			return
		}

		break
	}

	return
}

// Запрос к ВК
func (vk *Api) fullRequest(method string, params map[string]string) (ans Response, err error) {
	q := url.Values{}
	for k, v := range params {
		q.Add(k, v)
	}
	if params["v"] == "" {
		q.Add("v", API_VERSION)
	}
	q.Add("access_token", vk.AccessToken)

	// Формируем запрос
	req, err := http.NewRequest("POST", API_METHOD_URL+method, strings.NewReader(q.Encode()))
	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Отправляем запрос
	client := &http.Client{Transport: httpTr}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Если проблема с ответом
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		log.Println("[error]", resp.Status, resp.StatusCode)
		return
	}

	// Читаем ответ
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if !strings.Contains(err.Error(), "server sent GOAWAY") &&
			!strings.Contains(err.Error(), "unexpected EOF") {
			log.Println("[error]", err)
		}
		return
	}

	// Парсим ответ
	err = json.Unmarshal(content, &ans)
	if err != nil {
		log.Println("[error]", err, string(content))
		return
	}

	return
}

// Ждем между запросами если вк ответил что запросы слишком частые
func (vk *Api) floodWait(method string) (ok bool) {
	// Для некоторых методов повторять бессмысленно (https://vk.com/dev/data_limits)
	if vk.retryCount >= 3 &&
		(method == "wall.get" || method == "wall.search" || method == "newsfeed.search") {
		// Сбрасываем счетчик ожидания
		vk.Lock()
		vk.retryCount = 0
		vk.Unlock()
		return
	}

	// Определяем сколько времени будет ждать
	var sleepTime int
	if vk.retryCount < 5 {
		sleepTime = 1
	} else if vk.retryCount < 10 {
		sleepTime = 2
	} else if vk.retryCount < 20 {
		sleepTime = 3
	} else if vk.retryCount < 25 {
		sleepTime = 5
	} else {
		// Сбрасываем счетчик ожидания
		vk.Lock()
		vk.retryCount = 0
		vk.Unlock()
		return
	}

	// Увеличиваем счетчик
	vk.Lock()
	vk.retryCount++
	vk.Unlock()

	// Ждем
	time.Sleep(time.Duration(sleepTime) * time.Second)

	ok = true
	return
}

// Попытка повтора запроса при ошибки http
func (vk *Api) httpErrorWait(method string) (ok bool) {
	if method == "wall.post" || method == "wall.repost" {
		return
	}

	if vk.httpRetryCount >= 3 {
		vk.Lock()
		vk.httpRetryCount = 0
		vk.Unlock()
		return
	}

	vk.Lock()
	vk.httpRetryCount++
	vk.Unlock()

	// Ждем
	time.Sleep(1 * time.Second)

	ok = true
	return
}
