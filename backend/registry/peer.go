package registry

type Peer interface {
	Name() string
	// Metadata TODO: just map for now
	Metadata() map[string]string
	Addresses() []string
}
