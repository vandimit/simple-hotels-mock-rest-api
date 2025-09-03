package models

import (
	"time"
)

// Reservation represents a hotel reservation
type Reservation struct {
	ID           string    `json:"id"`
	HotelID      string    `json:"hotelId"`
	CustomerName string    `json:"customerName"`
	StartDate    string    `json:"startDate"` // ISO 8601 format: YYYY-MM-DD
	EndDate      string    `json:"endDate"`   // ISO 8601 format: YYYY-MM-DD
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ReservationResponse represents the response format for reservation data
type ReservationResponse struct {
	Reservations []Reservation `json:"reservations"`
}

// CreateReservationRequest represents the request body for creating a reservation
type CreateReservationRequest struct {
	CustomerName string `json:"customerName"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
}

// UpdateReservationRequest represents the request body for updating a reservation
type UpdateReservationRequest struct {
	CustomerName string `json:"customerName"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
}

// ParseDate parses a date string in YYYY-MM-DD format
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// IsOverlapping checks if two date ranges overlap
func IsOverlapping(start1, end1, start2, end2 time.Time) bool {
	// Two date ranges overlap if one starts before the other ends and ends after the other starts
	return (start1.Before(end2) || start1.Equal(end2)) && (end1.After(start2) || end1.Equal(start2))
}
