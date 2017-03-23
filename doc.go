/*
Package ghost provides functionality for exposing application monitoring
switches via HTTP handler.

Its idea is to provide primitives for enabling or disabling monitoring for parts
of the system during runtime.

To get started, the service needs to register a Monitor object which provides
enable, disable and state reporting functionality.

	var running bool
	runtimeMonitor := ghost.MonitorFuncs("runtime",
		func() {
			// Turn on the monitor.
			running = true
		},
		func() {
			// Turn off the monitor.
			running = false
		}, func() bool {
			// Report whether the monitor is enabled.
			return running
		})
	ghost.RegisterMonitor(runtimeMonitor)

If a monitor needs to have more sophisticated state it could be wrapped in a
struct that implements the Monitor interface. After registering all monitors,
the remote enabling functionality should be
exposed by registering the MonitorHandler.

	http.Handle("/monitors", ghost.MonitorHandler())
	http.ListenAndServe("localhost:8080", nil)

After starting the application, listing the available monitor targets and the
(un)monitor actions could be triggered via the ghost CLI.

	ghost -g http://localhost:8080/monitors targets
	Name		Status
	runtime		disabled

	ghost -g http://localhost:8080/monitors monitor runtime
	ghost -g http://localhost:8080/monitors targets
	Name		Status
	runtime		enabled
*/
package ghost // import "github.com/Bo0mer/ghost"
