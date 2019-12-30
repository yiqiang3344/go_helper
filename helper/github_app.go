package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const clientId = "Iv1.c638d2dd2221ab39"
const secret = "a4ae6a088f9c8868ef41c95093541ffdb4bb33f6"
const IssuesURL = "https://api.github.com/repos/yiqiang3344/test/issues"

var Token string
var TimeLocal *time.Location

type CreateIssueParams struct {
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	Assignees []string `json:"assignees,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	State     string   `json:"state,omitempty"`
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	ClosedAt  time.Time `json:"closed_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

func GetToken(w http.ResponseWriter, r *http.Request, code string, state string, redirectUri string) error {
	if Token != "" {
		WriteLog("already has Token:"+Token, "sidneyyi.com/helper/github_app.go::GetToken")
		return nil
	}

	if code == "" {
		//跳转用户授权页面
		_url := "https://github.com/login/oauth/authorize?client_id=" + clientId + "&redirect_uri=" + redirectUri + "&state=" + state
		http.Redirect(w, r, _url, http.StatusFound)
		return fmt.Errorf("")
	}

	//获取token
	jsonStr := "{\"client_id\":\"" + clientId + "\",\"client_secret\":\"" + secret + "\",\"code\":\"" + code + "\"}"
	resp, err := http.Post("https://github.com/login/oauth/access_token", "application/json", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("request failed:%s", resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	data, _ := url.ParseQuery(string(body))
	if _, ok := data["access_token"]; !ok {
		WriteLog(fmt.Sprintf("failed:%s", string(body)), "sidneyyi.com/helper/github_app.go::GetToken")
		return fmt.Errorf("get token failed:%s", string(body))
	}

	WriteLog(fmt.Sprintf("success:%s", data["access_token"][0]), "sidneyyi.com/helper/github_app.go::GetToken")
	Token = data["access_token"][0]
	return nil
}

// SearchIssues queries the GitHub issue tracker.
func GetIssues() ([]Issue, error) {
	resp, err := http.Get(IssuesURL)
	if err != nil {
		return nil, err
	}

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result []Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()

	WriteLog(fmt.Sprintf("%#v", result), "sidneyyi.com/helper/github_app.go::GetIssues")

	return result, nil
}

func CreateIssues(data []byte) (*Issue, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", IssuesURL, bytes.NewBuffer(data))

	req.Header.Add("Authorization", "token "+Token)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	resp, _ := client.Do(req)

	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		return nil, fmt.Errorf("create failed: %s %s", resp.Status, string(data))
	}

	//body, _ := ioutil.ReadAll(resp.Body)
	//return nil, fmt.Errorf("data:%s\nresult: %s", data, string(body))

	var result Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

func UpdateIssues(id string, data []byte) (*Issue, error) {
	client := &http.Client{}

	req, err := http.NewRequest("PATCH", IssuesURL+"/"+id, bytes.NewBuffer(data))

	req.Header.Add("Authorization", "token "+Token)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	resp, _ := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("create failed: %s %s", resp.Status, string(data))
	}

	//body, _ := ioutil.ReadAll(resp.Body)
	//return nil, fmt.Errorf("data:%s\nresult: %s", data, string(body))

	var result Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}
