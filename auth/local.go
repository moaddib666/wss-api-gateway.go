package auth

import "log"

type TestClient struct {
	id           string
	perms        []string
	Header       []byte
	isAdmin      bool
	isAuthorized bool
}

func (c *TestClient) IsAuthorised() bool {
	return c.isAuthorized
}

func (c *TestClient) IsAdmin() bool {
	return c.isAdmin
}

func (c *TestClient) Id() string {
	return c.id
}

func (c *TestClient) init(header []byte) {
	c.Header = header
	// TODO make api call to auth api to get client details.
	log.Printf("Constructing new client from %d header", len(c.Header))
	c.isAdmin = false
	c.isAuthorized = len(c.Header) > 0
	c.id = "TestClient"
}
