// Package ghost provides functionality for exposing application monitoring
// switches via HTTP handler.
package ghost

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Monitor represents a resource that could be (un)monitored.
type Monitor interface {
	// Name should return name describing the resource.
	Name() string
	// Enable should enable monitoring of the resource.
	Enable()
	// Disable should disable monitoring of the resource.
	Disable()
	// Enabled should return whether the resource is being monitored.
	Enabled() bool
}

// Action describes an action regarding a monitor.
type Action string

var (
	ActionEnable  Action = "enable"
	ActionDisable Action = "disable"
)

var initOnce sync.Once
var mu sync.Mutex // guards
var monitors map[string]Monitor

// RegisterMonitor registers m with the set of available monitors.
func RegisterMonitor(m Monitor) {
	initOnce.Do(func() {
		monitors = make(map[string]Monitor)
	})

	mu.Lock()
	defer mu.Unlock()
	monitors[m.Name()] = m
}

// MonitorHandler returns http.Handler that handles (un)monitor requests.
func MonitorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			listMonitors(w, req)
		case http.MethodPost:
			switchMonitor(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func listMonitors(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if err := json.NewEncoder(w).Encode(&monitors); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func switchMonitor(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := req.FormValue("name")
	mu.Lock()
	monitor, ok := monitors[name]
	mu.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch Action(req.FormValue("action")) {
	case ActionEnable:
		monitor.Enable()
		w.WriteHeader(http.StatusOK)
	case ActionDisable:
		monitor.Disable()
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}
