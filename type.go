package vkapi

import (
	"encoding/json"
	"log"
	"regexp"
	"sync"
)

var (
	executeErrorSkipReg *regexp.Regexp
	httpErrorReg        *regexp.Regexp
)

func init() {
	executeErrorSkipReg = regexp.MustCompile("server sent GOAWAY|User authorization failed|unexpected EOF|Database problems, try later|Internal Server Error|Bad Request|Gateway Timeout|Bad Gateway|could not check access_token now|connection reset by peer")
	httpErrorReg = regexp.MustCompile("unexpected EOF|server sent GOAWAY|Bad Request|Internal Server Error")
}

// API - главный объект
type API struct {
	AccessToken    string
	retryCount     int
	httpRetryCount int
	ErrorToSkip    []string
	sync.Mutex
}

// TokenData - объект получения токена
type TokenData struct {
	ClientID     int
	ClientSecret string
	Code         string
	RedirectURI  string
}

// Response - объект ответа VK
type Response struct {
	Response      json.RawMessage `json:"response"`
	Error         ResponseError   `json:"error"`
	ExecuteErrors []ExecuteErrors `json:"execute_errors"`
}

// ExecuteErrors - объект ошибок execute
type ExecuteErrors struct {
	Method    string `json:"method"`
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// ResponseError - объект ошибки выболнения запроса
type ResponseError struct {
	ErrorCode     int                 `json:"error_code"`
	ErrorMsg      string              `json:"error_msg"`
	RequestParams []map[string]string `json:"request_params"`
}

/*
	Получение токена
*/

// GetTokenAns - объект ответа при получении покена
type GetTokenAns struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	UserID           int    `json:"user_id"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

/*
	Users
*/

// UsersGetAns - объект ответа при запросе пользователей
type UsersGetAns struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Photo100  string `json:"photo_100"`
	Site      string `json:"site"`
	Sex       int    `json:"sex"`
	Status    string `json:"status"`
	Role      string `json:"role"`
}

/*
	Groups
*/

// GroupsGetAns - объект ответа при запросе групп
type GroupsGetAns struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ScreenName   string `json:"screen_name"`
	IsClosed     int    `json:"is_closed"`
	Deactivated  string `json:"deactivated"`
	IsAdmin      int    `json:"is_admin"`
	AdminLevel   int    `json:"admin_level"`
	IsMember     int    `json:"is_member"`
	InvitedBy    int    `json:"invited_by"`
	Type         string `json:"type"`
	Photo50      string `json:"photo_50"`
	Photo100     string `json:"photo_100"`
	Photo200     string `json:"photo_200"`
	AgeLimits    int    `json:"age_limits "`
	Description  string `json:"description"`
	MembersCount int    `json:"members_count"`
	Verified     int    `json:"verified"`
}

// GroupsGetMembersAns - объект ответа при запросе подписчиков групп
type GroupsGetMembersAns struct {
	Count int             `json:"count"`
	Items json.RawMessage `json:"items"`
}

// ScriptGroupsGetMembersAns - объект ответа при подписчиков (execute)
type ScriptGroupsGetMembersAns struct {
	Count  int   `json:"count"`
	Offset int   `json:"offset"`
	Users  []int `json:"users"`
}

/*
	Stats
*/

// StatsGetAns - объект ответа при запросе статистики группы
type StatsGetAns struct {
	GroupID          int             `json:"group_id"`
	Day              string          `json:"day"`
	Views            int             `json:"views"`
	Visitors         int             `json:"visitors"`
	Reach            int             `json:"reach"`
	ReachSubscribers int             `json:"reach_subscribers"`
	Subscribed       int             `json:"subscribed"`
	Unsubscribed     int             `json:"unsubscribed"`
	Sex              []StatsGetValue `json:"sex"`
	Age              []StatsGetValue `json:"age"`
	SexAge           []StatsGetValue `json:"sex_age"`
	Cities           []StatsGetValue `json:"cities"`
	Countries        []StatsGetValue `json:"countries"`
}

// StatsGetValue - объект статистики
type StatsGetValue struct {
	Visitors int         `json:"visitors"`
	Value    interface{} `json:"value"`
	Name     string      `json:"name"`
}

/*
	Wall
*/

// WallGetAns - объект ответа при запросе постов
type WallGetAns struct {
	Count int              `json:"count"`
	Items []WallGetByIDAns `json:"items"`
}

// WallGetByIDAns - обект постов
type WallGetByIDAns struct {
	ID           int              `json:"id"`
	OwnerID      int              `json:"owner_id"`
	FromID       int              `json:"from_id"`
	CreatedBy    int              `json:"created_by"`
	Date         int              `json:"date"`
	Text         string           `json:"text"`
	ReplyOwnerID int              `json:"reply_owner_id"`
	ReplyPostID  int              `json:"reply_post_id"`
	FriendsOnly  int              `json:"friends_only"`
	Comments     CommentData      `json:"comments"`
	Likes        LikeData         `json:"likes"`
	Reposts      LikeData         `json:"reposts"`
	Views        LikeData         `json:"views"`
	PostType     string           `json:"post_type"`
	Attachments  []Attachments    `json:"attachments"`
	SignerID     int              `json:"signer_id"`
	CopyHistory  []WallGetByIDAns `json:"copy_history"`
	IsPinned     int              `json:"is_pinned"`
	MarkedAsAds  int              `json:"marked_as_ads"`
}

// Attachments - объект аттача
type Attachments struct {
	Type        string           `json:"type"`
	Photo       *json.RawMessage `json:"photo"`
	Audio       *json.RawMessage `json:"audio"`
	Video       *json.RawMessage `json:"video"`
	Poll        *json.RawMessage `json:"poll"`
	Page        *json.RawMessage `json:"page"`
	Album       *json.RawMessage `json:"album"`
	Link        *json.RawMessage `json:"link"`
	Doc         *json.RawMessage `json:"doc"`
	Note        *json.RawMessage `json:"note"`
	Sticker     *json.RawMessage `json:"sticker"`
	PrettyCards *json.RawMessage `json:"pretty_cards"`
}

// AttachmentsPrettyCards - объект карточек аттача
type AttachmentsPrettyCards struct {
	Cards []AttachmentsPrettyCardsCards `json:"cards"`
}

// AttachmentsPrettyCardsCards - объект карточек аттача
type AttachmentsPrettyCardsCards struct {
	CardID  string `json:"card_id"`
	LinkURL string `json:"link_url"`
	Title   string `json:"title"`
}

// GetPrettyCards - Преобрахуем данные карточек в объекты
func (a *Attachments) GetPrettyCards() (t AttachmentsPrettyCards) {
	err := json.Unmarshal(*a.PrettyCards, &t)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// CommentData - объект комментариев
type CommentData struct {
	Count         int  `json:"count"`
	CanPost       int  `json:"can_post"`
	GroupsCanPost bool `json:"groups_can_post"`
}

// LikeData - объект лайков
type LikeData struct {
	Count int `json:"count"`
}

// WallGetCommentsAns - объект комментариев
type WallGetCommentsAns struct {
	Count  int                   `json:"count"`
	Offset int                   `json:"offset"`
	Items  []WallGetCommentsItem `json:"items"`
}

// WallGetCommentsItem - объект комментария
type WallGetCommentsItem struct {
	ID             int           `json:"id"`
	FromID         int           `json:"from_id"`
	Date           int           `json:"date"`
	Text           string        `json:"text"`
	ReplyToUser    int           `json:"reply_to_user"`
	ReplyToComment int           `json:"reply_to_comment"`
	Attachments    []Attachments `json:"attachments"`
	Likes          LikeData      `json:"likes"`
}

/*
	Likes
*/

// LikesGetListAns - объект лайков
type LikesGetListAns struct {
	Count  int   `json:"count"`
	Offset int   `json:"offset"`
	Items  []int `json:"items"`
}

/*
	Utils
*/

// UtilsGetShortLinkAns - объект ответа при запросе короткой ссылки
type UtilsGetShortLinkAns struct {
	ShortURL  string `json:"short_url"`
	URL       string `json:"url"`
	Key       string `json:"key"`
	AccessKey string `json:"access_key"`
}

// UtilsGetLinkStatsAns - объект ответа при запросе статистики короткой ссылки
type UtilsGetLinkStatsAns struct {
	Key   string                   `json:"key"`
	Stats []UtilsGetLinkStatsStats `json:"stats"`
}

// UtilsGetLinkStatsStats - объект статистики короткой ссылки
type UtilsGetLinkStatsStats struct {
	Timestamp int         `json:"timestamp"`
	Views     int         `json:"views"`
	SexAge    []SexAge    `json:"sex_age"`
	Countries []Countries `json:"countries"`
	Cities    []Cities    `json:"cities"`
}

// UtilsResolveScreenNameAns - объект ответа при запросе резольвинка короткого имени
type UtilsResolveScreenNameAns struct {
	Type     string `json:"type"`
	ObjectID int    `json:"object_id"`
}

/*
	Board
*/

// BoardGetTopicsAns - объект списка обсуждений
type BoardGetTopicsAns struct {
	Count  int                  `json:"count"`
	Offset int                  `json:"offset"`
	Items  []BoardGetTopicsItem `json:"items"`
}

// BoardGetTopicsItem - объект обсуждения
type BoardGetTopicsItem struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Created      int    `json:"created"`
	CreatedBy    int    `json:"created_by"`
	Updated      int    `json:"updated"`
	UpdatedBy    int    `json:"updated_by"`
	IsClosed     int    `json:"is_closed"`
	IsFixed      int    `json:"is_fixed"`
	Comments     int    `json:"comments"`
	FirstComment int    `json:"first_comment"`
	LastComment  int    `json:"last_comment"`
}

// BoardGetCommentsAns - объект списка комментариев обсуждения
type BoardGetCommentsAns struct {
	Count  int                   `json:"count"`
	Offset int                   `json:"offset"`
	Items  []WallGetCommentsItem `json:"items"`
}

/*
	Video
*/

// VideoGetAns - объект списка видео
type VideoGetAns struct {
	Count  int            `json:"count"`
	Offset int            `json:"offset"`
	Items  []VideoGetItem `json:"items"`
}

// VideoGetItem - объект видео
type VideoGetItem struct {
	ID         int      `json:"id"`
	OwnerID    int      `json:"owner_id"`
	Title      string   `json:"title"`
	Duration   int      `json:"duration"`
	Date       int      `json:"date"`
	Comments   int      `json:"comments"`
	Views      int      `json:"views"`
	Likes      LikeData `json:"likes"`
	Reposts    LikeData `json:"reposts"`
	Platform   string   `json:"platform"`
	Player     string   `json:"player"`
	AddingDate int      `json:"adding_date"`
}

// VideoGetCommentsAns - объект списка комментариев
type VideoGetCommentsAns struct {
	Count  int                   `json:"count"`
	Offset int                   `json:"offset"`
	Items  []WallGetCommentsItem `json:"items"`
}

/*
	Ads
*/

// AdsGetCampaignsAns - объект ответа при запросе кампаний
type AdsGetCampaignsAns struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// AdsGetAdsLayoutAns - объект ответа при запросе вида объявления
type AdsGetAdsLayoutAns struct {
	ID         string `json:"id"`
	CampaignID int    `json:"campaign_id"`
	Title      string `json:"title"`
	LinkURL    string `json:"link_url"`
}

// AdsGetStatisticsAns - объект ответа при запросе статистики
type AdsGetStatisticsAns struct {
	ID       int                        `json:"id"`
	Type     string                     `json:"type"`
	StatsBug []json.RawMessage          `json:"stats"`
	Stats    []AdsGetStatisticsAnsStats `json:"-"`
}

// AdsGetStatisticsAnsStats - объект статистики
type AdsGetStatisticsAnsStats struct {
	Day         string `json:"day"`
	Spent       string `json:"spent"`
	Impressions int    `json:"impressions"`
	Clicks      int    `json:"clicks"`
	Reach       int    `json:"reach"`
}

// AdsGetStatisticsAnsStatsBug - объект ответа при запросе статистики (если VK криво типы переменных сформировал)
type AdsGetStatisticsAnsStatsBug struct {
	Day         string `json:"day"`
	Spent       string `json:"spent"`
	Impressions string `json:"impressions"`
	Clicks      int    `json:"clicks"`
	Reach       int    `json:"reach"`
}

/*
	Other
*/

// SexAge - объект пола/возраста
type SexAge struct {
	AgeRange string `json:"age_range"`
	Female   int    `json:"female"`
	Male     int    `json:"male"`
}

// Countries - объект статистики по странам
type Countries struct {
	CountryID int `json:"country_id"`
	Views     int `json:"views"`
}

// Cities - объект статистики по городам
type Cities struct {
	CityID int `json:"city_id"`
	Views  int `json:"views"`
}
