package janitor

import (
	"net"
	"net/http"

	"github.com/armon/go-proxyproto"
	"golang.org/x/net/context"
)

type JanitorServer struct {
	config Config

	UpstreamLoader *UpstreamLoader
	EventChan      chan *TargetChangeEvent

	httpServer *http.Server
}

func NewJanitorServer(Config Config) *JanitorServer {
	server := &JanitorServer{
		config: Config,
	}

	server.EventChan = make(chan *TargetChangeEvent, 1024)
	server.UpstreamLoader = NewUpstreamLoader(server.EventChan)

	server.httpServer = &http.Server{Handler: NewHTTPProxy(&http.Transport{},
		server.config.HttpHandler,
		server.config.ListenAddr,
		server.UpstreamLoader)}

	return server
}

func (server *JanitorServer) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", server.config.ListenAddr)
	if err != nil {
		return err
	}

	go server.UpstreamLoader.Start(ctx)

	return server.httpServer.Serve(&proxyproto.Listener{Listener: TcpKeepAliveListener{ln.(*net.TCPListener)}})
}
