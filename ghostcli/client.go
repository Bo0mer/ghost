package ghostcli

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Bo0mer/ghost"
	"github.com/pkg/errors"
)

func EnableMonitor(url, name string) error {
	return (&Ghost{
		url:  url,
		http: http.DefaultClient,
	}).EnableMonitor(name)
}

func DisableMonitor(url, name string) error {
	return (&Ghost{
		url:  url,
		http: http.DefaultClient,
	}).DisableMonitor(name)
}

func Monitors(url string) (map[string]bool, error) {
	return (&Ghost{
		url:  url,
		http: http.DefaultClient,
	}).Monitors()
}

type Ghost struct {
	url  string
	http *http.Client
}

func (g *Ghost) EnableMonitor(name string) error {
	return g.action(ghost.ActionEnable, name)
}

func (g *Ghost) DisableMonitor(name string) error {
	return g.action(ghost.ActionDisable, name)
}

func (g *Ghost) Monitors() (map[string]bool, error) {
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

func (g *Ghost) action(action ghost.Action, name string) error {
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
