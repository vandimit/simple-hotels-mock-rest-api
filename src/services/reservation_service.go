package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vandimit/simple-hotels-mock-rest-api/src/models"
)

// ReservationService handles reservation operations
type ReservationService struct {
	hotelService *HotelService
	reservations map[string][]models.Reservation // map[hotelID][]Reservation
	mutex        sync.RWMutex
}

// NewReservationService creates a new instance of ReservationService
func NewReservationService(hotelService *HotelService) *ReservationService {
	return &ReservationService{
		hotelService: hotelService,
		reservations: make(map[string][]models.Reservation),
		mutex:        sync.RWMutex{},
	}
}

// GetReservationsByHotelID returns all reservations for a hotel
func (s *ReservationService) GetReservationsByHotelID(hotelID string) ([]models.Reservation, error) {
	// Check if the hotel exists
	if _, err := s.hotelService.GetHotelByID(hotelID); err != nil {
		return nil, errors.New("hotel not found")
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Return empty slice if no reservations for this hotel
	if reservations, ok := s.reservations[hotelID]; ok {
		return reservations, nil
	}

	return []models.Reservation{}, nil
}

// GetReservationByID returns a reservation by its ID
func (s *ReservationService) GetReservationByID(hotelID, reservationID string) (*models.Reservation, error) {
	// Check if the hotel exists
	if _, err := s.hotelService.GetHotelByID(hotelID); err != nil {
		return nil, errors.New("hotel not found")
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Find the reservation
	if reservations, ok := s.reservations[hotelID]; ok {
		for _, reservation := range reservations {
			if reservation.ID == reservationID {
				return &reservation, nil
			}
		}
	}

	return nil, errors.New("reservation not found")
}

// CreateReservation creates a new reservation for a hotel
func (s *ReservationService) CreateReservation(hotelID string, req models.CreateReservationRequest) (*models.Reservation, error) {
	// Check if the hotel exists
	if _, err := s.hotelService.GetHotelByID(hotelID); err != nil {
		return nil, errors.New("hotel not found")
	}

	// Validate dates
	startDate, err := models.ParseDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := models.ParseDate(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Ensure end date is after start date
	if !endDate.After(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	// Check for overlapping reservations
	if s.hasOverlappingReservations(hotelID, startDate, endDate, "") {
		return nil, errors.New("reservation dates overlap with an existing booking")
	}

	// Create the reservation
	now := time.Now().UTC()
	reservation := models.Reservation{
		ID:           uuid.New().String(),
		HotelID:      hotelID,
		CustomerName: req.CustomerName,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Store the reservation
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.reservations[hotelID]; !ok {
		s.reservations[hotelID] = []models.Reservation{}
	}
	s.reservations[hotelID] = append(s.reservations[hotelID], reservation)

	return &reservation, nil
}

// UpdateReservation updates an existing reservation
func (s *ReservationService) UpdateReservation(hotelID, reservationID string, req models.UpdateReservationRequest) (*models.Reservation, error) {
	// Check if the hotel exists
	if _, err := s.hotelService.GetHotelByID(hotelID); err != nil {
		return nil, errors.New("hotel not found")
	}

	// Validate dates
	startDate, err := models.ParseDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := models.ParseDate(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Ensure end date is after start date
	if !endDate.After(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	// Check for overlapping reservations (excluding the current reservation)
	if s.hasOverlappingReservations(hotelID, startDate, endDate, reservationID) {
		return nil, errors.New("reservation dates overlap with an existing booking")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Find and update the reservation
	if reservations, ok := s.reservations[hotelID]; ok {
		for i, reservation := range reservations {
			if reservation.ID == reservationID {
				// Update reservation
				reservations[i].CustomerName = req.CustomerName
				reservations[i].StartDate = req.StartDate
				reservations[i].EndDate = req.EndDate
				reservations[i].UpdatedAt = time.Now().UTC()
				return &reservations[i], nil
			}
		}
	}

	return nil, errors.New("reservation not found")
}

// DeleteReservation deletes a reservation
func (s *ReservationService) DeleteReservation(hotelID, reservationID string) error {
	// Check if the hotel exists
	if _, err := s.hotelService.GetHotelByID(hotelID); err != nil {
		return errors.New("hotel not found")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Find and delete the reservation
	if reservations, ok := s.reservations[hotelID]; ok {
		for i, reservation := range reservations {
			if reservation.ID == reservationID {
				// Remove reservation from slice
				s.reservations[hotelID] = append(reservations[:i], reservations[i+1:]...)
				return nil
			}
		}
	}

	return errors.New("reservation not found")
}

// hasOverlappingReservations checks if a date range overlaps with any existing reservations
// excludeReservationID is an optional parameter to exclude a specific reservation from the check (used during updates)
func (s *ReservationService) hasOverlappingReservations(hotelID string, startDate, endDate time.Time, excludeReservationID string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if reservations, ok := s.reservations[hotelID]; ok {
		for _, reservation := range reservations {
			// Skip the excluded reservation
			if reservation.ID == excludeReservationID {
				continue
			}

			// Parse existing reservation dates
			existingStartDate, err := models.ParseDate(reservation.StartDate)
			if err != nil {
				continue
			}

			existingEndDate, err := models.ParseDate(reservation.EndDate)
			if err != nil {
				continue
			}

			// Check for overlap
			if models.IsOverlapping(startDate, endDate, existingStartDate, existingEndDate) {
				return true
			}
		}
	}

	return false
}