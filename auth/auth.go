package auth

import "log"

type Client struct {
	Id           string
	perms        []string
	Header       []byte
	isAdmin      bool
	isAuthorized bool
}

func (c *Client) IsAuthorised() bool {
	return c.isAuthorized
}

func (c *Client) IsAdmin() bool {
	return c.isAdmin
}

func (c *Client) init() {
	// TODO make api call to auth api to get client details.
	log.Printf("Constructing new client from %d header", len(c.Header))
	c.isAdmin = false
	c.isAuthorized = len(c.Header) > 0
	c.Id = "TestClient"
}

func NewClient(header []byte) *Client {
	client := &Client{Header: header}
	client.init()
	return client
}
