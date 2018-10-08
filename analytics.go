package vkapi

var (
	statRqCollect  bool
	statRqQueueLen int
	statRqChan     chan RqStatObj
)

// RqStatObj - объект статистики запроса
type RqStatObj struct {
	Method  string
	Error   error
	Timeout int64
}

// InitStatChan - Создаем канал для аналитики
func InitStatChan(l int) chan RqStatObj {

	statRqChan = make(chan RqStatObj, l)
	statRqQueueLen = int(float64(l) * 0.1)
	statRqCollect = true

	return statRqChan
}
