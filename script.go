package vkapi

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

//ScriptWallGetByID - Получаем список постов по их ID (execute)
func (vk *API) ScriptWallGetByID(posts []string) (ans []WallGetByIDAns, err error) {
	// Разбиваем посты на нужное кол-во
	arr := chunkSliceString(posts, 100)
	// Формируем массив для запроса
	tmpArr := make([]string, len(arr))
	for i, v := range arr {
		tmpArr[i] = strings.Join(v, ",")
	}

	b, err := json.Marshal(tmpArr)
	if err != nil {
		log.Println(err)
		return
	}

	script := fmt.Sprintf(`
		var arr = %s;
		var ans = [];

		while(arr.length > 0) {
			var str = arr.shift();

			var posts = API.wall.getById({
				posts: str,
				copy_history_depth: 1,
			});

			if(posts) {
				ans = ans + posts;
			}
		}

		return ans;
	`, b)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptGroupsGetByID - Получаем группы по их ID (execute)
func (vk *API) ScriptGroupsGetByID(groupIds []string, fields string) (ans []GroupsGetAns, err error) {
	// Разбиваем посты на нужное кол-во
	arr := chunkSliceString(groupIds, 500)
	// Формируем массив для запроса
	tmpArr := make([]string, len(arr))
	for i, v := range arr {
		tmpArr[i] = strings.Join(v, ",")
	}

	b, err := json.Marshal(tmpArr)
	if err != nil {
		log.Println(err)
		return
	}

	script := fmt.Sprintf(`
		var fields = "%s";
		var arr = %s;
		var ans = [];

		while(arr.length > 0) {
			var str = arr.shift();

			var groups = API.groups.getById({
				group_ids: str,
				fields: fields,
			});

			if(groups) {
				ans = ans + groups;
			}
		}

		return ans;
	`, fields, b)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptStatsGet - Получаем статистику групп. Максимум 25. (execute)
func (vk *API) ScriptStatsGet(groupIds []string, dateFrom, dateTo time.Time) (ans []StatsGetAns, err error) {
	b, err := json.Marshal(groupIds)
	if err != nil {
		log.Println(err)
		return
	}

	script := fmt.Sprintf(`
		var arr = %s;
		var ans = [];

		while(arr.length > 0) {
			var grid = arr.shift();

			var stats = API.stats.get({
				group_id  : grid,
				date_from : "%s",
				date_to   : "%s"
			});

			if(stats) {
				var st = stats[0];
				st.group_id = parseInt(grid);
				ans.push(st);
			}
		}

		return ans;
	`, b, dateFrom.Format("2006-02-01"), dateTo.Format("2006-02-01"))

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptUtilsResolveScreenName - Резольвим короткие имена в айдишники. максимум 25. (execute)
func (vk *API) ScriptUtilsResolveScreenName(ids []string) (ans []UtilsResolveScreenNameAns, err error) {
	b, err := json.Marshal(ids)
	if err != nil {
		log.Println(err)
		return
	}

	script := fmt.Sprintf(`
		var arr = %s;
		var ans = [];

		while(arr.length > 0) {
			var sn = arr.shift();

			var res = API.utils.resolveScreenName({
				screen_name: sn
			});

			if(res) {
				ans.push(res);
			}
		}

		return ans;
	`, b)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptGroupsGetMembers - Получаем подписчиков группы. (execute)
func (vk *API) ScriptGroupsGetMembers(groupID, offset int, s string) (ans ScriptGroupsGetMembersAns, err error) {
	if s == "" {
		s = "id_asc"
	}

	script := fmt.Sprintf(`
		var group_id = %d;
		var offset   = %d;
		var sort     = "%s";

		var cnt   = 25;
		var count = offset + 1;
		var users = [];

		while(cnt > 0 && offset < count){
			var res = API.groups.getMembers({ 
				group_id : group_id, 
				offset   : offset, 
				sort     : sort, 
				count    : 1000
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count  = res.count; 
				users  = users + res.items;
				offset = offset + 1000;
			}
			else {
				cnt = 0;
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items	 : users
		};
		
		return result;
	`, groupID, offset, s)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptUsersGetFollowers - Получаем подписчиков человека. (execute)
func (vk *API) ScriptUsersGetFollowers(userID, offset int) (ans ScriptGroupsGetMembersAns, err error) {

	script := fmt.Sprintf(`
		var user_id = %d;
		var offset  = %d;

		var cnt   = 25;
		var count = offset + 1;
		var users = [];

		while(cnt > 0 && offset < count){
			var res = API.users.getFollowers({ 
				user_id : user_id, 
				offset  : offset, 
				count   : 1000
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count  = res.count; 
				users  = users + res.items;
				offset = offset + 1000;
			}
			else {
				cnt = 0;
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items	 : users
		};
		
		return result;
	`, userID, offset)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptFriendsGet - Получаем друзей человека. (execute)
func (vk *API) ScriptFriendsGet(userID, offset int) (ans ScriptGroupsGetMembersAns, err error) {

	script := fmt.Sprintf(`
		var user_id = %d;
		var offset  = %d;

		var cnt   = 25;
		var count = offset + 1;
		var users = [];

		while(cnt > 0 && offset < count){
			var res = API.friends.get({ 
				user_id : user_id, 
				offset  : offset, 
				count   : 5000
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count  = res.count; 
				users  = users + res.items;
				offset = offset + 5000;
			}
			else {
				cnt = 0;
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items	 : users
		};
		
		return result;
	`, userID, offset)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptWallGetComments - Получаем комментарии поста. (execute)
func (vk *API) ScriptWallGetComments(ownerID, postID, startCommentID int) (ans WallGetCommentsAns, err error) {

	script := fmt.Sprintf(`
		var owner_id         = %d;
		var post_id          = %d;
		var start_comment_id = %d;

		var cnt         = 25;
		var real_offset = 0;
		var offset      = 0;
		var count       = offset + 1;
		var comments    = [];
		var limit       = 100;

		while(cnt > 0 && offset < count){
			var res = API.wall.getComments({ 
				owner_id         : owner_id, 
				post_id          : post_id,
				start_comment_id : start_comment_id,
				offset           : offset,
				sort             : "desc",
				need_likes       : 1,
				count            : limit
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count       = res.count; 
				comments    = comments + res.items;
				offset      = offset + limit;
				real_offset = res.real_offset + limit;
			}
			else {
				cnt = 0;
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items  : comments
		};
		
		return result;
	`, ownerID, postID, startCommentID)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptLikesGetList - Получаем лайки. (execute)
func (vk *API) ScriptLikesGetList(ownerID, itemID int, t, filter, pageURL string, offset int) (ans LikesGetListAns, err error) {

	script := fmt.Sprintf(`
		var owner_id = %d;
		var item_id  = %d;
		var type     = "%s";
		var filter   = "%s";
		var page_url = "%s";
		var offset   = %d;
	
		var cnt   = 25;
		var count = offset + 1;
		var users = [];

		while(cnt > 0 && offset < count){
			var res = API.likes.getList({ 
				owner_id : owner_id, 
				item_id  : item_id,
				type     : type,
				filter   : filter,
				page_url : page_url,
				offset   : offset,
				count    : 1000
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count  = res.count; 
				users  = users + res.items;
				offset = offset + 1000;
			}
			else {
				cnt = 0;
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items  : users
		};
		
		return result;
	`, ownerID, itemID, t, filter, pageURL, offset)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptBoardGetTopics - Получаем обсуждения. (execute)
func (vk *API) ScriptBoardGetTopics(groupID, offset int) (ans BoardGetTopicsAns, err error) {

	script := fmt.Sprintf(`
		var group_id = %d;
		var offset   = %d;
	
		var cnt   = 25;
		var count = offset + 1;
		var topics = [];

		while(cnt > 0 && offset < count){
			var res = API.board.getTopics({ 
				group_id : group_id, 
				order    : -2,
				offset   : offset,
				count    : 100
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count  = res.count; 
				topics = topics + res.items;
				offset = offset + 100;
			}
			else {
				cnt = 0;
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items  : topics
		};
		
		return result;
	`, groupID, offset)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptBoardGetComments - Получаем комментарии обсуждений. (execute)
func (vk *API) ScriptBoardGetComments(groupID, topicID, startCommentID, cnt int) (ans BoardGetCommentsAns, err error) {

	script := fmt.Sprintf(`
		var group_id         = %d;
		var topic_id         = %d;
		var start_comment_id = %d;
	
		var cnt         = %d;
		var real_offset = 0;
		var offset      = 0;
		var count       = offset + 1;
		var comments    = [];
		var limit       = 100;

		while(cnt > 0 && real_offset < count){
			var res = API.board.getComments({ 
				group_id         : group_id, 
				topic_id         : topic_id,
				need_likes       : 1,
				start_comment_id : start_comment_id,
				offset           : offset,
				sort             : "desc",
				count            : limit
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count       = res.count; 
				comments    = comments + res.items;
				offset      = offset + limit;
				real_offset = res.real_offset + limit;
			}
			else {
				cnt = 0;
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items  : comments
		};
		
		return result;
	`, groupID, topicID, startCommentID, cnt)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptVideoGet - Получаем видео сообщества или пользователя. (execute)
func (vk *API) ScriptVideoGet(ownerID, offset int) (ans VideoGetAns, err error) {

	script := fmt.Sprintf(`
		var owner_id = %d;
		var offset   = %d;
	
		var cnt    = 25;
		var count  = offset + 1;
		var videos = [];
		var limit  = 200;

		while(cnt > 0 && offset < count){
			var res = API.video.get({ 
				owner_id   : owner_id,
				offset     : offset,
				count      : limit,
				extended   : 1,
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count  = res.count; 
				videos = videos + res.items;
				offset = offset + limit;
			}
			else {
				cnt = 0;
				if count == 1 {
					offset = count;
				}
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items  : videos
		};
		
		return result;
	`, ownerID, offset)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptVideoGetComments - Получаем комментарии к видео. (execute)
func (vk *API) ScriptVideoGetComments(ownerID, videoID, startCommentID int) (ans VideoGetCommentsAns, err error) {

	script := fmt.Sprintf(`
		var owner_id         = %d;
		var video_id         = %d;
		var start_comment_id = %d;
	
		var cnt         = 25;
		var offset      = 0;
		var real_offset = 0;
		var count       = offset + 1;
		var comments    = [];
		var limit       = 100;

		while(cnt > 0 && real_offset < count){
			var res = API.video.getComments({ 
				owner_id         : owner_id,
				video_id         : video_id,
				need_likes       : 1,
				start_comment_id : start_comment_id,
				sort             : "desc",
				offset           : offset,
				count            : limit
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count       = res.count; 
				comments    = comments + res.items;
				offset      = offset + limit;
		 		real_offset = res.real_offset + limit;
			}
			else {
				cnt = 0;
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items  : comments
		};
		
		return result;
	`, ownerID, videoID, startCommentID)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptPhotosGet - Получаем фото из альбома. (execute)
func (vk *API) ScriptPhotosGet(ownerID, albumID, offset int) (ans PhotosGetAns, err error) {

	script := fmt.Sprintf(`
		var owner_id = %d;
		var album_id = %d;
		var offset   = %d;
	
		var cnt    = 25;
		var count  = offset + 1;
		var photos = [];
		var limit  = 1000;

		while(cnt > 0 && offset < count){
			var res = API.photos.get({ 
				owner_id   : owner_id,
				album_id   : album_id,
				rev        : 1,
				extended   : 1,
				offset     : offset,
				count      : limit
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count  = res.count; 
				photos = photos + res.items;
				offset = offset + limit;
			}
			else {
				cnt = 0;
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items  : photos
		};
		
		return result;
	`, ownerID, albumID, offset)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// ScriptPhotosGetComments - Получаем комментарии фото. (execute)
func (vk *API) ScriptPhotosGetComments(ownerID, photoID, StartCommentID int) (ans PhotosGetCommentsAns, err error) {

	script := fmt.Sprintf(`
		var owner_id         = %d;
		var photo_id         = %d;
		var start_comment_id = %d;
	
		var cnt         = 25;
		var offset      = 0;
		var real_offset = 0;
		var count       = offset + 1;
		var comments    = [];
		var limit       = 100;

		while(cnt > 0 && real_offset < count){
			var res = API.photos.getComments({ 
				owner_id         : owner_id,
				photo_id         : photo_id,
				start_comment_id : start_comment_id,
				sort             : "desc",
				need_likes       : 1,
				offset           : offset,
				count            : limit
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count       = res.count; 
				comments    = comments + res.items;
				offset      = offset + limit;
				real_offset = res.real_offset + limit;
			}
			else {
				cnt = 0;
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items  : comments
		};
		
		return result;
	`, ownerID, photoID, StartCommentID)

	r, err := vk.Execute(script)
	if err != nil {
		if !executeErrorSkipReg.MatchString(err.Error()) {
			if !vk.checkErrorSkip(err.Error()) {
				log.Println("[error]", err)
			}
		}
		return
	}

	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}
