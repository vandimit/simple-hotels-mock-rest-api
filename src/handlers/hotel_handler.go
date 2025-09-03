package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vandimit/simple-hotels-mock-rest-api/src/models"
	"github.com/vandimit/simple-hotels-mock-rest-api/src/services"
)

// HotelHandler handles HTTP requests for hotel data
type HotelHandler struct {
	Service *services.HotelService
}

// NewHotelHandler creates a new instance of HotelHandler
func NewHotelHandler(service *services.HotelService) *HotelHandler {
	return &HotelHandler{
		Service: service,
	}
}

// GetHotels handles GET requests for searching hotels
func (h *HotelHandler) GetHotels(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for search filters
	params := parseSearchParams(r)

	// Search hotels based on parameters
	hotels := h.Service.SearchHotels(params)

	// Return results
	sendJSONResponse(w, models.HotelResponse{Hotels: hotels})
}

// GetHotelByID handles GET requests for a specific hotel by ID
func (h *HotelHandler) GetHotelByID(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL parameters
	vars := mux.Vars(r)
	id := vars["hotelId"]

	// Find hotel by ID
	hotel, err := h.Service.GetHotelByID(id)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "Hotel not found")
		return
	}

	// Return the hotel
	sendJSONResponse(w, hotel)
}

// parseSearchParams extracts search parameters from the HTTP request
func parseSearchParams(r *http.Request) models.SearchParams {
	query := r.URL.Query()
	params := models.SearchParams{
		Name:        query.Get("name"),
		City:        query.Get("city"),
		CountryCode: query.Get("countryCode"),
	}

	// Parse numeric parameters with proper error handling
	if minRate := query.Get("minRate"); minRate != "" {
		if val, err := strconv.ParseFloat(minRate, 64); err == nil {
			params.MinRate = val
		}
	}

	if maxRate := query.Get("maxRate"); maxRate != "" {
		if val, err := strconv.ParseFloat(maxRate, 64); err == nil {
			params.MaxRate = val
		}
	}

	if minRating := query.Get("minRating"); minRating != "" {
		if val, err := strconv.ParseFloat(minRating, 64); err == nil {
			params.MinRating = val
		}
	}

	if maxRating := query.Get("maxRating"); maxRating != "" {
		if val, err := strconv.ParseFloat(maxRating, 64); err == nil {
			params.MaxRating = val
		}
	}

	if amenityMask := query.Get("amenityMask"); amenityMask != "" {
		if val, err := strconv.Atoi(amenityMask); err == nil {
			params.AmenityMask = val
		}
	}

	// Parse pagination parameters
	limit := 20 // Default limit
	if limitParam := query.Get("limit"); limitParam != "" {
		if val, err := strconv.Atoi(limitParam); err == nil && val > 0 {
			limit = val
		}
	}
	params.Limit = limit

	offset := 0 // Default offset
	if offsetParam := query.Get("offset"); offsetParam != "" {
		if val, err := strconv.Atoi(offsetParam); err == nil && val >= 0 {
			offset = val
		}
	}
	params.Offset = offset

	return params
}

// sendJSONResponse sends a JSON response with the provided data
func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends a JSON error response with the provided status code and message
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Code:    statusCode,
		Message: message,
	})
}
