package timer

import (
	"encoding/json"
	"fmt"
	"github.com/zwk-app/go-tools/logs"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type TargetListJson struct {
	Items []struct {
		TargetTime  string `json:"targetTime"`
		TargetLabel string `json:"targetLabel"`
		AlarmSound  string `json:"alarmSound"`
	} `json:"targets"`
}

type TargetList struct {
	items []*TargetListItem
	last  struct {
		error error
	}
}

func NewTargetList() *TargetList {
	t := TargetList{}
	t.Append(NewTargetListItem("L1200", "Lunch Time", ""))
	t.Append(NewTargetListItem("L1800", "Let's go Home", ""))
	t.last.error = nil
	return &t
}

func (t *TargetList) Error() error {
	lastError := t.last.error
	t.last.error = nil
	return lastError
}

func (t *TargetList) NextTargetListItem() *TargetListItem {
	if len(t.items) > 0 {
		n := NewTime().Set(time.Now().In(time.Local)).TimeString()
		for _, r := range t.items {
			if n > r.timeString {
				continue
			} else {
				return r
			}
		}
		return t.items[0]
	}
	return NewTargetListItem("000000", "", "")
}

func (t *TargetList) String() string {
	s := ""
	for i, r := range t.items {
		s += fmt.Sprintf("[%02d] %6s %-16s %s\n", i+1, r.timeString, r.textLabel, r.alarmSound)
	}
	return s
}

func (t *TargetList) clear() *TargetList {
	t.items = nil
	t.last.error = nil
	return t
}

func (t *TargetList) sort() *TargetList {
	sort.Slice(t.items, func(i, j int) bool { return t.items[i].timeString < t.items[j].timeString })
	return t
}

func (t *TargetList) Append(i *TargetListItem) *TargetList {
	t.items = append(t.items, i)
	return t
}

func (t *TargetList) loadJson(jsonReader io.Reader) *TargetList {
	t.clear()
	jsonData := TargetListJson{}
	d := json.NewDecoder(jsonReader)
	t.last.error = d.Decode(&jsonData)
	if t.last.error == nil {
		for _, r := range jsonData.Items {
			t.Append(NewTargetListItem(r.TargetTime, r.TargetLabel, r.AlarmSound))
		}
	}
	t.sort()
	return t
}

func (t *TargetList) loadJsonData(jsonData string) *TargetList {
	t.loadJson(strings.NewReader(jsonData))
	return t
}

func (t *TargetList) loadJsonFile(jsonFileName string) *TargetList {
	f, e := os.Open(jsonFileName)
	if e != nil {
		t.last.error = e
		return t
	}
	t.loadJson(f)
	return t
}

func (t *TargetList) loadJsonURL(jsonURL string) *TargetList {
	u, e := url.Parse(jsonURL)
	c := http.Client{}
	r, e := c.Get(u.String())
	if e != nil {
		t.last.error = e
		return t
	}
	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()
	if r.StatusCode == 200 {
		t.loadJson(r.Body)
	} else {
		t.last.error = fmt.Errorf("HTTP %d %s", r.StatusCode, r.Status)
	}
	return t
}

func (t *TargetList) LoadJson(json string) *TargetList {
	if strings.HasPrefix(json, "{") && strings.HasSuffix(json, "}") {
		t.loadJsonData(json)
	} else if strings.HasPrefix(json, "http") {
		t.loadJsonURL(json)
	} else {
		t.loadJsonFile(json)
	}
	logs.Debug("TargetList", fmt.Sprintf("LoadJson: Items:\n%s", t.String()), nil)
	return t
}
