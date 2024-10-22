package alerts

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Ip          string
	Trigger     string
	TriggerPort int
	AlertLink   string
	HostLink    string
	Desc        string
	Timestamp   time.Time

	Loaded  bool
	AlertId string
	Name    string
	Ports   map[int][]Cve
}

func NewEventFromItem(item Item) Event {
	splitTitle := strings.Split(item.Title, " ")
	ip := splitTitle[0]
	trigger := strings.ReplaceAll(splitTitle[len(splitTitle)-1], "`", "")
	port, _ := strconv.Atoi(splitTitle[3])

	timestamp, err := time.Parse("Sun, 20 Oct 2024 14:33:09 +0000", item.PubDate)
	if err != nil {
		timestamp = time.Now()
	}

	return Event{
		Ip:          ip,
		Trigger:     trigger,
		TriggerPort: port,
		AlertLink:   item.Link,
		HostLink:    "https://www.shodan.io/host/" + ip,
		Desc:        item.Description + " on port " + string(port),
		Timestamp:   timestamp,
		Ports:       make(map[int][]Cve),
		Loaded:      false,
	}
}

func (e *Event) Load() Event {
	if e.Loaded == true {
		return *e
	}

	go func(e *Event) {
		e.getAlertId(0)
		e.getName(0)
	}(e)

	bannerChannel := make(chan Banner)

	go func(e *Event, ch chan Banner) {
		banner := e.getBanner(0)
		ch <- banner
	}(e, bannerChannel)

	banner := <-bannerChannel
	for _, d := range banner.Data {
		for name, v := range d.Vulns {
			cve := NewCve(name, v)
			e.Ports[d.Port] = append(e.Ports[d.Port], cve)
		}
	}

	e.Loaded = true
	return *e
}

func (e *Event) getAlertId(retries int) {
	url := e.AlertLink + "?key=" + os.Getenv("API_KEY")
	response, err := http.Get(url)
	if err != nil {
		e.AlertId = "Could not get AlertID" + err.Error()
		return
	}

	if response.StatusCode == http.StatusTooManyRequests {
		if retries == 5 {
			e.AlertId = fmt.Sprintf("http response error: %s", response.Status)
		} else {
			time.Sleep(time.Second * time.Duration((retries + 1)))
			e.getAlertId(retries + 1)
		}
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		e.AlertId = fmt.Sprintf("Error reading alert id: %s", err.Error())
		return
	}

	splitData := strings.Split(string(body), "let data =")
	e.AlertId = strings.Split(splitData[1], "\"")[3]
}

func (e *Event) getName(retries int) {
	url := "https://api.shodan.io/shodan/alert/" + e.AlertId + "/info?key=" + os.Getenv("API_KEY")
	response, err := http.Get(url)
	if err != nil {
		e.Name = "Could not get name: " + err.Error()
		return
	}

	if response.StatusCode == http.StatusTooManyRequests {
		if retries == 5 {
			e.Name = fmt.Sprintf("http response error: %s", response.Status)
		} else {
			time.Sleep(time.Second * time.Duration((retries + 1)))
			e.getName(retries + 1)
		}
		return
	}
	if response.StatusCode != http.StatusOK {
		e.Name = fmt.Sprintf("http response error: %s", response.Status)
		return
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	alertString := make(map[string]string)
	json.Unmarshal(body, &alertString)

	e.Name = alertString["name"]
}

func (e *Event) getBanner(retries int) Banner {
	url := "https://api.shodan.io/shodan/host/" + e.Ip + "?key=" + os.Getenv("API_KEY")
	response, err := http.Get(url)
	if err != nil {
		return Banner{}
	}

	if response.StatusCode == http.StatusTooManyRequests {
		if retries == 5 {
			return Banner{}
		} else {
			time.Sleep(time.Second * time.Duration((retries + 1)))
			return e.getBanner(retries + 1)
		}
	}
	if response.StatusCode != http.StatusOK {
		return Banner{}
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	banner := Banner{}

	json.Unmarshal(body, &banner)

	return banner
}

type Vuln struct {
	Cvss   float32 `json:"cvss,omitempty"`
	CvssV2 float32 `json:"cvss_v2,omitempty"`
	Epss   float32 `json:"epss,omitempty"`
	Kev    bool    `json:"kev,omitempty"`
}

type Banner struct {
	Data []struct {
		Port    int             `json:"port"`
		Vulns   map[string]Vuln `json:"vulns,omitempty"`
		Product string          `json:"product,omitempty"`
	} `json:"data"`
}

type Item struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Guid        struct {
		Text string `xml:",chardata"`
	} `xml:"guid"`
	PubDate string `xml:"pubDate"`
}

type Rss struct {
	Channel struct {
		Item []Item `xml:"item"`
	} `xml:"channel"`
}

func Download() []Event {
	cache := NewEventCache()

	apiKey := os.Getenv("API_KEY")
	response, _ := http.Get("https://monitor.shodan.io/events.rss?key=" + apiKey)

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", response.StatusCode)
		return []Event{}
	}
	defer response.Body.Close()

	var rss Rss

	decoder := xml.NewDecoder(response.Body)
	decoder.Decode(&rss)

	events := []Event{}

	for _, item := range rss.Channel.Item {
		newEvent := NewEventFromItem(item)
		if cache.HasEventBeenSeen(newEvent) == false {
			events = append(events, newEvent)
			cache.InsertEvent(newEvent)
		}
	}

	return events
}
