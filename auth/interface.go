package auth

type Client interface {
	Id() string
	IsAuthorised() bool
	IsAdmin() bool
	init(header []byte)
}

func NewClient(header []byte) Client {
	client := &JWTClient{}
	client.init(header)
	return client
}
