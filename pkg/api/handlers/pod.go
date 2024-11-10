package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"

	"gokube/pkg/api"
	"gokube/pkg/registry"
)

// PodHandler handles Pod-related requests
type PodHandler struct {
	podRegistry *registry.PodRegistry
}

// NewPodHandler creates a new instance of PodHandler
func NewPodHandler(podRegistry *registry.PodRegistry) *PodHandler {
	return &PodHandler{podRegistry: podRegistry}
}

// CreatePod handles POST requests to create a new Pod
func (h *PodHandler) CreatePod(request *restful.Request, response *restful.Response) {
	pod := new(api.Pod)
	if err := request.ReadEntity(pod); err != nil {
		api.WriteError(response, http.StatusBadRequest, err)
		return
	}

	if err := h.podRegistry.CreatePod(request.Request.Context(), pod); err != nil {
		switch {
		case errors.Is(err, registry.ErrPodAlreadyExists):
			api.WriteError(response, http.StatusConflict, err)
		case errors.Is(err, registry.ErrPodInvalid):
			api.WriteError(response, http.StatusBadRequest, err)
			return
		default:
			api.WriteError(response, http.StatusInternalServerError, err)
			return
		}
	}

	api.WriteResponse(response, http.StatusCreated, pod)
}

// ListPods handles GET requests to list all Pods
func (h *PodHandler) ListPods(request *restful.Request, response *restful.Response) {
	pods, err := h.podRegistry.ListPods(request.Request.Context())
	if err != nil {
		api.WriteError(response, http.StatusInternalServerError, err)
		return
	}

	api.WriteResponse(response, http.StatusOK, pods)
}

// GetPod handles GET requests to retrieve a Pod
func (h *PodHandler) GetPod(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	pod, err := h.podRegistry.GetPod(request.Request.Context(), name)
	if err != nil {
		switch {
		case errors.Is(err, registry.ErrPodNotFound):
			api.WriteError(response, http.StatusNotFound, err)
		default:
			api.WriteError(response, http.StatusInternalServerError, err)
		}
		return
	}

	api.WriteResponse(response, http.StatusOK, pod)
}

// UpdatePod handles PUT requests to update a Pod
func (h *PodHandler) UpdatePod(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	pod := new(api.Pod)
	if err := request.ReadEntity(pod); err != nil {
		api.WriteError(response, http.StatusBadRequest, err)
		return
	}

	if name != pod.Name {
		api.WriteError(response, http.StatusBadRequest, fmt.Errorf("pod name in URL does not match pod name in request body"))
		return
	}

	if err := h.podRegistry.UpdatePod(request.Request.Context(), pod); err != nil {
		switch {
		case errors.Is(err, registry.ErrPodInvalid):
			api.WriteError(response, http.StatusBadRequest, err)
			return
		default:
			api.WriteError(response, http.StatusInternalServerError, err)
			return
		}
	}

	api.WriteResponse(response, http.StatusOK, pod)
}

// DeletePod handles DELETE requests to remove a Pod
func (h *PodHandler) DeletePod(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("name")
	if err := h.podRegistry.DeletePod(request.Request.Context(), name); err != nil {
		api.WriteError(response, http.StatusInternalServerError, err)
		return
	}

	api.WriteResponse(response, http.StatusNoContent, nil)
}

// ListUnassignedPods handles GET requests to list all unassigned Pods
func (h *PodHandler) ListUnassignedPods(request *restful.Request, response *restful.Response) {
	pods, err := h.podRegistry.ListUnassignedPods(request.Request.Context())
	if err != nil {
		api.WriteError(response, http.StatusInternalServerError, err)
		return
	}

	api.WriteResponse(response, http.StatusOK, pods)
}

func RegisterPodRoutes(ws *restful.WebService, podHandler *PodHandler) {
	ws.Route(ws.POST("/pods").To(podHandler.CreatePod))
	ws.Route(ws.GET("/pods").To(podHandler.ListPods))
	ws.Route(ws.GET("/pods/{name}").To(podHandler.GetPod))
	ws.Route(ws.PUT("/pods/{name}").To(podHandler.UpdatePod))
	ws.Route(ws.DELETE("/pods/{name}").To(podHandler.DeletePod))
	ws.Route(ws.GET("/pods/unassigned").To(podHandler.ListUnassignedPods))
}
