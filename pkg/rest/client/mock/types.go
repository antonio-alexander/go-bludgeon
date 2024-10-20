package mock

type Mock interface {
	//MockDoRequest
	MockDoRequest([]byte, int, error)
}
