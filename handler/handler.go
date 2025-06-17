package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/devdevaraj/bender/creator"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type BridgeRequest struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Subnet string `json:"subnet"`
	Driver string `json:"driver"`
}

type BridgeResponse struct {
	NID     string `json:"nid,omitempty"`
	CID     string `json:"cid,omitempty"`
	IP      string `json:"ip,omitempty"`
	Gateway string `json:"gateway,omitempty"`
	Subnet  string `json:"subnet,omitempty"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func CreateBridge(w http.ResponseWriter, r *http.Request, rdb *redis.Client, ctx context.Context) {
	var req BridgeRequest

	// vars := mux.Vars(r)
	// image := vars["image"]

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Image == "" {
		sendErrorResponse(w, "Image name is required", http.StatusBadRequest)
		return
	}

	if req.Id == "" {
		sendErrorResponse(w, "ID is required", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		sendErrorResponse(w, "Name and subnet are required", http.StatusBadRequest)
		return
	}

	nid, cid, ip, subnet, gateway, err := creator.CreateDockerBridge(req.Name, req.Image)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = rdb.Set(ctx, req.Id, ip, time.Minute).Err()
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, BridgeResponse{
		NID:     nid,
		CID:     cid,
		IP:      ip,
		Gateway: gateway,
		Subnet:  subnet,
		Message: "Bridge network created successfully",
	}, http.StatusCreated)
}

func DeleteBridge(w http.ResponseWriter, r *http.Request, rdb *redis.Client, ctx context.Context) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		sendErrorResponse(w, "Bridge name is required", http.StatusBadRequest)
		return
	}

	err := creator.DeleteDockerBridge(name)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = rdb.Del(ctx, name).Err()
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, BridgeResponse{
		Message: "Bridge network deleted successfully",
	}, http.StatusOK)
}

func ListBridges(w http.ResponseWriter, r *http.Request) {
	bridges, err := creator.ListDockerBridges()
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bridges)
}

func sendErrorResponse(w http.ResponseWriter, errMsg string, statusCode int) {
	sendResponse(w, BridgeResponse{
		Error: errMsg,
	}, statusCode)
}

func sendResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
