// Copyright 2015, Timothy Bogdala <tdb@animal-machine.com>
// See the LICENSE file for more details.

/*

Package groggy is a library that makes it easier to setup custom
logging channels.

To use:

1) Call Register() with a log name and an optional function to handle the log events.
2) Call Log() with this log name and the data objects to log

If no optional log handlers are supplied, the default handler writes the data
objects out to stdout using fmt.Print(). To do this, the objects should be
strings or implement the fmt.Stringer interface.

If Log() is called with a log name that is not registered, it will not be able
to call a handler, and an error will be returned.

Clients can call Deregister() to remove a log handler.

*/
package groggy

import (
	"fmt"
	"sync"
	"time"
)

// GroggyEvent defines a function handler for incoming logging events.
// It's setup with a variadic argument for extra flexibility in custom handlers.
type GroggyEvent func(logName string, data ...interface{}) error

var (
	// handlers is a global registry of event handlers
	handlers map[string]GroggyEvent

	// make a mutex for the default sync handler
	handlerMutex sync.Mutex
)

func init() {
	// make sure to initialize the global map
	handlers = make(map[string]GroggyEvent)
}

// Register adds a new log handler to the global registry and assigns it
// the handler function passed in. If handler is nil, then DefaultHandler is used.
// An existing log handler can be replaced using this function.
func Register(newLogName string, handler GroggyEvent) {
	// use DefaultHandler if a nil handler was supplied
	var h GroggyEvent = handler
	if h == nil {
		h = DefaultHandler
	}
	handlers[newLogName] = h
}

// Deregister removes the log handler from the global registry so that
// further calls to Log with the log name do not get handled.
func Deregister(logName string) {
	delete(handlers, logName)
}

// DefaultHandler writes out the information assuming data members are strings
// or anything that implements the GoStringer interface. This is not
// considered safe for concurrency.
func DefaultHandler(logName string, data ...interface{}) error {
	const layout = "15:04:05.000"
	now := time.Now()
	fmt.Printf("%s %s: ", now.Format(layout), logName)
	for _, ds := range data {
		switch v := ds.(type) {
		case string:
			fmt.Print(v)
		case fmt.Stringer:
			fmt.Print(v.String())
		default:
			fmt.Printf("<unknown log data type %v>", ds)
		}
	}
	fmt.Print("\n")
	return nil
}

// DefaultSyncHandler writes out the information assuming data members are strings
// or anything that implements the GoStringer interface. This is considered safe
// for concurrency.
func DefaultSyncHandler(logName string, data ...interface{}) error {
	handlerMutex.Lock()
	DefaultHandler(logName, data...)
	handlerMutex.Unlock()
	return nil
}

// Log sends the data to the handler specified by the logName. This is not
// considered safe for concurrency by default.
func Log(logName string, data ...interface{}) error {
	h, okay := handlers[logName]
	if okay == false {
		return fmt.Errorf("No log handler found for %s.", logName)
	}

	return h(logName, data...)
}

// Logsf uses the second parameter as the format string for fmt.Spritnf and
// then sends the rest of the parameters to it. The resulting string is then
// passed to the handler specified by logName.
func Logsf(logName string, data ...interface{}) error {
	h, okay := handlers[logName]
	if okay == false {
		return fmt.Errorf("No log handler found for %s.", logName)
	}

	sprintStr, okay := data[0].(string)
	if okay == false {
		return fmt.Errorf("A format string was not passed as the second parameter.")
	}

	s := fmt.Sprintf(sprintStr, data[1:]...)
	return h(logName, s)
}
