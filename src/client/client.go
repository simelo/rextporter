package client

type Client interface {
	GetRemoteInfo() ([]byte, error)
}
