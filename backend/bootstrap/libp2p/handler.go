package libp2p

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
)

func (bs *BootstrapServer) ListPeersHandler() server.Handler {
	return server.NewHandler("listBoostrapPeers", "/cgi_bin/bootstrap",
		func(c context.Context, ctx *app.RequestContext) {
			listPeers := bs.ListPeers(c)
			fmt.Printf("Connected listPeers (%d):\n", len(listPeers))
			for _, peer := range listPeers {
				fmt.Println(" -", peer)
			}

			ctx.String(http.StatusOK, "hello world, from hertz handler")
		},
	)
}
