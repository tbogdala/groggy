// Copyright 2015, Timothy Bogdala <tdb@animal-machine.com>
// See the LICENSE file for more details.

package groggy

import (
	"fmt"
	"testing"
	"sync"
	"time"
)


func TestDefaultLogger(t *testing.T) {
	// use the default handler
	Register("defaultTest", nil)

	// send a test message
	err := Log("defaultTest", "This is a test of: ", "Hello World!")
	if err != nil {
		t.Log("Log failed when it shouldn't have.\n")
		t.Fail()
	}

	// remove the handler
	Deregister("defaultTest")
}

func TestDefaultSyncLogger(t *testing.T) {
	// use the default handler
	Register("SYNC", DefaultSyncHandler)
	var wg sync.WaitGroup

	// send a test message
	for i:=0; i<42; i++ {
		wg.Add(1)
		go func(i int) {
			Log("SYNC", fmt.Sprintf("This is a synchronized test number %d!", i))
			wg.Done()
		}(i)
	}
	wg.Wait()

	// remove the handler
	Deregister("SYNC")
}

func TestNoLogger(t *testing.T) {
	// use the default handler
	Register("defaultTest", nil)

	// send a test message
	err := Log("NOT_THERE", "This is a test of: ", "Hello World!")
	if err == nil {
		t.Log("Log() should have failed with a non-registered log name.\n")
		t.Fail()
	}

	// remove the handler
	Deregister("defaultTest")
}


type EventTest struct {
	Data int
	Desc string
}

func LogEventHandler(logName string, data ...interface{}) error {
	const layout = "15:04:05.000"
	now := time.Now()
	for _, d := range data {
		switch dt := d.(type) {
		case *EventTest:
			fmt.Printf("%s %s: %s %d\n", now.Format(layout), logName, dt.Desc, dt.Data)
		default:
			fmt.Printf("<unknown log data type %v>",dt)
		}
	}
	return nil
}


func TestEventLogger(t *testing.T) {
	// use the default handler
	Register("EVENT", LogEventHandler)

	var e EventTest
	e.Data = 42
	e.Desc = "The answer is"
	err := Log("EVENT", &e)
	if err != nil {
		t.Log("Log() returned an error for a custom object.")
		t.Fail()
	}

	// remove the handler
	Deregister("EVENT")
}
