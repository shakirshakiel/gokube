package server

import (
	"fmt"
	"log"
	"net/http"

	"etcdtest/pkg/api"
	"etcdtest/pkg/registry"
	"github.com/emicklei/go-restful/v3"
)

// APIServer represents the API server
type APIServer struct {
	nodeRegistry *registry.NodeRegistry
}

// NewAPIServer creates a new instance of APIServer
func NewAPIServer(nodeRegistry *registry.NodeRegistry) *APIServer {
	return &APIServer{
		nodeRegistry: nodeRegistry,
	}
}

// Start initializes and starts the API server
func (s *APIServer) Start(address string) error {
	container := restful.NewContainer()
	s.registerNodeRoutes(container)

	// TODO: Register routes for other types here

	return http.ListenAndServe(address, container)
}

// registerNodeRoutes adds Node-related routes to the container
func (s *APIServer) registerNodeRoutes(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/api/v1").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/nodes").To(s.createNode))
	ws.Route(ws.GET("/nodes/{name}").To(s.getNode))
	ws.Route(ws.PUT("/nodes/{name}").To(s.updateNode))
	ws.Route(ws.DELETE("/nodes/{name}").To(s.deleteNode))
	ws.Route(ws.GET("/nodes").To(s.listNodes))

	container.Add(ws)
}

// writeResponse is a helper function to write the response and log any errors
func writeResponse(response *restful.Response, status int, entity interface{}) {
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

// writeError is a helper function to write an error response and log any errors
func writeError(response *restful.Response, status int, err error) {
	writeErr := response.WriteError(status, err)
	if writeErr != nil {
		log.Printf("Error writing error response: %v", writeErr)
	}
}

// createNode handles POST requests to create a new Node
func (s *APIServer) createNode(request *restful.Request, response *restful.Response) {
	node := new(api.Node)
	err := request.ReadEntity(node)
	if err != nil {
		writeError(response, http.StatusBadRequest, err)
		return
	}

	err = s.nodeRegistry.CreateNode(request.Request.Context(), node)
	if err != nil {
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	writeResponse(response, http.StatusCreated, node)
}

// getNode handles GET requests to retrieve a Node
func (s *APIServer) getNode(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	node, err := s.nodeRegistry.GetNode(request.Request.Context(), name)
	if err != nil {
		writeError(response, http.StatusNotFound, err)
		return
	}

	writeResponse(response, http.StatusOK, node)
}

// updateNode handles PUT requests to update a Node
func (s *APIServer) updateNode(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	node := new(api.Node)
	err := request.ReadEntity(node)
	if err != nil {
		writeError(response, http.StatusBadRequest, err)
		return
	}

	if name != node.Name {
		writeError(response, http.StatusBadRequest,
			fmt.Errorf("Node name in URL does not match the name in the request body"))
		return
	}

	err = s.nodeRegistry.UpdateNode(request.Request.Context(), node)
	if err != nil {
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	writeResponse(response, http.StatusOK, node)
}

// deleteNode handles DELETE requests to remove a Node
func (s *APIServer) deleteNode(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	err := s.nodeRegistry.DeleteNode(request.Request.Context(), name)
	if err != nil {
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	writeResponse(response, http.StatusNoContent, nil)
}

// listNodes handles GET requests to list all Nodes
func (s *APIServer) listNodes(request *restful.Request, response *restful.Response) {
	nodes, err := s.nodeRegistry.ListNodes(request.Request.Context())
	if err != nil {
		writeError(response, http.StatusInternalServerError, err)
		return
	}

	writeResponse(response, http.StatusOK, nodes)
}
