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
		var cnt      = 25;

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
			users	 : users
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
		var cnt     = 25;

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
			users	 : users
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
		var cnt     = 25;

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
			users	 : users
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
