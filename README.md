groggy: a simple Go library creating custom logging routines
============================================================

Groggy is a library that makes it easier to setup custom
logging channels.

At the most basic level, this can be done by calling
`Register()` with a log name and then calling `Log()` with that name
as a parameter.

```go
groggy.Register("DEBUG", nil)
err := groggy.Log("DEBUG", "This is a test of: ", "Hello World!")
if err != nil {
  panic()
}
groggy.Deregister("DEBUG")
```

This will output something like:

```bash
14:32:07.197 DEBUG: This is a test of: Hello World!
```

Groggy manages a non-exported map that will be used to dispatch events to
handlers when `Log()` is called, which means that client code doesn't need
to maintain any object pointer to the log itself. This makes it much easier
to call from any code in the project.

Also, determining which log names you register will determine what log events
get passed on. If you call `Log("DEBUG", ...)` a lot, but then don't register
it at the start of the program, your DEBUG messages won't go through to a handler --
a convenient way to not show any DEBUG messages in a production environment.


Installation
------------

You can get the latest copy of the library using this command:

```bash
go get github.com/tbogdala/groggy
```

You can then use it in your code by importing it:

```go
import "github.com/tbogdala/groggy"
```


Usage
-----

Two basic steps are required ot use this package:

1. Call Register() with a log name and an optional function to handle the log events.
2. Call Log() with this log name and the data objects to log

If no optional log handlers are supplied, the default handler writes the data
objects out to stdout using fmt.Print(). To do this, the objects should be
strings or implement the fmt.Stringer interface.

If Log() is called with a log name that is not registered, it will not be able
to call a handler, and an error will be returned.

Clients can call Deregister() to remove a log handler.

Besides the basic `DefaultHandler` function, a `DefaultSyncHandler` function
is supplied as a drop-in replacement that locks a mutex during writes as
an example of how to handle events in a synchronous way.

Since the event handler map is not protected by `sync` objects, it is not
considered safe to call `Log()` concurrently with possible calls to
`Register()` and `Deregister()`

The following is a sample based on a test case that uses a custom handler:

```go
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

func TestEventLogger() {
  Register("EVENT", LogEventHandler)
  var e EventTest
  e.Data = 42
  e.Desc = "The answer is"
  Log("EVENT", &e)
}
```

License
-------

Groggy is released under the BSD license. See the `LICENSE` file for more details.
