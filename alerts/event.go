package alerts

import (
	"encoding/xml"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Event struct {
	Ip        string
	Trigger   string
	AlertLink string
	HostLink  string
	Desc      string
	Timestamp time.Time

	AlertId   string
	Name string
	Ports map[string][]Cve
	// cves  []Cve
}

func NewEventFromItem(item Item) Event {
	splitTitle := strings.Split(item.Title, " ")
	ip := splitTitle[0]
	trigger := strings.ReplaceAll(splitTitle[len(splitTitle)-1], "`", "")
	port :=	splitTitle[3]

	timestamp, err := time.Parse("Sun, 20 Oct 2024 14:33:09 +0000", item.PubDate)
	if err != nil {
		timestamp = time.Now()
	}

	return Event{
		Ip:        ip,
		Trigger:   trigger,
		AlertLink: item.Link,
		HostLink:  "https://www.shodan.io/host/" + ip,
		Desc:      item.Description + " on port " + port,
		Timestamp: timestamp,
	}
}

func (e *Event) GetAlertId() {
	// url := 
}

func (e *Event) GetName(retries int) {
	url := "https://api.shodan.io/shodan/alert/" + e.AlertId + "/info?key=" + os.Getenv("API_KEY")
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		e.Name = "Could not get name: " +err.Error()
		return
	}

	if response.StatusCode == http.StatusTooManyRequests {
		if retries == 5 {
			e.Name = fmt.Sprintf("http response error: %s", response.Status)
			return
		} else {
			time.Sleep(time.Second * time.Duration((retries + 1)))
			e.GetName(retries + 1)
		}
		return
	}
	if response.StatusCode != http.StatusOK {
		e.Name = fmt.Sprintf("http response error: %s", response.Status)
		return
	}
	defer response.Body.Close()

	alertString := make(map[string]string)

	if err := json.NewDecoder(response.Body).Decode(&alertString); err != nil {
	    e.Name = "Error unmarshalling JSON: " + err.Error()
	    return
	}

	e.Name = alertString["name"]
}

type Cve struct {
	Name string
	Rank int
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
		events = append(events, NewEventFromItem(item))
	}

	return events
}
