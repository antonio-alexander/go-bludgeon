package bludgeonrestserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

//Router provides a struct that houses the REST server and pointers to communicate with
// the server
type restServer struct {
	sync.RWMutex                //mutex for threadsafe operations
	sync.WaitGroup              //waitgroup to manage goRoutines
	router         *mux.Router  //mux router
	server         *http.Server //rest Server
	log            *log.Logger  //logger
}

type Server interface {
	//
	Close()

	//
	BuildRoutes(routes []HandleFuncConfig) (err error)

	//
	UpdateLog(log *log.Logger)

	//
	Start(address, port string) (err error)

	//
	Stop() (err error)
}

//NewRouter creates a router pointer (from scratch) and requires a worker pool as an input
func NewServer() interface {
	Server
} {
	return &restServer{
		router: mux.NewRouter(),
	}
}

//Close will set all of the internal pointers to nil
func (r *restServer) Close() {
	r.Lock()
	defer r.Unlock()

	//set configuration to default
	//close internal pointers
	//set internal pointers to nil
	r.router, r.server = nil, nil
}

func (r *restServer) UpdateLog(log *log.Logger) {
	r.Lock()
	defer r.Unlock()

	//update the log
	r.log = log
}

//Start uses the configured mux/router to start listening to responses via REST
func (r *restServer) Start(address, port string) (err error) {
	r.Lock()
	defer r.Unlock()

	//create the routes to "answer" rest calls
	// r.buildRoutes()
	//create server pointer for REST
	r.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", address, port),
		Handler: r.router,
	}
	//start the goRoutine to run the server (listen and serve blocks)
	r.Add(1)
	go func() {
		defer r.Done()

		if err := r.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Printf("Httpserver: ListenAndServe() error: %s", err)
			}
		}
		//Do we need to account for a situation where the rest server kills itself
		// unexepctedly?
	}()

	return
}

//Stop will shutdown the rest server
func (r *restServer) Stop() (err error) {
	r.Lock()
	defer r.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), ConfigShutdownTimeout)
	defer cancel()

	r.server.Shutdown(ctx)
	r.Wait()

	return
}

func (r *restServer) BuildRoutes(routes []HandleFuncConfig) (err error) {
	r.Lock()
	defer r.Unlock()

	//
	for _, route := range routes {
		r.router.HandleFunc(route.Route, route.HandleFx).Methods(route.Method)
	}

	return
}
