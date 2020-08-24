package bludgeonconfig

import (
	"encoding/json"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client/config"
	metajson "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json/config"
	metamysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql/config"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/config"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server/config"
)

func DefaultClient(pwd string, c *Client) {
	//
	c.MetaType = DefaultMetaType
	c.Meta[bludgeon.MetaTypeJSON], _ = json.Marshal(metajson.Default())
	c.Meta[bludgeon.MetaTypeMySQL], _ = json.Marshal(metamysql.Default())
	c.RemoteType = DefaultRemoteType
	c.Remote[bludgeon.RemoteTypeRest], _ = json.Marshal(rest.Default())
	c.Rest = rest.Default()
	c.Client = client.Default()
}

func DefaultServer(pwd string, s *Server) {
	//
	s.MetaType = DefaultMetaType
	s.Meta[bludgeon.MetaTypeJSON], _ = json.Marshal(metajson.Default())
	s.Meta[bludgeon.MetaTypeMySQL], _ = json.Marshal(metamysql.Default())
	s.Rest = rest.Default()
	s.Server = server.Default()
}

//TODO: cli options to Server
//TODO: cli options to Client
