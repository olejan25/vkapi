package vkapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// APIVersion - используемая версия API
	APIVersion = "5.92"
	// APIMethodURL - URL запросов к API
	APIMethodURL = "https://api.vk.com/method/"
	// APITokenURL - URL oauth авторизации
	APITokenURL = "https://oauth.vk.com/access_token"
	// APIAuthURL - URL для oauth авторизации
	APIAuthURL = "https://oauth.vk.com/authorize"
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

// GetAuthURL - получаем ссылку для авторизации
func GetAuthURL(d AuthURLData) string {
	str := APIAuthURL + fmt.Sprintf("?client_id=%d&redirect_uri=%s&response_type=code", d.ClientID, d.RedirectURI)
	if d.V != 0 {
		str += fmt.Sprintf("&v=%g", d.V)
	} else {
		str += "&v=" + APIVersion
	}

	if d.Scope != "" {
		str += "&scope=" + d.Scope
	}
	if d.GroupIDs != "" {
		str += "&group_ids=" + d.GroupIDs
	}
	if d.Display != "" {
		str += "&display=" + d.Display
	}
	if d.State != "" {
		str += "&state=" + d.State
	}

	return str
}

// GetTokenGroup - Получение токена группы
func (vk *API) GetTokenGroup(d TokenData) (ans map[string]interface{}, err error) {
	content, err := getToken(d)
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

// GetToken - Получение токена
func (vk *API) GetToken(d TokenData) (ans GetTokenAns, err error) {
	content, err := getToken(d)
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

func getToken(d TokenData) (content []byte, err error) {
	q := url.Values{}
	q.Add("code", d.Code)
	q.Add("client_id", strconv.Itoa(d.ClientID))
	q.Add("client_secret", d.ClientSecret)
	q.Add("redirect_uri", d.RedirectURI)
	q.Add("v", APIVersion)

	// Формируем запрос
	req, err := http.NewRequest("POST", APITokenURL, strings.NewReader(q.Encode()))
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
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

/*
	Users
*/

// UsersGet - Получаем информацию о пользователях
func (vk *API) UsersGet(params map[string]string) (ans []UsersGetAns, err error) {

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

// UsersGetSubscriptions - Получаем информацию о пользователях
func (vk *API) UsersGetSubscriptions(params map[string]string) (ans UsersGetSubscriptionsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("users.getSubscriptions", params)
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

// GroupsJoin - Присоединяемся к группе
func (vk *API) GroupsJoin(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.join", params)
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

// GroupsGet - Получаем информацию о группах
func (vk *API) GroupsGet(params map[string]string) (ans GroupsGetAns, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.get", params)
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

// GroupsGetByID - Получаем информацию о группах
func (vk *API) GroupsGetByID(params map[string]string) (ans []GroupsGetByIDAns, err error) {

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

// GroupsGetMembers - Получаем информацию о подписчиках
func (vk *API) GroupsGetMembers(params map[string]string) (ans GroupsGetMembersAns, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.getMembers", params)
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

// GroupsIsMember - Получаем информацию о подписчиках
// При запросе нескольких человек одновременно результат может быть не верным. баг ВК
func (vk *API) GroupsIsMember(params map[string]string) (ans []GroupsIsMemberAns, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.isMember", params)
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

// GroupsIsMemberOne - Получаем информацию о подписчиках
// При запросе нескольких человек одновременно результат может быть не верным. баг ВК
func (vk *API) GroupsIsMemberOne(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.isMember", params)
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

// GroupsGetTokenPermissions - Получаем информацию о правах токена
func (vk *API) GroupsGetTokenPermissions() (ans GroupsGetTokenPermissionsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.getTokenPermissions", map[string]string{})
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

// GroupsGetCallbackServers - Получаем информацию о callback серверах
func (vk *API) GroupsGetCallbackServers(params map[string]string) (ans GroupsGetCallbackServersAns, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.getCallbackServers", params)
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

// GroupsGetCallbackSettings - Получаем настройки callback сервера
func (vk *API) GroupsGetCallbackSettings(params map[string]string) (ans GroupsGetCallbackSettingsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.getCallbackSettings", params)
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

// GroupsAddCallbackServer - Добавляем callback сервер
func (vk *API) GroupsAddCallbackServer(params map[string]string) (ans GroupsAddCallbackServerAns, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.addCallbackServer", params)
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

// GroupsEditCallbackServer - редактирование callback сервер
func (vk *API) GroupsEditCallbackServer(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.editCallbackServer", params)
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

// GroupsDeleteCallbackServer - удаляем callback сервер
func (vk *API) GroupsDeleteCallbackServer(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.deleteCallbackServer", params)
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

// GroupsSetCallbackSettings - настройка callback сервер
func (vk *API) GroupsSetCallbackSettings(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.setCallbackSettings", params)
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

// GroupsGetCallbackConfirmationCode - Получаем код подтверждения для сервера callback
func (vk *API) GroupsGetCallbackConfirmationCode(params map[string]string) (ans GroupsGetCallbackConfirmationCodeAns, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.getCallbackConfirmationCode", params)
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

// GroupsBan - баним в сообществе
func (vk *API) GroupsBan(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.ban", params)
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

// GroupsGetBanned - Получаем инфу по забаненым
func (vk *API) GroupsGetBanned(params map[string]string) (ans GroupsGetBannedAns, err error) {

	// Отправляем запрос
	r, err := vk.request("groups.getBanned", params)
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

// WallGet - Возвращает список записей со стен пользователей или сообществ по их идентификаторам.
func (vk *API) WallGet(params map[string]string) (ans WallGetAns, err error) {

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

// WallGetByID - Возвращает список записей со стен пользователей или сообществ по их идентификаторам.
func (vk *API) WallGetByID(params map[string]string) (ans []WallGetByIDAns, err error) {

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

// WallGetComment - Возвращает список комментариев к посту.
func (vk *API) WallGetComment(params map[string]string) (ans WallGetCommentsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("wall.getComment", params)
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

// WallGetComments - Возвращает список комментариев к посту.
func (vk *API) WallGetComments(params map[string]string) (ans WallGetCommentsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("wall.getComments", params)
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

// WallDelete - Удаляем пост со стены
func (vk *API) WallDelete(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("wall.delete", params)
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

// WallRestore - Восстанавливаем пост на стене
func (vk *API) WallRestore(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("wall.restore", params)
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

// WallDeleteComment - Удаляем комментарий со стены
func (vk *API) WallDeleteComment(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("wall.deleteComment", params)
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

// WallRestoreComment - Восстанавливаем комментарий на стене
func (vk *API) WallRestoreComment(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("wall.restoreComment", params)
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
	Likes
*/

// LikesGetList - Возвращает список лайков.
func (vk *API) LikesGetList(params map[string]string) (ans LikesGetListAns, err error) {

	// Отправляем запрос
	r, err := vk.request("likes.getList", params)
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
	Board
*/

// BoardGetTopics - Возвращает список обсуждений.
func (vk *API) BoardGetTopics(params map[string]string) (ans BoardGetTopicsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("board.getTopics", params)
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

// BoardGetComments - Возвращает список комментариев обсуждения.
func (vk *API) BoardGetComments(params map[string]string) (ans BoardGetCommentsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("board.getComments", params)
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

// BoardDeleteComment - Удаляем комментарий из обсуждения
func (vk *API) BoardDeleteComment(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("board.deleteComment", params)
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

// BoardRestoreComment - восстанавливаем комментарий из обсуждения
func (vk *API) BoardRestoreComment(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("board.restoreComment", params)
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
	Photos
*/

// PhotosGetAlbums - Возвращает список видео.
func (vk *API) PhotosGetAlbums(params map[string]string) (ans PhotosGetAlbumsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("photos.getAlbums", params)
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

// PhotosGet - Возвращает список фотографий.
func (vk *API) PhotosGet(params map[string]string) (ans PhotosGetAns, err error) {

	// Отправляем запрос
	r, err := vk.request("photos.get", params)
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

// PhotosGetAll - Возвращает список фотографий.
func (vk *API) PhotosGetAll(params map[string]string) (ans PhotosGetAns, err error) {

	// Отправляем запрос
	r, err := vk.request("photos.getAll", params)
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

// PhotosGetByID - Возвращает список фотографий.
func (vk *API) PhotosGetByID(params map[string]string) (ans []PhotosGetItem, err error) {

	// Отправляем запрос
	r, err := vk.request("photos.getById", params)
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

// PhotosGetComments - Возвращает список комментариев фотографии.
func (vk *API) PhotosGetComments(params map[string]string) (ans PhotosGetCommentsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("photos.getComments", params)
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

// PhotosGetAllComments - Возвращает список комментариев фотографии.
func (vk *API) PhotosGetAllComments(params map[string]string) (ans PhotosGetCommentsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("photos.getAllComments", params)
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

// PhotosDelete - Удаление фотки
func (vk *API) PhotosDelete(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("photos.delete", params)
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

// PhotosRestore - Восстановление фотки
func (vk *API) PhotosRestore(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("photos.restore", params)
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

// PhotosDeleteComment - Удаление комментарий фотки
func (vk *API) PhotosDeleteComment(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("photos.deleteComment", params)
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

// PhotosRestoreComment - Восстановление комментарий фотки
func (vk *API) PhotosRestoreComment(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("photos.restoreComment", params)
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
	Video
*/

// VideoGet - Возвращает список видео.
func (vk *API) VideoGet(params map[string]string) (ans VideoGetAns, err error) {

	// Отправляем запрос
	r, err := vk.request("video.get", params)
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

// VideoGetComments - Возвращает список комментариев видео.
func (vk *API) VideoGetComments(params map[string]string) (ans VideoGetCommentsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("video.getComments", params)
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

// VideoDeleteComment - Удаление комментарий видео
func (vk *API) VideoDeleteComment(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("video.deleteComment", params)
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

// VideoRestoreComment - Восстановление комментарий видео
func (vk *API) VideoRestoreComment(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("video.restoreComment", params)
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
	Message
*/

// MessagesSend - отправка сообщений
func (vk *API) MessagesSend(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("messages.send", params)
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

// MessagesIsMessagesFromGroupAllowed - проверяем разрешена ли отправка сообщений от имени сообщества
func (vk *API) MessagesIsMessagesFromGroupAllowed(params map[string]string) (ans MessagesIsMessagesFromGroupAllowedAns, err error) {

	// Отправляем запрос
	r, err := vk.request("messages.isMessagesFromGroupAllowed", params)
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

// UtilsGetShortLink - Получаем сокращенную ссылку
func (vk *API) UtilsGetShortLink(params map[string]string) (ans UtilsGetShortLinkAns, err error) {

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

// UtilsGetLinkStats - Получаем статистику по ссылке
func (vk *API) UtilsGetLinkStats(params map[string]string) (ans UtilsGetLinkStatsAns, err error) {

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

// UtilsResolveScreenName - Получаем сокращенную ссылку
func (vk *API) UtilsResolveScreenName(params map[string]string) (ans UtilsResolveScreenNameAns, err error) {

	// Отправляем запрос
	r, err := vk.request("utils.resolveScreenName", params)
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
	Market
*/

// MarketGet - получаем список товаров
func (vk *API) MarketGet(params map[string]string) (ans MarketGetAns, err error) {

	// Отправляем запрос
	r, err := vk.request("market.get", params)
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

// MarketDeleteComment - удаляем комментарий у товаров
func (vk *API) MarketDeleteComment(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("market.deleteComment", params)
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

// MarketRestoreComment - восстанавливаем комментарий у товаров
func (vk *API) MarketRestoreComment(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("market.restoreComment", params)
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
	Ads
*/

// AdsGetAccounts - Получаем список аккаунтов
func (vk *API) AdsGetAccounts(params map[string]string) (ans []AdsGetAccountsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.getAccounts", params)
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

// AdsСreateTargetGroup - Создаем группу ретаргетинга
func (vk *API) AdsСreateTargetGroup(params map[string]string) (ans AdsСreateTargetGroupAns, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.createTargetGroup", params)
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

// AdsDeleteTargetGroup - удаляем группу ретаргетинга
func (vk *API) AdsDeleteTargetGroup(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.deleteTargetGroup", params)
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

// AdsImportTargetContacts - добавиление контактов в группу ретаргета
func (vk *API) AdsImportTargetContacts(params map[string]string) (ans int, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.importTargetContacts", params)
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

// AdsGetSuggestions - получение подсказок к рекламе
func (vk *API) AdsGetSuggestions(params map[string]string) (ans []AdsGetSuggestionsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.getSuggestions", params)
	if err != nil {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		// VK BUG FIX
		if strings.Contains(err.Error(), "cannot unmarshal string") {
			var arr []AdsGetSuggestionsAnsStr
			err = json.Unmarshal(r.Response, &arr)
			if err != nil {
				log.Println("[error]", err, string(r.Response))
				return
			}

			ans = []AdsGetSuggestionsAns{}
			for _, a := range arr {
				id, _ := strconv.ParseInt(a.ID, 10, 64)

				ans = append(ans, AdsGetSuggestionsAns{
					ID:     int(id),
					Name:   a.Name,
					Type:   a.Type,
					Parent: a.Parent,
				})
			}
			return
		}

		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

// AdsGetTargetGroups - получение групп ретаргета
func (vk *API) AdsGetTargetGroups(params map[string]string) (ans []AdsGetTargetGroupsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.getTargetGroups", params)
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

// AdsGetTargetingStats - Смотрим размер аудитории
func (vk *API) AdsGetTargetingStats(params map[string]string) (ans AdsGetTargetingStatsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.getTargetingStats", params)
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

// AdsGetCampaigns - Получаем список кампаний
func (vk *API) AdsGetCampaigns(params map[string]string) (ans []AdsGetCampaignsAns, err error) {

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

// AdsGetAds - Получаем список объявлений
func (vk *API) AdsGetAds(params map[string]string) (ans []AdsGetAdsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.getAds", params)
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

// AdsGetAdsLayout - Получаем список список объявлений
func (vk *API) AdsGetAdsLayout(params map[string]string) (ans []AdsGetAdsLayoutAns, err error) {

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

// AdsGetStatistics - Получаем статистику объявлений
func (vk *API) AdsGetStatistics(params map[string]string) (ans []AdsGetStatisticsAns, err error) {

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

// AdsGetDemographics - Получаем статистику объявлений демографическую
func (vk *API) AdsGetDemographics(params map[string]string) (ans []AdsGetDemographicsAns, err error) {

	// Отправляем запрос
	r, err := vk.request("ads.getDemographics", params)
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
	Stats
*/

// StatsGet - Получаем стату страницы
func (vk *API) StatsGet(params map[string]string) (ans []StatsGetAns, err error) {

	// Отправляем запрос
	r, err := vk.request("stats.get", params)
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

// StatsGetPostReach - Получаем стату поста
func (vk *API) StatsGetPostReach(params map[string]string) (ans []StatsGetPostReachAns, err error) {

	// Отправляем запрос
	r, err := vk.request("stats.getPostReach", params)
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
	Execute
*/

// Execute - пакетное выполнение запросов
func (vk *API) Execute(code string) (r Response, err error) {

	// Отправляем запрос
	r, err = vk.request("execute", map[string]string{"code": code})
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
				log.Println(code)
			}
		}
		return
	}

	if len(r.ExecuteErrors) > 0 {
		vk.ExecuteErrors = r.ExecuteErrors
		vk.ExecuteCode = code
	}

	return
}

/*
	Запрос к ВК
*/

// Обертка для запроса к ВК
func (vk *API) request(method string, params map[string]string) (ans Response, err error) {
	// прометей
	if promInited {
		promRqCount.WithLabelValues(method).Inc()
		tn := time.Now()
		defer func(tn time.Time) {
			promRq.WithLabelValues(method).Observe(float64(time.Now().Sub(tn) / time.Millisecond))
		}(tn)
	}

	if vk.AccessToken == "" {
		err = errors.New("no access token")
		log.Println("[error]", err)
		return
	}

	for {
		ans, err = vk.fullRequest(method, params)
		if err != nil {
			if httpErrorReg.MatchString(err.Error()) {
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
func (vk *API) fullRequest(method string, params map[string]string) (ans Response, err error) {
	if statRqCollect {
		// Проверим что очередь не переполнена
		if len(statRqChan) <= statRqQueueLen {
			rqStart := time.Now().UnixNano() // Время начала запроса
			defer func() {
				rqTimeout := (time.Now().UnixNano() - rqStart) / int64(time.Millisecond)
				statRqChan <- RqStatObj{
					Method:  method,
					Error:   err,
					Timeout: rqTimeout,
				}
			}()
		}
	}

	q := url.Values{}
	for k, v := range params {
		q.Add(k, v)
	}
	if params["v"] == "" {
		q.Add("v", APIVersion)
	}
	q.Add("access_token", vk.AccessToken)

	// Формируем запрос
	req, err := http.NewRequest("POST", APIMethodURL+method, strings.NewReader(q.Encode()))
	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Добавляем контекст
	ctx, cancel := context.WithCancel(context.Background())
	key := vk.AccessToken + "_" + strconv.FormatInt(time.Now().UnixNano(), 32)
	contMap.Lock()
	contMap.h[key] = cancel
	contMap.Unlock()
	defer func() {
		contMap.Lock()
		delete(contMap.h, key)
		contMap.Unlock()
	}()

	if exited {
		err = errors.New("context canceled")
		return
	}

	// Отправляем запрос
	client := &http.Client{Transport: httpTr}
	resp, err := client.Do(req.WithContext(ctx))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if !strings.Contains(err.Error(), "connection reset by peer") && !strings.Contains(err.Error(), "context canceled") {
			log.Println("[error]", err)
		}
		return
	}

	// Если проблема с ответом
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		if !httpErrorReg.MatchString(err.Error()) {
			log.Println("[error]", resp.Status, resp.StatusCode)
		}
		return
	}

	// Читаем ответ
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if !httpErrorReg.MatchString(err.Error()) {
			log.Println("[error]", err)
		}
		return
	}

	// Парсим ответ
	err = json.Unmarshal(content, &ans)
	if err != nil {
		log.Println("[error]", method, err, string(content))
		return
	}

	return
}

// Ждем между запросами если вк ответил что запросы слишком частые
func (vk *API) floodWait(method string) (ok bool) {

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
func (vk *API) httpErrorWait(method string) (ok bool) {
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
