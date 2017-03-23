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
	// ActionEnable implies enabling a monitor.
	ActionEnable Action = "enable"
	// ActionDisable implies disabling a monitor.
	ActionDisable Action = "disable"
)

// MonitorState describes state of a monitor.
type MonitorState bool

var (
	// MonitorStateEnabled implies that a monitor is enabled.
	MonitorStateEnabled = true
	// MonitorStateDisabled implies that a monitor is disabled.
	MonitorStateDisabled = false
)

var initOnce sync.Once // gaurds monitors map initialization
var mu sync.Mutex      // guards
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

// MonitorFuncs is an adaptor func to allow use of ordinary functions as
// a Monitor.
func MonitorFuncs(name string, enable, disable func(), state func() bool) Monitor {
	return &funcMonitor{
		name:    name,
		enable:  enable,
		disable: disable,
		state:   state,
	}
}

type funcMonitor struct {
	name    string
	enable  func()
	disable func()
	state   func() bool
}

func (m *funcMonitor) Enable()       { m.enable() }
func (m *funcMonitor) Disable()      { m.disable() }
func (m *funcMonitor) Enabled() bool { return m.state() }
func (m *funcMonitor) Name() string  { return m.name }

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

	mons := make(map[string]bool, len(monitors))
	for name, monitor := range monitors {
		mons[name] = monitor.Enabled()
	}

	if err := json.NewEncoder(w).Encode(&mons); err != nil {
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
