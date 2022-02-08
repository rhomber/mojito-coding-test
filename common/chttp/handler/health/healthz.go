package health

import (
	"mojito-coding-test/common/chttp"
	"mojito-coding-test/common/config"
	"net/http"
	"time"
)

const (
	StatusOK   = "OK"
	StatusFAIL = "FAIL"

	ItemNameDbConn = "DB_CONN"
)

type Item struct {
	Name    string `json:"name" yaml:"name"`
	Status  string `json:"status" yaml:"status"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Time    string `json:"time" yaml:"time"`
}

func NewOverall(cfg *config.Config) *Overall {
	appKey := cfg.GetString("playground.app.key")
	if appKey == "" {
		appKey = cfg.GetString("app.key")
	}

	return &Overall{
		App:    appKey,
		Status: StatusOK,
		Items:  make(map[string]Item),
	}
}

type Overall struct {
	App    string          `json:"app" yaml:"app"`
	Status string          `json:"status" yaml:"status"`
	Items  map[string]Item `json:"items" yaml:"items"`
}

func (o *Overall) AddItem(item Item) {
	o.Items[item.Name] = item
}

func (o *Overall) SetStatus() {
	status := StatusOK
	for _, item := range o.Items {
		if item.Status != StatusOK {
			status = StatusFAIL
			break
		}
	}
	o.Status = status
}

func (o *Overall) GetHttpStatus() int {
	statusCode := http.StatusOK
	if o.Status != StatusOK {
		statusCode = http.StatusInternalServerError
	}

	return statusCode
}

func TimedCheck(ctx *chttp.Context, fun func(*chttp.Context) Item) Item {
	start := time.Now()
	item := fun(ctx)
	item.Time = time.Since(start).String()

	return item
}

func CheckDbConn() func(ctx *chttp.Context) Item {
	return func(ctx *chttp.Context) Item {
		item := Item{
			Name:   ItemNameDbConn,
			Status: StatusOK,
		}

		// TODO: Test Db Connection

		return item
	}
}
