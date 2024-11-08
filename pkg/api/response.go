package api

import (
	"log"

	"github.com/emicklei/go-restful/v3"
)

// WriteResponse is a helper function to write the response and log any errors
func WriteResponse(response *restful.Response, status int, entity interface{}) {
	var err error
	if entity != nil {
		err = response.WriteHeaderAndEntity(status, entity)
	} else {
		response.WriteHeader(status)
	}
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// WriteError is a helper function to write an error response and log any errors
func WriteError(response *restful.Response, status int, err error) {
	writeErr := response.WriteError(status, err)
	if writeErr != nil {
		log.Printf("Error writing error response: %v", writeErr)
	}
}
