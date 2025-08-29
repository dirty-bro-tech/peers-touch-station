package family

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/server"
)

var (
	_ server.Subserver = (*PhotoSaveSubServer)(nil)
)

// familyRouterURL implements server.RouterURL for family endpoints
type familyRouterURL struct {
	name string
	url  string
}

func (f familyRouterURL) Name() string {
	return f.name
}

func (f familyRouterURL) URL() string {
	return f.url
}

// PhotoSaveSubServer handles photo upload requests
type PhotoSaveSubServer struct {
	opts *Options

	addrs  []string      // Populated from configuration
	status server.Status // Track server status
}

// Name returns the subserver identifier
func (s *PhotoSaveSubServer) Name() string {
	return "photo-save"
}

// Type returns the subserver type (HTTP in this case)
func (s *PhotoSaveSubServer) Type() server.SubserverType {
	return server.SubserverTypeHTTP
}

// Address returns the listening addresses
func (s *PhotoSaveSubServer) Address() server.SubserverAddress {
	return server.SubserverAddress{
		Address: s.addrs,
	}
}

// Handlers defines the upload endpoint
func (s *PhotoSaveSubServer) Handlers() []server.Handler {
	return []server.Handler{
		server.NewHandler(
			familyRouterURL{name: "sync", url: "/family/photo/sync"},
			s.handlePhotoUpload,            // Handler function
			server.WithMethod(server.POST), // HTTP method
		),
	}
}

// Init initializes the subserver (e.g., load configuration)
func (s *PhotoSaveSubServer) Init(ctx context.Context, opts ...option.Option) error {
	// Apply configuration options (e.g., set addresses from opts)
	for _, opt := range opts {
		s.opts.Apply(opt)
	}
	return nil
}

// Start begins listening for requests
func (s *PhotoSaveSubServer) Start(ctx context.Context, opts ...option.Option) error {
	s.status = server.StatusRunning
	return nil // Actual server start would be handled by the main server manager
}

// Stop shuts down the subserver
func (s *PhotoSaveSubServer) Stop(ctx context.Context) error {
	s.status = server.StatusStopped
	return nil
}

// Status returns current server status
func (s *PhotoSaveSubServer) Status() server.Status {
	return s.status
}

// handlePhotoUpload processes multipart file uploads and saves to photos-directory
func (s *PhotoSaveSubServer) handlePhotoUpload(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("photo")
	if err != nil {
		c.String(consts.StatusBadRequest, "Missing photo file: %v", err)
		return
	}

	// Create photos-directory if it doesn't exist
	uploadDir := "photos-directory"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.String(consts.StatusInternalServerError, "Failed to create upload directory: %v", err)
		return
	}
	// Construct full save path
	savePath := filepath.Join(uploadDir, file.Filename)

	// Save the uploaded file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.String(consts.StatusInternalServerError, "Failed to save photo: %v", err)
		return
	}

	c.String(consts.StatusOK, "Photo received: %s (size: %d bytes)", file.Filename, file.Size)
}

func NewPhotoSaveSubServer(opts ...option.Option) server.Subserver {
	s := &PhotoSaveSubServer{
		opts: option.GetOptions(opts...).Ctx().Value(serverOptionsKey{}).(*Options),
	}
	return s
}
