package yourls

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// UserData is a struct representing the request with user data
type UserData struct {
	Action  string
	URL     string
	Filter  string
	Limit   string
	token   string
	baseURL string
	format  string
}

// YourlsData is a struct representing the response return data Yourls.
type YourlsData struct {
	Args struct {
		Keyword string
		URL     string
		Title   string
		Date    string
		IP      string
		Clicks  string
	} `json:"url"`
	Links      map[string]link
	Link       link
	Stats      totalStats `json:"stats"`
	DBStats    totalStats `json:"db-stats"`
	Status     string
	LongURL    string
	Message    string
	Title      string
	Shorturl   string
	Code       string
	ErrorCode  int
	StatusCode int
}

type totalStats struct {
	TLinks  string `json:"total_links"`
	TClicks string `json:"total_clicks"`
}

type link struct {
	Shorturl  string
	URL       string
	Title     string
	Timestamp string
	IP        string
	Clicks    string
}

// New initiated and returned only format(json) signature "YOURS" and baseURL
func New(token, baseURL string) *UserData {
	return &UserData{
		format:  "json",
		baseURL: baseURL,
		token:   token,
	}
}

// SetData set user data
func (u *UserData) SetData(data map[string]string) {
	u.Action = data["action"]
	u.Filter = data["filter"]
	u.Limit = data["limit"]
	u.URL = data["url"]
}

func (u *UserData) prepareRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", u.baseURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("format", u.format)
	q.Add("action", u.Action)
	q.Add("signature", u.token)

	switch u.Action {
	case "shorturl":
		q.Add("url", u.URL)
	case "expand", "url-stats":
		q.Add("shorturl", u.URL)
	case "stats":
		{
			q.Add("filter", u.Filter)
			q.Add("limit", u.Limit)
		}
	case "db-stats":
		q.Add("db-stats", "")
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
}

// Get receives data from your YourlsData
func (u *UserData) Get() (*YourlsData, error) {
	r, err := u.prepareRequest()
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var y YourlsData
	if err = json.Unmarshal(body, &y); err != nil {
		return nil, err
	}
	resp.Body.Close()
	return &y, nil
}
