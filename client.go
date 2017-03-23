package ghost

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// EnableMonitor enables a remote monitor 'name' located at url.
func EnableMonitor(url, name string) error {
	return (&ghost{
		url:  url,
		http: http.DefaultClient,
	}).EnableMonitor(name)
}

// DisableMonitor disables a remote monitor 'name' located at url.
func DisableMonitor(url, name string) error {
	return (&ghost{
		url:  url,
		http: http.DefaultClient,
	}).DisableMonitor(name)
}

// Monitors returns all registered monitors and their respective states.
func Monitors(url string) (map[string]MonitorState, error) {
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

func (g *ghost) Monitors() (map[string]MonitorState, error) {
	resp, err := g.http.Get(g.url)
	if err != nil {
		return nil, errors.Wrap(err, "error doing monitors requst")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	monitors := make(map[string]MonitorState)
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
