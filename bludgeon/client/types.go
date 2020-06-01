package bludgeonclient

// "text/tabwriter"

//error constants
const (
	ErrStarted    string = "client started"
	ErrNotStarted string = "client not started"
)

//SerializedData
type SerializedData struct {
	//
}

type Configuration struct {
	ServerAddress string
	ServerPort    string
	ClientAddress string
	ClientPort    string
	Task          int64
	Employee      int64
}

//common constants
const (
	SQL_DRIVER = "sqlite3"
	HELP       = "Help Goes Here!"
	// tabwriterFlag = tabwriter.Debug //tabwriter.AlignRight | tabwriter.Debug
)
