package speechkit

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	// RU Russian language
	RU string = "ru-RU"
	// EN English language
	EN string = "en-US"
	// TR Turkish language
	TR string = "tr-TR"
	// UK Ukrainian language
	UK string = "uk-UK"
)

const (
	defaultHost   = "tts.voicetech.yandex.net"
	defaultPath   = "/generate"
	defaultScheme = "https"
	defaultMethod = "GET"
)

//Client for https://tech.yandex.ru/speechkit/cloud/
type Client struct {
	ApiKey  string
	Lang    string
	Format  string
	Speaker string
	Emotion string
	Cl      *http.Client
}

//SaveToAudio - get audio file and save it
//@text - text to speech
//@path - path of save file with name file
func (c *Client) SaveToAudio(text, path string, perm os.FileMode) error {
	b, err := c.Get(text)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, b, perm); err != nil {
		return err
	}
	return nil
}

func (c *Client) Get(text string) ([]byte, error) {
	req, err := c.makeReq(text)
	if err != nil {
		return nil, err
	}
	res, err := c.Cl.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return b, nil
}

func (c *Client) makeReq(text string) (*http.Request, error) {
	values := url.Values{}
	values.Set("text", text)
	values.Set("format", c.Format)
	values.Set("lang", c.Lang)
	values.Set("speaker", c.Speaker)
	values.Set("key", c.ApiKey)
	values.Set("emotion", c.Emotion)
	u := url.URL{}
	u.Host = defaultHost
	u.Scheme = defaultScheme
	u.Path = defaultPath
	u.RawQuery = values.Encode()
	req, err := http.NewRequest(defaultMethod, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

//DefaultClient lang: ru, format: mp3, speaker: ermil
func DefaultClient(apiKey string) *Client {
	cl := &http.Client{Timeout: time.Second * 20}
	return &Client{ApiKey: apiKey, Lang: RU, Format: "mp3", Speaker: "ermil",
		Emotion: "good", Cl: cl}
}
