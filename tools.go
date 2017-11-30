package vkapi

import (
	"strings"
)

// Разбиваем массив строк на несколько
func chunkSliceString(arr []string, size int) (ans [][]string) {
	msize := len(arr) / size
	if len(arr)%size != 0 {
		msize++
	}
	ans = make([][]string, msize)

	l := len(arr)
	now := 0
	i := 0
	for {
		next := now + size

		if now+size > l {
			next = l
		}

		ans[i] = arr[now:next]
		i++

		if next == l {
			break
		}
		now = next
	}

	return
}

// Проверяем надо ли пропустить ошибу
func (vk *Api) checkErrorSkip(str string) bool {
	for _, e := range vk.ErrorToSkip {
		if strings.Contains(str, e) {
			return true
		}
	}

	return false
}
