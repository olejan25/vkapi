package vkapi

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func (vk *Api) Script_Wall_GetById(posts []string) (ans []WallGetByIdAns, err error) {
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
			log.Println("[error]", err)
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

func (vk *Api) Script_Groups_GetById(groupIds []string, fields string) (ans []GroupsGetAns, err error) {
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
			log.Println("[error]", err)
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
