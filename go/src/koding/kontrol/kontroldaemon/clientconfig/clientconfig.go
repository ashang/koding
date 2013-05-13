package clientconfig

import (
	"labix.org/v2/mgo/bson"
	"koding/tools/config"
	"labix.org/v2/mgo"
	"log"
)

type ServerInfo struct {
	BuildNumber string
	GitBranch   string
	ConfigUsed  string
	Config      *ConfigFile
	Hostname    Hostname
	IP          IP
}

type ConfigFile struct {
	Mongo string
	Mq    struct {
		Host          string
		Port          int
		ComponentUser string
		Password      string
		Vhost         string
	}
}

type Hostname struct {
	Public string
	Local  string
}

type IP struct {
	Public string
	Local  string
}

type ClientConfig struct {
	Hostname        string
	RegisteredHosts map[string][]string
	Session         *mgo.Session
	Collection      *mgo.Collection
}

func Connect() (*ClientConfig, error) {
	host := config.Current.Kontrold.Mongo.Host
	session, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Strong, true)
	col := session.DB("kontrol").C("clients")

	cc := &ClientConfig{
		Session:    session,
		Collection: col,
	}

	return cc, nil
}

func (c *ClientConfig) AddClient(info ServerInfo) {
	_, err := c.Collection.Upsert(bson.M{"buildnumber": info.BuildNumber}, info)
	if err != nil {
		log.Println(err)
	}

}
