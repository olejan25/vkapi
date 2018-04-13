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
func (vk *API) ScriptGroupsGetByID(groupIDs []string, fields string) (ans []GroupsGetAns, err error) {
	// Разбиваем посты на нужное кол-во
	arr := chunkSliceString(groupIDs, 500)
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
		var limit = 1000;

		while(cnt > 0 && offset < count){
			var res = API.users.getFollowers({ 
				user_id : user_id, 
				offset  : offset, 
				count   : limit
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count  = res.count; 
				users  = users + res.items;
				offset = offset + limit;
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

// ScriptMultiUsersGetFollowers - Получаем подписчиков человека. (execute)
func (vk *API) ScriptMultiUsersGetFollowers(arr []map[string]interface{}) (ans ScriptMultiFriendsGetAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr     = %s;
		var users   = [];
		var rq_data = [];
		var limit   = 1000;

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.users.getFollowers({ 
				user_id : h.user_id, 
				count   : limit
			}); 

			if(res.count) {
				res.offset = limit;
				users.push(res);
				rq_data.push(h);
			}
		}

		var result = {
			items   : users,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptMultiFriendsGet - Получаем друзей человеков. (execute)
func (vk *API) ScriptMultiFriendsGet(arr []map[string]interface{}) (ans ScriptMultiFriendsGetAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr     = %s;
		var users   = [];
		var rq_data = [];
		var limit   = 5000;

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.friends.get({
				user_id : h.user_id,
				count   : limit
			}); 

			if(res.count) {
				res.offset = limit;
				users.push(res);
				rq_data.push(h);
			}
		}

		var result = {
			items   : users,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptMultiWallGet - Получаем посты разных сообществ и людей. (execute)
func (vk *API) ScriptMultiWallGet(arr []map[string]interface{}) (ans MultiWallGetAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr     = %s;
		var posts   = [];
		var rq_data = [];
		var limit   = 100;

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.wall.get({ 
				owner_id : h.owner_id,
				sort     : "desc",
				count    : limit
			}); 

			if(res.count) {
				res.offset = limit;
				posts.push(res);
				rq_data.push(h);
			}
		}

		var result = {
			items   : posts,
			rq_data : rq_data,
		};
		
		return result;
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

		while(cnt > 0 && real_offset < count){
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

// ScriptMultiWallGetComments - Получаем комментарии нескольких постов. (execute)
func (vk *API) ScriptMultiWallGetComments(arr []map[string]interface{}) (ans MultiWallGetCommentsAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr      = %s;
		var comments = [];
		var rq_data  = [];
		var limit    = 100;

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.wall.getComments({ 
				owner_id   : h.owner_id, 
				post_id    : h.post_id,
				sort       : "desc",
				need_likes : 1,
				count      : limit
			}); 

			if(res.count) {
				res.offset = limit;
				comments.push(res);
				rq_data.push(h);
			}
		}

		var result = {
			items   : comments,
			rq_data : rq_data,
		};
		
		return result;
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
		var limit = 1000;

		while(cnt > 0 && offset < count){
			var res = API.likes.getList({ 
				owner_id : owner_id, 
				item_id  : item_id,
				type     : type,
				filter   : filter,
				page_url : page_url,
				offset   : offset,
				count    : limit
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count  = res.count; 
				users  = users + res.items;
				offset = offset + limit;
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

// ScriptMultiLikesGetList - Получаем лайки у нескольких объектов. (execute)
func (vk *API) ScriptMultiLikesGetList(arr []map[string]interface{}) (ans MultiLikesGetListAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr     = %s;
		var users   = [];
		var rq_data = [];
		var limit   = 1000;

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.likes.getList({ 
				owner_id : h.owner_id, 
				item_id  : h.item_id,
				type     : h.type,
				filter   : h.filter,
				page_url : h.page_url,
				count    : limit
			}); 

			if(res.count) {
				res.offset = limit;
				users.push(res);
				rq_data.push(h);
			}
		}

		var result = {
			items   : users,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptMultiBoardGetTopics - Получаем обсуждения. (execute)
func (vk *API) ScriptMultiBoardGetTopics(arr []map[string]interface{}) (ans MultiBoardGetTopicsAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr     = %s;
		var topics  = [];
		var rq_data = [];
		var limit   = 100;

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.board.getTopics({ 
				group_id : h.group_id, 
				order    : -2,
				count    : limit
			}); 

			if(res.count) {
				res.offset = limit;
				topics.push(res);
				rq_data.push(h);
			}
		}

		var result = {
			items   : topics,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptMultiBoardGetComments - Получаем комментарии нескольких обсуждений. (execute)
func (vk *API) ScriptMultiBoardGetComments(arr []map[string]interface{}) (ans MultiBoardGetCommentsAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr      = %s;
		var comments = [];
		var limit    = 100;
		var rq_data  = [];

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.board.getComments({ 
				group_id   : h.group_id, 
				topic_id   : h.topic_id,
				need_likes : 1,
				sort       : "desc",
				count      : limit
			}); 

			if(res.count) {
				res.offset = limit;
				comments.push(res);
				rq_data.push(h);
			}
		}

		var result = {
			items   : comments,
			rq_data : rq_data,
		};
		
		return result;
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
				if(count == offset + 1) {
					count = offset;
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

// ScriptMultiVideoGet - Получаем видео сообщества или пользователя. (execute)
func (vk *API) ScriptMultiVideoGet(arr []map[string]interface{}) (ans MultiVideoGetAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr     = %s;
		var rq_data = [];
		var videos  = [];
		var limit   = 200;
	
		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.video.get({ 
				owner_id   : h.owner_id,
				count      : limit,
				extended   : 1,
			}); 

			if(res.count) {
				res.offset = limit;
				videos.push(res);
		 		rq_data.push(h);
			}
		}

		var result = {
			items   : videos,
			rq_data : rq_data,
		};
		
		return result;
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

//ScriptVideoGetByID - Получаем список видео по их ID (execute)
func (vk *API) ScriptVideoGetByID(videos []string) (ans VideoGetAns, err error) {
	// Разбиваем посты на нужное кол-во
	arr := chunkSliceString(videos, 100)
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
		var arr    = %s;
		var videos = [];
		var count  = 0;

		while(arr.length > 0) {
			var str = arr.shift();

			var res = API.video.get({
				videos:   str,
				extended: 1,
			});

			if(res.count) {
				count = res.count;
				videos = videos + res.items;
			}
		}

		var ans = {
			count: count,
			items: videos,
		};

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

// ScriptMultiVideoGetComments - Получаем комментарии к нескольким видео. (execute)
func (vk *API) ScriptMultiVideoGetComments(arr []map[string]interface{}) (ans MultiVideoGetCommentsAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr      = %s;
		var rq_data  = [];
		var comments = [];
		var limit    = 100;

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.video.getComments({ 
				owner_id   : h.owner_id,
				video_id   : h.video_id,
				need_likes : 1,
				sort       : "desc",
				count      : limit
			});

			if(res.count) {
				res.offset = limit;
				comments.push(res);
		 		rq_data.push(h);
			}
		}

		var result = {
			items   : comments,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptMultiPhotosGetAlbums - Получаем фото альбомы. (execute)
func (vk *API) ScriptMultiPhotosGetAlbums(arr []map[string]interface{}) (ans MultiPhotosGetAlbumsAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr     = %s;
		var rq_data = [];
		var albums  = [];
	
		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.photos.getAlbums({ 
				owner_id : h.owner_id,
			}); 

			if(res.count) {
				albums.push(res);
		 		rq_data.push(h);
			}
		}

		var result = {
			items   : albums,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptMultiPhotosGet - Получаем фото из альбома. (execute)
func (vk *API) ScriptMultiPhotosGet(arr []map[string]interface{}) (ans MultiPhotosGetAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr     = %s;
		var rq_data = [];
		var photos  = [];
		var limit   = 1000;


		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.photos.get({ 
				owner_id : h.owner_id,
				album_id : h.album_id,
				rev      : 1,
				extended : 1,
				count    : limit,
			});

			if(res.count) {
				res.offset = limit;
				photos.push(res);
		 		rq_data.push(h);
			}
		}

		var result = {
			items   : photos,
			rq_data : rq_data,
		};
		
		return result;
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

//ScriptPhotosGetByID - Получаем список фото по их ID (execute)
func (vk *API) ScriptPhotosGetByID(photos []string) (ans PhotosGetAns, err error) {
	// Разбиваем посты на нужное кол-во
	arr := chunkSliceString(photos, 100)
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
		var arr    = %s;
		var photos = [];

		while(arr.length > 0) {
			var str = arr.shift();

			var res = API.photos.getById({
				photos:   str,
				extended: 1,
			});

			if(res) {
				photos = photos + res;
			}
		}

		var ans = {
			items: photos,
		};

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

// ScriptMultiPhotosGetComments - Получаем комментарии нескольких фото. (execute)
func (vk *API) ScriptMultiPhotosGetComments(arr []map[string]interface{}) (ans MultiPhotosGetCommentsAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr      = %s;
		var rq_data  = [];
		var comments = [];
		var limit    = 100;

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.photos.getComments({ 
				owner_id         : h.owner_id,
				photo_id         : h.photo_id,
				sort             : "desc",
				need_likes       : 1,
				count            : limit
			});

			if(res.count) {
				res.offset = limit;
				comments.push(res);
		 		rq_data.push(h);
			}
		}

		var result = {
			items   : comments,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptMultiUsersGetSubscriptions - Получаем подписки нескольких людей. (execute)
func (vk *API) ScriptMultiUsersGetSubscriptions(arr []map[string]interface{}) (ans MultiUsersGetSubscriptionsAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr           = %s;
		var rq_data       = [];
		var subscriptions = [];

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.users.getSubscriptions({ 
				user_id : h.user_id,
			});

			if(res) {
				subscriptions.push(res);
		 		rq_data.push(h);
			}
		}

		var result = {
			items   : subscriptions,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptUsersGet - Получаем пользователей по ID (execute)
func (vk *API) ScriptUsersGet(userIDs []string, fields string) (ans []UsersGetAns, err error) {
	// Разбиваем посты на нужное кол-во
	arr := chunkSliceString(userIDs, 1000)
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
		var arr    = %s;
		var ans    = [];

		while(arr.length > 0) {
			var str = arr.shift();
			var res = API.users.get({
				user_ids: str,
				fields: fields,
			});

			if(res) {
				ans = ans + res;
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

// ScriptMultiUsersGet - Получаем пользователей по ID (execute)
func (vk *API) ScriptMultiUsersGet(arr []map[string]interface{}) (ans ScriptUsersMultiGetAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr     = %s;
		var rq_data = [];
		var users   = [];

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.users.get({
				user_ids : h.user_ids,
				fields   : h.fields,
			});

			if(res) {
				users = users + res;
				rq_data.push(h);
			}
		}

		var result = {
			items   : users,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptMultiMarketGet - Получаем товары (execute)
func (vk *API) ScriptMultiMarketGet(arr []map[string]interface{}) (ans ScriptMultiMarketGetAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr     = %s;
		var rq_data = [];
		var goods   = [];
		var limit   = 200;

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.market.get({
				owner_id : h.owner_id,
				extended : 1,
				count    : limit,
			});

			if(res) {
				res.offset = limit;
				goods.push(res);
				rq_data.push(h);
			}
		}

		var result = {
			items   : goods,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptMarketGet - Получаем товары сообщества или пользователя. (execute)
func (vk *API) ScriptMarketGet(ownerID, offset int) (ans MarketGetAns, err error) {

	script := fmt.Sprintf(`
		var owner_id = %d;
		var offset   = %d;
	
		var cnt   = 25;
		var count = offset + 1;
		var goods = [];
		var limit = 200;

		while(cnt > 0 && offset < count){
			var res = API.market.get({ 
				owner_id   : owner_id,
				offset     : offset,
				count      : limit,
				extended   : 1,
			}); 
			cnt = cnt - 1;

			if(res.count) {
				count  = res.count; 
				goods = goods + res.items;
				offset = offset + limit;
			}
			else {
				cnt = 0;
				if(count == offset + 1) {
					count = offset;
				}
			}
		}

		var result = {
			count	 : count,
			offset : offset,
			items  : goods
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
		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

// ScriptMultiMarketGetComments - Получаем комментарии нескольких фото. (execute)
func (vk *API) ScriptMultiMarketGetComments(arr []map[string]interface{}) (ans MultiMarketGetCommentsAns, err error) {
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	script := fmt.Sprintf(`
		var arr      = %s;
		var rq_data  = [];
		var comments = [];
		var limit    = 100;

		while(arr.length > 0) {
			var h   = arr.shift();
			var res = API.market.getComments({ 
				owner_id         : h.owner_id,
				item_id          : h.item_id,
				sort             : "desc",
				need_likes       : 1,
				count            : limit
			});

			if(res.count) {
				res.offset = limit;
				comments.push(res);
		 		rq_data.push(h);
			}
		}

		var result = {
			items   : comments,
			rq_data : rq_data,
		};
		
		return result;
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

// ScriptMarketGetComments - Получаем комментарии товара. (execute)
func (vk *API) ScriptMarketGetComments(ownerID, itemID, startCommentID int) (ans WallGetCommentsAns, err error) {

	script := fmt.Sprintf(`
		var owner_id         = %d;
		var item_id          = %d;
		var start_comment_id = %d;

		var cnt         = 25;
		var real_offset = 0;
		var offset      = 0;
		var count       = offset + 1;
		var comments    = [];
		var limit       = 100;

		while(cnt > 0 && real_offset < count){
			var res = API.wall.getComments({ 
				owner_id         : owner_id, 
				item_id          : item_id,
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
	`, ownerID, itemID, startCommentID)

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
