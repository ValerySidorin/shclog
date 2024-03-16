# slog: hclog adapter

Allows to use [slog](https://pkg.go.dev/log/slog) as an implementation of [hclog](https://github.com/hashicorp/go-hclog). Useful, when using hashicorp libraries, for example [raft](https://github.com/hashicorp/raft), or other projects, adopted [hclog](https://github.com/hashicorp/go-hclog) interface as their default logger.

```
go get github.com/ValerySidorin/shclog
```

## 💡 Usage

### Hashicorp raft custom logging

```go
import (
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/ValerySidorin/shclog"
	"github.com/hashicorp/raft"
	"github.com/lmittmann/tint"
)

raftCfg := raft.DefaultConfig()
raftCfg.LocalID = raft.ServerID("node_1")

w := os.Stderr
slogger := slog.New(tint.NewHandler(w, &tint.Options{
	AddSource: true,
	Level:     shclog.SlogLevelTrace,
}))
l := shclog.New(slogger)

// Use shclog here
raftCfg.Logger = l

addr, err := net.ResolveTCPAddr("tcp", "localhost:9876")
if err != nil {
	l.Error(err.Error())
	os.Exit(1)
}

transport, err := raft.NewTCPTransport("localhost:9876", addr, 3, 10*time.Second, os.Stderr)
if err != nil {
	l.Error(err.Error())
	os.Exit(1)
}

_, err = raft.NewRaft(raftCfg, nil, // Pass your fsm implementation here
	raft.NewInmemStore(), raft.NewInmemStore(), raft.NewDiscardSnapshotStore(),
	transport)
if err != nil {
	l.Error(err.Error())
	os.Exit(1)
}
```

### Default logging

```go
import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ValerySidorin/shclog"
	"github.com/hashicorp/go-hclog"
)

level := shclog.SlogLevelTrace
slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
	Level:     level,
}))
l := shclog.New(slogger)

l.Log(hclog.Error, "error log", "key", "value")
l.Log(hclog.Off, "no log")
l.Info("info log")
l.Trace("debug log")

fmt.Println("Log level:", l.GetLevel())
fmt.Println("Log is trace:", l.IsTrace())
fmt.Println("Log is debug:", l.IsDebug())
fmt.Println("Log name:", l.Name())

l = l.Named("name_1")
fmt.Println("Log name:", l.Name())

l = l.With("key", "value")
fmt.Println("Implied args:", l.ImpliedArgs())

l = l.Named("name_2")
fmt.Println("Log name:", l.Name())
l = l.ResetNamed("name")
fmt.Println("Log name:", l.Name())

l.SetLevel(hclog.Debug)
fmt.Println("Level can not be set through shclog. Instead set it through slog handler. Log level:", l.GetLevel())

stdLog := l.StandardLogger(&hclog.StandardLoggerOptions{})
stdLog.Println("log from std logger")

w := l.StandardWriter(&hclog.StandardLoggerOptions{})
w.Write([]byte("log from std writer"))

```