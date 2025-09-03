package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vandimit/simple-hotels-mock-rest-api/src/models"
	"github.com/vandimit/simple-hotels-mock-rest-api/src/services"
)

// ReservationHandler handles HTTP requests for reservation data
type ReservationHandler struct {
	Service *services.ReservationService
}

// NewReservationHandler creates a new instance of ReservationHandler
func NewReservationHandler(service *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		Service: service,
	}
}

// GetReservations handles GET requests for all reservations of a hotel
func (h *ReservationHandler) GetReservations(w http.ResponseWriter, r *http.Request) {
	// Get hotelId from URL parameters
	vars := mux.Vars(r)
	hotelID := vars["hotelId"]

	// Get reservations from service
	reservations, err := h.Service.GetReservationsByHotelID(hotelID)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	// Return results
	sendJSONResponse(w, models.ReservationResponse{Reservations: reservations})
}

// GetReservationByID handles GET requests for a specific reservation
func (h *ReservationHandler) GetReservationByID(w http.ResponseWriter, r *http.Request) {
	// Get parameters from URL
	vars := mux.Vars(r)
	hotelID := vars["hotelId"]
	reservationID := vars["reservationId"]

	// Get the reservation
	reservation, err := h.Service.GetReservationByID(hotelID, reservationID)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	// Return the reservation
	sendJSONResponse(w, reservation)
}

// CreateReservation handles POST requests to create a new reservation
func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	// Get hotelId from URL parameters
	vars := mux.Vars(r)
	hotelID := vars["hotelId"]

	// Parse request body
	var req models.CreateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.CustomerName == "" || req.StartDate == "" || req.EndDate == "" {
		sendErrorResponse(w, http.StatusBadRequest, "CustomerName, StartDate, and EndDate are required fields")
		return
	}

	// Create the reservation
	reservation, err := h.Service.CreateReservation(hotelID, req)
	if err != nil {
		if err.Error() == "hotel not found" {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "reservation dates overlap with an existing booking" {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	// Return the created reservation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reservation)
}

// UpdateReservation handles PUT requests to update an existing reservation
func (h *ReservationHandler) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	// Get parameters from URL
	vars := mux.Vars(r)
	hotelID := vars["hotelId"]
	reservationID := vars["reservationId"]

	// Parse request body
	var req models.UpdateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.CustomerName == "" || req.StartDate == "" || req.EndDate == "" {
		sendErrorResponse(w, http.StatusBadRequest, "CustomerName, StartDate, and EndDate are required fields")
		return
	}

	// Update the reservation
	reservation, err := h.Service.UpdateReservation(hotelID, reservationID, req)
	if err != nil {
		if err.Error() == "hotel not found" || err.Error() == "reservation not found" {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "reservation dates overlap with an existing booking" {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	// Return the updated reservation
	sendJSONResponse(w, reservation)
}

// DeleteReservation handles DELETE requests to remove a reservation
func (h *ReservationHandler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	// Get parameters from URL
	vars := mux.Vars(r)
	hotelID := vars["hotelId"]
	reservationID := vars["reservationId"]

	// Delete the reservation
	err := h.Service.DeleteReservation(hotelID, reservationID)
	if err != nil {
		if err.Error() == "hotel not found" || err.Error() == "reservation not found" {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete reservation")
		}
		return
	}

	// Return success with no content
	w.WriteHeader(http.StatusNoContent)
}
