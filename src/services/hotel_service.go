package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vandimit/simple-hotels-mock-rest-api/src/models"
)

// HotelService handles hotel data operations
type HotelService struct {
	Hotels []models.Hotel
}

// NewHotelService creates a new instance of HotelService
func NewHotelService() *HotelService {
	return &HotelService{
		Hotels: []models.Hotel{},
	}
}

// LoadHotelsFromFile loads hotel data from the specified JSON file
func (s *HotelService) LoadHotelsFromFile(filePath string) error {
	// Get the absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %w", err)
	}

	// Read file contents
	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Parse JSON into struct
	var response models.HotelResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	// Store hotels
	s.Hotels = response.Hotels
	return nil
}

// GetHotelByID returns a hotel by its ID
func (s *HotelService) GetHotelByID(id string) (*models.Hotel, error) {
	for _, hotel := range s.Hotels {
		if hotel.ID == id {
			return &hotel, nil
		}
	}
	return nil, errors.New("hotel not found")
}

// SearchHotels filters hotels based on search parameters
func (s *HotelService) SearchHotels(params models.SearchParams) []models.Hotel {
	var results []models.Hotel

	// Apply filters
	for _, hotel := range s.Hotels {
		// Skip if doesn't match name filter (case-insensitive, partial match)
		if params.Name != "" && !strings.Contains(strings.ToLower(hotel.Name), strings.ToLower(params.Name)) {
			continue
		}

		// Skip if doesn't match city filter (case-insensitive)
		if params.City != "" && !strings.EqualFold(hotel.City, params.City) {
			continue
		}

		// Skip if doesn't match country code filter
		if params.CountryCode != "" && hotel.CountryCode != params.CountryCode {
			continue
		}

		// Skip if below minimum rate
		if params.MinRate > 0 && hotel.LowRate < params.MinRate {
			continue
		}

		// Skip if above maximum rate
		if params.MaxRate > 0 && hotel.HighRate > params.MaxRate {
			continue
		}

		// Skip if below minimum rating
		if params.MinRating > 0 && hotel.HotelRating < params.MinRating {
			continue
		}

		// Skip if above maximum rating
		if params.MaxRating > 0 && hotel.HotelRating > params.MaxRating {
			continue
		}

		// Skip if doesn't match amenity mask (bitwise AND)
		if params.AmenityMask > 0 && (hotel.AmenityMask&params.AmenityMask) != params.AmenityMask {
			continue
		}

		// Add hotel to results if it passed all filters
		results = append(results, hotel)
	}

	// Apply pagination
	if params.Limit <= 0 {
		params.Limit = 20 // Default limit
	}

	if params.Offset < 0 {
		params.Offset = 0
	}

	// Handle pagination boundaries
	end := params.Offset + params.Limit
	if end > len(results) {
		end = len(results)
	}

	// Return paginated results
	if params.Offset < len(results) {
		return results[params.Offset:end]
	}

	// Return empty slice if offset is out of bounds
	return []models.Hotel{}
}