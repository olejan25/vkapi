package vkapi

import "encoding/json"

// CallBackObj - Объект запроса через callback
type CallBackObj struct {
	Type    string          `json:"type"`
	Object  json.RawMessage `json:"object"`
	GroupID int             `json:"group_id"`
	Secret  string          `json:"secret"`

	Message             MessagesGetAns        `json:"-"`
	MessageAllow        CallbackMessageAllow  `json:"-"`
	Photo               PhotosGetItem         `json:"-"`
	PhotoComment        WallGetCommentsItem   `json:"-"`
	PhotoCommentDelete  CallbackCommentDelete `json:"-"`
	Video               VideoGetItem          `json:"-"`
	VideoComment        WallGetCommentsItem   `json:"-"`
	VideoCommentDelete  CallbackCommentDelete `json:"-"`
	Wall                WallGetByIDAns        `json:"-"`
	WallComment         WallGetCommentsItem   `json:"-"`
	WallCommentDelete   CallbackCommentDelete `json:"-"`
	Board               BoardGetTopicsItem    `json:"-"`
	BoardDelete         CallbackCommentDelete `json:"-"`
	MarketComment       WallGetCommentsItem   `json:"-"`
	MarketCommentDelete CallbackCommentDelete `json:"-"`
	UserChange          CallBackUserChange    `json:"-"`
}

// CallbackMessageAllow - объект подписки на сообщения
type CallbackMessageAllow struct {
	UserID int    `json:"user_id"`
	Key    string `json:"key"`
}

// CallbackCommentDelete - объект инфы о удаленной фотке
type CallbackCommentDelete struct {
	OwnerID      int `json:"owner_id"`
	ID           int `json:"id"`
	UserID       int `json:"user_id"`
	DeleterID    int `json:"deleter_id"`
	PhotoID      int `json:"photo_id"`
	VideoID      int `json:"video_id"`
	PostID       int `json:"post_id"`
	TopicID      int `json:"topic_id"`
	TopicOwnerID int `json:"topic_owner_id"`
	ItemID       int `json:"item_id"`
}

// CallBackUserChange - объект выходя или входа человека
type CallBackUserChange struct {
	UserID   int    `json:"user_id"`
	Self     int    `json:"self"`
	JoinType string `json:"join_type"`
}
