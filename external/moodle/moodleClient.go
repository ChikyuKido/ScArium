package moodle

import (
	"ScArium/common/log"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type MoodleClient struct {
	Client       *http.Client
	ServiceUrl   string
	Token        string
	PrivateToken string
}

func NewMoodleClient(serviceUrl string, username string, password string) *MoodleClient {
	client := &MoodleClient{
		ServiceUrl: serviceUrl,
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.Client = &http.Client{Transport: tr}
	err := client.login(username, password)
	if err != nil {
		log.E.Errorf("Failed to login to moodle (%s) with the username (%s): %v", serviceUrl, username, err)
		return nil
	}
	return client
}
func (mc *MoodleClient) login(username, password string) error {
	loginURL := fmt.Sprintf("%s/login/token.php", mc.ServiceUrl)
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("service", "moodle_mobile_app")
	req, err := http.NewRequest("POST", loginURL, nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = data.Encode()

	resp, err := mc.Client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response struct {
		Token        string `json:"token"`
		PrivateToken string `json:"privatetoken"`
		Error        string `json:"error"`
		ErrorCode    string `json:"errorcode"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshel json: %v", err)
	}

	if response.Error != "" {
		return fmt.Errorf("failed to obtain token: %s", response.Error)
	}

	mc.Token = response.Token
	mc.PrivateToken = response.PrivateToken
	return nil
}

func (mc *MoodleClient) makeRequest(function string, params map[string]string, url string) ([]byte, error) {
	webserviceURL := fmt.Sprintf("%s%s", mc.ServiceUrl, url)

	params["wstoken"] = mc.Token
	params["wsfunction"] = function
	params["moodlewsrestformat"] = "json"
	req, err := http.NewRequest("GET", webserviceURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := mc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (mc *MoodleClient) MakeWebserviceRequest(function string, params map[string]string) ([]byte, error) {
	return mc.makeRequest(function, params, "/webservice/rest/server.php")
}

func (mc *MoodleClient) MakeModRequest(function string, params map[string]string) ([]byte, error) {
	return mc.makeRequest(function, params, "/mod/assign/view.php")
}
func (mc *MoodleClient) DownloadFile(url string, path string) error {
	_, err := os.Stat(path)
	if err == nil {
		log.E.Info("File with the name already exists. Skipping it: %s", path)
		return fmt.Errorf("file %s already exists", path)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("token", mc.Token)
	req.URL.RawQuery = q.Encode()

	resp, err := mc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
