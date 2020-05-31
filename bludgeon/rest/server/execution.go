package bludgeonrestserver

// import (
// 	"io/ioutil"
// 	"net/http"
// )

// //executeRequest is a function that encapsulates the functionality to respond to a
// // request that's been received.
// func executeClientRequest(fun) {
// 	//attempt to read the body of the request
// 	if bytesIn, err := ioutil.ReadAll(req.Body); err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 	} else {
// 		//close the body of the request
// 		if err := req.Body.Close(); err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(err.Error()))
// 		} else {
// 			//execute the command and return the response
// 			if bytesOut, err := r.executeCommand(CommandServerStop, bytesIn); err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				w.Write([]byte(err.Error()))
// 			} else {
// 				w.Write(bytesOut)
// 			}
// 		}
// 	}
// }
