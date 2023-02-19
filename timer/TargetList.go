package timer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

type TargetListJson struct {
	Items []struct {
		TargetTime  string `json:"targetTime"`
		TargetLabel string `json:"targetLabel"`
		AlertSound  string `json:"alertSound"`
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
	t.clear()
	return &t
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

func (t *TargetList) Error() error {
	lastError := t.last.error
	t.last.error = nil
	return lastError
}

func (t *TargetList) String() string {
	s := ""
	for i, r := range t.items {
		s += fmt.Sprintf("[%02d] %6s %-16s %s\n", i+1, r.timeString, r.textLabel, r.alertSound)
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
	if t.last.error != nil {
		log.Printf("TargetList->loadJson error: %s", t.last.error.Error())
	} else {
		for _, r := range jsonData.Items {
			t.Append(NewTargetListItem(r.TargetTime, r.TargetLabel, r.AlertSound))
		}
	}
	t.sort()
	log.Printf("TargetList->items:\n%s", t.String())
	return t
}

func (t *TargetList) LoadJsonFile(jsonFileName string) *TargetList {
	if jsonFileName == "" {
		return t
	}
	log.Printf("TargetList->LoadJsonFile: '%s'", jsonFileName)
	f, e := os.Open(jsonFileName)
	if e != nil {
		t.last.error = e
		log.Printf("TargetList->LoadJsonFile error: %s", t.last.error.Error())
		return t
	}
	t.loadJson(f)
	if t.last.error != nil {
		log.Printf("TargetList->LoadJsonFile error: %s", t.last.error.Error())
	}
	return t
}

func (t *TargetList) LoadJsonURL(jsonURL string) *TargetList {
	if jsonURL == "" {
		return t
	}
	log.Printf("TargetList->LoadJsonURL '%s'", jsonURL)
	u, e := url.Parse(jsonURL)
	c := http.Client{}
	r, e := c.Get(u.String())
	if e != nil {
		t.last.error = e
		log.Printf("TargetList->LoadJsonURL error: %s", t.last.error.Error())
		return t
	}
	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()
	if r.StatusCode == 200 {
		t.loadJson(r.Body)
		if t.last.error != nil {
			log.Printf("TargetList->LoadJsonURL error: %s", t.last.error.Error())
		}
		return t
	} else {
		t.last.error = fmt.Errorf("HTTP %d %s", r.StatusCode, r.Status)
		log.Printf("TargetList->LoadJsonURL error: %s", t.last.error.Error())
	}
	return t
}
