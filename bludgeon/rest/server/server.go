package bludgeonrestserver

import (
	"context"
	"log"
	"net/http"
	"sync"

	client "github.com/antonio-alexander/go-bludgeon/bludgeon/client"
	rest "github.com/antonio-alexander/go-bludgeon/bludgeon/rest"
	server "github.com/antonio-alexander/go-bludgeon/bludgeon/server"

	"github.com/gorilla/mux"
)

//Router provides a struct that houses the REST server and pointers to communicate with
// the server
type restServer struct {
	sync.RWMutex                //mutex for threadsafe operations
	sync.WaitGroup              //waitgroup to manage goRoutines
	router         *mux.Router  //mux router
	server         *http.Server //rest Server
}

type Server interface {
	//
	Close()

	//
	BuildRoutes(i interface{}) (err error)

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

func (r *restServer) BuildRoutes(i interface{}) (err error) {
	r.Lock()
	defer r.Unlock()

	switch v := i.(type) {
	case client.Functional:
		r.buildRoutesClient(v)
	case server.Functional:
		r.buildRoutesServer(v)
	default:
		//TODO: generate error
	}

	return
}

//Start uses the configured mux/router to start listening to responses via REST
func (r *restServer) Start(address, port string) (err error) {
	r.Lock()
	defer r.Unlock()

	//create the routes to "answer" rest calls
	// r.buildRoutes()
	//create server pointer for REST
	r.server = &http.Server{
		Addr:    address + ":" + port,
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

//buildRoutes will create all the routes and their functions to execute when received
func (r *restServer) buildRoutesClient(client client.Functional) {
	// //Timer
	// r.router.HandleFunc(bludgeon.RouteTimerCreate, r.TimerCreate).Methods(RestPost)
	// r.router.HandleFunc(bludgeon.RouteTimerRead, r.TimerRead).Methods(RestPost)
	// r.router.HandleFunc(bludgeon.RouteTimersRead, r.TimersRead).Methods(RestPost)
	// r.router.HandleFunc(bludgeon.RouteTimerUpdate, r.TimerUpdate).Methods(RestPost)
	// r.router.HandleFunc(bludgeon.RouteTimerDelete, r.TimerDelete).Methods(RestPost)

	return
}

//buildRoutes will create all the routes and their functions to execute when received
func (r *restServer) buildRoutesServer(server server.Functional) {
	//admin
	//server
	//timer
	r.router.HandleFunc(rest.RouteTimerCreate, serverTimerCreate(server)).Methods(POST)
	r.router.HandleFunc(rest.RouteTimerRead, serverTimerRead(server)).Methods(POST)
	r.router.HandleFunc(rest.RouteTimerUpdate, serverTimerUpdate(server)).Methods(POST)
	r.router.HandleFunc(rest.RouteTimerDelete, serverTimerDelete(server)).Methods(POST)
	//time slice
	r.router.HandleFunc(rest.RouteTimeSliceCreate, serverTimeSliceCreate(server)).Methods(POST)
	r.router.HandleFunc(rest.RouteTimeSliceRead, serverTimeSliceRead(server)).Methods(POST)
	r.router.HandleFunc(rest.RouteTimeSliceUpdate, serverTimeSliceUpdate(server)).Methods(POST)
	r.router.HandleFunc(rest.RouteTimeSliceDelete, serverTimeSliceDelete(server)).Methods(POST)

	return
}
