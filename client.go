package ghost

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func EnableMonitor(url, name string) error {
	return (&ghost{
		url:  url,
		http: http.DefaultClient,
	}).EnableMonitor(name)
}

func DisableMonitor(url, name string) error {
	return (&ghost{
		url:  url,
		http: http.DefaultClient,
	}).DisableMonitor(name)
}

func Monitors(url string) (map[string]bool, error) {
	return (&ghost{
		url:  url,
		http: http.DefaultClient,
	}).Monitors()
}

type ghost struct {
	url  string
	http *http.Client
}

func (g *ghost) EnableMonitor(name string) error {
	return g.action(ActionEnable, name)
}

func (g *ghost) DisableMonitor(name string) error {
	return g.action(ActionDisable, name)
}

func (g *ghost) Monitors() (map[string]bool, error) {
	resp, err := g.http.Get(g.url)
	if err != nil {
		return nil, errors.Wrap(err, "error doing monitors requst")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	monitors := make(map[string]bool)
	if err := json.NewDecoder(resp.Body).Decode(&monitors); err != nil {
		return nil, errors.Wrap(err, "error decoding monitors response body")
	}
	return monitors, nil
}

func (g *ghost) action(action Action, name string) error {
	var vals url.Values = make(map[string][]string)
	vals.Set("action", string(action))
	vals.Set("name", name)
	resp, err := g.http.PostForm(g.url, vals)
	if err != nil {
		return errors.Wrapf(err, "error doing %s request", action)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected response status code %d", resp.StatusCode)
	}
	return nil
}
