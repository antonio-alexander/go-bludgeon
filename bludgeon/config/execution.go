package bludgeonconfig

import (
	"encoding/json"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client/config"
	metajson "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json/config"

	// metamysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/mysql"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest/config"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server/config"
)

func DefaultClient(directory string, c *Client) {
	//
	c.MetaType = DefaultMetaType
	if c.Meta == nil {
		c.Meta = make(map[bludgeon.MetaType]json.RawMessage)
	}
	c.Meta[bludgeon.MetaTypeJSON], _ = json.Marshal(metajson.Default(directory))
	// c.Meta[bludgeon.MetaTypeMySQL], _ = json.Marshal(metamysql.Default())
	c.RemoteType = DefaultRemoteType
	if c.Remote == nil {
		c.Remote = make(map[bludgeon.RemoteType]json.RawMessage)
	}
	c.Remote[bludgeon.RemoteTypeRest], _ = json.Marshal(rest.Default())
	c.Rest = rest.Default()
	c.Client = client.Default()
}

func DefaultServer(directory string, s *Server) {
	//
	s.MetaType = DefaultMetaType
	if s.Meta == nil {
		s.Meta = make(map[bludgeon.MetaType]json.RawMessage)
	}
	s.Meta[bludgeon.MetaTypeJSON], _ = json.Marshal(metajson.Default(directory))
	// s.Meta[bludgeon.MetaTypeMySQL], _ = json.Marshal(metamysql.Default())
	s.Rest = rest.Default()
	s.Server = server.Default()
}

//TODO: cli options to Server
//TODO: cli options to Client
