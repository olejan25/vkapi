package vkapi

import (
	"encoding/json"
	"log"
)

// CallBackObj - Объект запроса через callback
type CallBackObj struct {
	Type    string          `json:"type"`
	Object  json.RawMessage `json:"object"`
	GroupID int             `json:"group_id"`
	Secret  string          `json:"secret"`

	Message       MessagesGetAns        `json:"-"`
	MessageAllow  CallbackMessageAllow  `json:"-"`
	Photo         PhotosGetItem         `json:"-"`
	PhotoComment  WallGetCommentsItem   `json:"-"`
	Video         VideoGetItem          `json:"-"`
	VideoComment  WallGetCommentsItem   `json:"-"`
	Wall          WallGetByIDAns        `json:"-"`
	WallComment   WallGetCommentsItem   `json:"-"`
	Board         BoardGetTopicsItem    `json:"-"`
	MarketComment WallGetCommentsItem   `json:"-"`
	UserChange    CallBackUserChange    `json:"-"`
	CommentDelete CallbackCommentDelete `json:"-"`
}

// Parse - Парсим объект
func (cbo *CallBackObj) Parse() (err error) {

	switch cbo.Type {
	case "message_new", "message_reply", "message_edit":
		err = json.Unmarshal(cbo.Object, &cbo.Message)
	case "message_allow", "message_deny":
		err = json.Unmarshal(cbo.Object, &cbo.MessageAllow)
	case "photo_new":
		err = json.Unmarshal(cbo.Object, &cbo.Photo)
	case "photo_comment_new", "photo_comment_edit", "photo_comment_restore":
		err = json.Unmarshal(cbo.Object, &cbo.PhotoComment)
	case "video_new":
		err = json.Unmarshal(cbo.Object, &cbo.Video)
	case "video_comment_new", "video_comment_edit", "video_comment_restore":
		err = json.Unmarshal(cbo.Object, &cbo.VideoComment)
	case "wall_post_new", "wall_repost":
		err = json.Unmarshal(cbo.Object, &cbo.Wall)
	case "wall_reply_new", "wall_reply_edit", "wall_reply_restore":
		err = json.Unmarshal(cbo.Object, &cbo.WallComment)
	case "board_post_new", "board_post_edit", "board_post_restore":
		err = json.Unmarshal(cbo.Object, &cbo.Board)
	case "market_comment_new", "market_comment_edit", "market_comment_restore":
		err = json.Unmarshal(cbo.Object, &cbo.MarketComment)
	case "group_leave", "group_join":
		err = json.Unmarshal(cbo.Object, &cbo.UserChange)
	case "photo_comment_delete", "video_comment_delete", "wall_reply_delete", "board_post_delete", "market_comment_delete":
		err = json.Unmarshal(cbo.Object, &cbo.CommentDelete)
	}

	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
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
