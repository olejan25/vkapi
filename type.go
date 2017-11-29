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
	executeErrorSkipReg = regexp.MustCompile("server sent GOAWAY|User authorization failed|unexpected EOF|Database problems, try later|Internal Server Error|Bad Request")
	httpErrorReg = regexp.MustCompile("unexpected EOF|server sent GOAWAY|Bad Request|Internal Server Error")
}

type Api struct {
	AccessToken    string
	retryCount     int
	httpRetryCount int
	sync.Mutex
}

type TokenData struct {
	ClientId     int
	ClientSecret string
	Code         string
	RedirectUri  string
}

type Response struct {
	Response      json.RawMessage `json:"response"`
	Error         ResponseError   `json:"error"`
	ExecuteErrors []ExecuteErrors `json:"execute_errors"`
}

type ExecuteErrors struct {
	Method     string `json:"method"`
	Error_code int    `json:"error_code"`
	Error_msg  string `json:"error_msg"`
}

type ResponseError struct {
	ErrorCode     int                 `json:"error_code"`
	ErrorMsg      string              `json:"error_msg"`
	RequestParams []map[string]string `json:"request_params"`
}

/*
	Получение токена
*/

type GetTokenAns struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	UserId           int    `json:"user_id"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

/*
	Users
*/

type UsersGetAns struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Photo100  string `json:"photo_100"`
	Site      string `json:"site"`
	Sex       int    `json:"sex"`
	Status    string `json:"status"`
}

/*
	Groups
*/

type GroupsGetAns struct {
	Id           int    `json:"id"`
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

/*
	Wall
*/

type WallGetAns struct {
	Count int              `json:"count"`
	Items []WallGetByIdAns `json:"items"`
}

type WallGetByIdAns struct {
	Id           int              `json:"id"`
	OwnerId      int              `json:"owner_id"`
	FromId       int              `json:"from_id"`
	CreatedBy    int              `json:"created_by"`
	Date         int              `json:"date"`
	Text         string           `json:"text"`
	ReplyOwnerId int              `json:"reply_owner_id"`
	ReplyPostId  int              `json:"reply_post_id"`
	FriendsOnly  int              `json:"friends_only"`
	Comments     LikeData         `json:"comments"`
	Likes        LikeData         `json:"likes"`
	Reposts      LikeData         `json:"reposts"`
	Views        LikeData         `json:"views"`
	PostType     string           `json:"post_type"`
	Attachments  []Attachments    `json:"attachments"`
	SignerId     int              `json:"signer_id"`
	CopyHistory  []WallGetByIdAns `json:"copy_history"`
	IsPinned     int              `json:"is_pinned"`
	MarkedAsAds  int              `json:"marked_as_ads"`
}

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

type AttachmentsPrettyCards struct {
	Cards []AttachmentsPrettyCardsCards `json:"cards"`
}

type AttachmentsPrettyCardsCards struct {
	CardId  string `json:"card_id"`
	LinkUrl string `json:"link_url"`
	Title   string `json:"title"`
}

func (a *Attachments) GetPrettyCards() (t AttachmentsPrettyCards) {
	err := json.Unmarshal(*a.PrettyCards, &t)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

type LikeData struct {
	Count int `json:"count"`
}

/*
	Utils
*/
type UtilsGetShortLinkAns struct {
	ShortUrl  string `json:"short_url"`
	Url       string `json:"url"`
	Key       string `json:"key"`
	AccessKey string `json:"access_key"`
}

type UtilsGetLinkStatsAns struct {
	Key   string                   `json:"key"`
	Stats []UtilsGetLinkStatsStats `json:"stats"`
}

type UtilsGetLinkStatsStats struct {
	Timestamp int         `json:"timestamp"`
	Views     int         `json:"views"`
	SexAge    []SexAge    `json:"sex_age"`
	Countries []Countries `json:"countries"`
	Cities    []Cities    `json:"cities"`
}

/*
	Ads
*/

type AdsGetCampaignsAns struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type AdsGetAdsLayoutAns struct {
	Id         string `json:"id"`
	CampaignId int    `json:"campaign_id"`
	Title      string `json:"title"`
	LinkUrl    string `json:"link_url"`
}

type AdsGetStatisticsAns struct {
	Id       int                        `json:"id"`
	Type     string                     `json:"type"`
	StatsBug []json.RawMessage          `json:"stats"`
	Stats    []AdsGetStatisticsAnsStats `json:"-"`
}

type AdsGetStatisticsAnsStats struct {
	Day         string `json:"day"`
	Spent       string `json:"spent"`
	Impressions int    `json:"impressions"`
	Clicks      int    `json:"clicks"`
	Reach       int    `json:"reach"`
}

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

type SexAge struct {
	AgeRange string `json:"age_range"`
	Female   int    `json:"female"`
	Male     int    `json:"male"`
}

type Countries struct {
	CountryId int `json:"country_id"`
	Views     int `json:"views"`
}

type Cities struct {
	CityId int `json:"city_id"`
	Views  int `json:"views"`
}
