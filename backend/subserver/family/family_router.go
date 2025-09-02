package family

import (
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
)

const (
	RouterURLFamilyPhotoSync RouterPath = "/family/photo/sync"
	RouterURLFamilyPhotoList RouterPath = "/family/photo/list"
	RouterURLFamilyPhotoGet  RouterPath = "/family/photo/get"
)

const (
	RoutersNameFamily = "family"
)

type RouterPath string

func (rp RouterPath) Name() string {
	return string(rp)
}

func (rp RouterPath) SubPath() string {
	return string(rp)
}

// FamilyRouters provides family photo endpoints for the service
type FamilyRouters struct{}

// Ensure FamilyRouters implements server.Routers interface
var _ server.Routers = (*FamilyRouters)(nil)

// Handlers registers all family-related handlers
func (fr *FamilyRouters) Handlers() []server.Handler {
	handlerInfos := GetFamilyHandlers()
	handlers := make([]server.Handler, len(handlerInfos))

	for i, info := range handlerInfos {
		handlers[i] = server.NewHandler(
			info.RouterURL,
			info.Handler,
			server.WithMethod(info.Method),
			server.WithWrappers(info.Wrappers...),
		)
	}

	return handlers
}

func (fr *FamilyRouters) Name() string {
	return RoutersNameFamily
}

// NewFamilyRouter creates a new family router instance
func NewFamilyRouter() server.Routers {
	return &FamilyRouters{}
}