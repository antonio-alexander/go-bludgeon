package bludgeon

//--------------------------------------------------------------------------------------------
// utilities.go contains types and functions that are core functions to be used in bludgeon
//--------------------------------------------------------------------------------------------

import (
	"log"
	"sync"
	"time"
)

//--------------------------------------------------------------------------------------------
// Queue
//--------------------------------------------------------------------------------------------

var (
	empty struct{} //empty provides a variable that can be used to signal easily
)

//Queue is a basic queue data structure that allows you to buffer finite data of a configured size
type Queue struct {
	sync.RWMutex
	data    []interface{}
	signal  chan struct{}
	head    int
	foot    int
	maxSize int
}

//newQueue creates the internal pointers for a Queue
func newQueue(maxSize int) (data []interface{}, signal chan struct{}) {
	data = make([]interface{}, maxSize)
	signal = make(chan struct{})

	return
}

//NewQueue creates a queue of the configured size and returns a pointer
func NewQueue(maxSize int) *Queue {
	data, signal := newQueue(maxSize)

	return &Queue{
		maxSize: maxSize,
		head:    0,
		foot:    0,
		data:    data,
		signal:  signal}
}

//Close will set all internal pointers to nil
func (q *Queue) Close() {
	q.Lock()
	q.Unlock()

	q.data = nil
	q.head, q.foot = 0, 0
	close(q.signal)
}

//New re-purposes a queue pointer with the configured size
func (q *Queue) New(maxSize int) {
	q.Lock()
	defer q.Unlock()

	q.data = nil
	q.data, q.signal = newQueue(maxSize)
}

//GetSignal can get the signal used when data has been placed into the queue
func (q *Queue) GetSignal() (signal chan struct{}) {
	q.RLock()
	defer q.RUnlock()

	signal = q.signal
	return
}

//Enqueue puts an element into the queue and returns an overflow if the queue is full
func (q *Queue) Enqueue(element interface{}) (overflow bool) {
	q.Lock()
	defer q.Unlock()
	if q.head-q.foot == q.maxSize {
		overflow = true
		return
	}
	q.rotate()
	q.data[0] = element
	q.head++
	q.signal <- empty

	return
}

//Dequeue returns one element from the queue
func (q *Queue) Dequeue() (element interface{}, underflow bool) {
	q.Lock()
	defer q.Unlock()

	if q.head-q.foot == 0 {
		underflow = true
		return
	}
	element = q.data[q.head-1]
	q.head--

	return
}

//Length returns the number of elements in the queue
func (q *Queue) Length() (length int) {
	q.RLock()
	defer q.RUnlock()

	length = q.head - q.foot

	return
}

//rotate will rotate the queue by the
func (q *Queue) rotate() {
	l := len(q.data)
	n := 1

	q.data = append(q.data[l-n:], q.data[:l-n]...)
}

//--------------------------------------------------------------------------------------------
// Pool
//--------------------------------------------------------------------------------------------

//WorkerFx is the function type used for worker functions (executed for each worker within a pool)
type WorkerFx func(id int, input interface{})

//Pool provides a type to store data pertinent to the worker pool
type Pool struct {
	sync.RWMutex
	wg         sync.WaitGroup
	workerFx   WorkerFx
	nWorkers   int
	inputQueue *Queue
	stopper    chan struct{}
	log        *Log
}

//newPool is a private function that creates all the internal pointers for a pool
func newPool() (stopper chan (struct{}), log *Log) {
	stopper = make(chan struct{})
	log = NewLog()

	return
}

//NewPool will return a pointer of type Pool
func NewPool() *Pool {
	stopper, log := newPool()

	return &Pool{
		stopper: stopper,
		log:     log}
}

//Close will set all internal pointers to nil
func (p *Pool) Close() {
	p.Lock()
	defer p.Unlock()

	p.inputQueue.Close()
	p.log.Close()

	p.stopper, p.inputQueue, p.log = nil, nil, nil
}

//Start will create all the workers and queues and will block until the workers are within their
// business logic
func (p *Pool) Start(nWorkers int, workerFx WorkerFx) {
	p.Lock()
	defer p.Unlock()

	p.nWorkers = nWorkers
	p.workerFx = workerFx
	p.inputQueue = NewQueue(p.nWorkers)

	for i := 0; i < p.nWorkers; i++ {
		started := make(chan struct{})
		p.wg.Add(1)
		go p.goWorker(i, started)
		<-started
	}
}

//goWorker is the goRoutine that will handle the function of the worker pool as data is input to it
func (p *Pool) goWorker(workerID int, started chan struct{}) {
	defer p.wg.Done()

	signal := p.inputQueue.GetSignal()
	close(started)

	for {
		select {
		case <-p.stopper:
			return
		case <-signal:
			if input, underflow := p.inputQueue.Dequeue(); !underflow {
				p.workerFx(workerID, input)
			}
		}
	}
}

//Input will place data onto the internal queue
func (p *Pool) Input(element interface{}) (overflow bool) {
	overflow = p.inputQueue.Enqueue(element)

	return
}

//Stop will stop the pool
func (p *Pool) Stop() {
	p.Lock()
	defer p.Unlock()

	close(p.stopper)
	p.wg.Wait()
}

//--------------------------------------------------------------------------------------------
// Log
//--------------------------------------------------------------------------------------------

const (
	//DefaultLogAlias provides a default alias when none is given
	DefaultLogAlias = "<No alias>"
	//DefaultVerbosity provides a default verbosity when none is given
	DefaultVerbosity = 0
)

//Log is a type used to wrap the log structure and provide alias/enable storage as well as future
// functionality
type Log struct {
	log       *log.Logger
	alias     string
	verbosity int
}

//NewLog will return a pointer of type Log that can be used to log any errors or messages
func NewLog() *Log {
	return &Log{
		alias:     DefaultLogAlias,
		verbosity: DefaultVerbosity}
}

//Close will set any internal pointers to nil
func (l *Log) Close() {
	l.log = nil
}

//UpdateLog allows setting of the interal log pointer
func (l *Log) UpdateLog(log *log.Logger) {
	l.log = log
}

//Configure allows setting of the internally stored alias and verbosity
func (l *Log) Configure(alias string, verbosity int) {
	l.alias, l.verbosity = alias, verbosity
}

//Print will print the info string with the provided alias if the alias is enabled along with a
// timestamp
func (l *Log) Print(verbosity int, info string) {
	if l.log != nil {
		if verbosity >= l.verbosity {
			l.log.Print(l.alias + " (" + time.Now().String() + "): " + info)
		}
	}
}

//Printf will print the format/interface string with the provided alias if the alias is enabled
// along with a timestamp
func (l *Log) Printf(verbosity int, format string, v ...interface{}) {
	if l.log != nil {
		if verbosity >= l.verbosity {
			l.log.Printf(l.alias+" ("+time.Now().String()+"): "+format, v...)
		}
	}
}
