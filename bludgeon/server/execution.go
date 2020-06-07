package bludgeonserver

import (
	"fmt"
	"log"
	"net/http"

	bludgeon "github.com/antonio-alexander/go-bludgeon/bludgeon"
	bjson "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/json"
	bmysql "github.com/antonio-alexander/go-bludgeon/bludgeon/meta/sql/mysql"
)

func println(log *log.Logger, v ...interface{}) {
	if log != nil {
		log.Println(v...)
	}
}

func printf(log *log.Logger, format string, v ...interface{}) {
	if log != nil {
		log.Printf(format, v...)
	}
}

func print(log *log.Logger, v ...interface{}) {
	if log != nil {
		log.Print(v...)
	}
}

func getToken(request *http.Request) (token bludgeon.Token, err error) {
	//TODO: get token from request

	return
}

func handleResponse(writer http.ResponseWriter, errIn error, bytes []byte) (err error) {
	//check for errors, if so, write 500 internal server error
	if errIn != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(errIn.Error()))

		return
	}
	//if no error, write bytes
	_, err = writer.Write(bytes)

	return
}

func initMeta(config Configuration) (interface {
	bludgeon.MetaTimer
	bludgeon.MetaTimeSlice
}, func(), error) {
	// filepath.Join(pwd, "bludgeon.json")
	switch config.Meta.Type {
	case "json":
		//create metajson
		m := bjson.NewMetaJSON()
		//initialize metajson
		if err := m.Initialize(config.Meta.JSON.File); err != nil {
			return nil, nil, err
		}
		deferFx := func() {
			m.Close()
		}

		return m, deferFx, nil
	case "mysql":
		m := bmysql.NewMetaMySQL()
		//connect
		if err := m.Connect(config.Meta.MySQL); err != nil {
			return nil, nil, err
		}
		//create defer function
		deferFx := func() {
			//disconnect
			if err := m.Disconnect(); err != nil {
				fmt.Println(err)
			}
			//close
			m.Close()
		}

		return m, deferFx, nil
	default:
		return nil, func() {}, nil
	}
}
