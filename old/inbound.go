package api

import (
	"net/rpc"
	"net"
	"net/rpc/jsonrpc"
    "github.com/backer/bacsrv/config"
    log "github.com/Sirupsen/logrus" 
	"os"

)


func init(){
     // Log as JSON instead of the default ASCII formatter.
  log.SetFormatter(&log.JSONFormatter{})

  // Output to stderr instead of stdout, could also be a file.
  log.SetOutput(os.Stderr)

  // Only log the warning severity or above.
  log.SetLevel(log.DebugLevel)
}
type Args struct {
	A, B int
}

type Arith int

type Result int

func (t *Arith) Multiply(args *Args, result *Result) error {
	log.Printf("Multiplying %d with %d\n", args.A, args.B)
	*result = Result(args.A * args.B)
	return nil
}


// StartInboundInterface is able to serve an interface in seperated goroutine
func StartInboundInterface(){
    log.Debug("Starting inboud interface...")
    config := config.ReadConfigFile()
    server := rpc.NewServer()
    arith := new(Arith)
    server.Register(arith)
    server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
    l, e := net.Listen("tcp", ":" + config.MgmtPort)
    defer l.Close()
    if e != nil {
        log.Fatal("Listen error: ", e)
    }
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go server.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}