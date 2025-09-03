package models

// Hotel represents a hotel entity with all its properties
type Hotel struct {
	ID                   string   `json:"id"`
	Type                 string   `json:"type"`
	Name                 string   `json:"name"`
	Created              int64    `json:"created"`
	Modified             int64    `json:"modified"`
	Address1             string   `json:"address1"`
	AirportCode          string   `json:"airportCode"`
	AmenityMask          int      `json:"amenityMask"`
	City                 string   `json:"city"`
	ConfidenceRating     int      `json:"confidenceRating"`
	CountryCode          string   `json:"countryCode"`
	DeepLink             string   `json:"deepLink"`
	HighRate             float64  `json:"highRate"`
	HotelID              int      `json:"hotelId"`
	HotelInDestination   bool     `json:"hotelInDestination"`
	HotelRating          float64  `json:"hotelRating"`
	Location             Location `json:"location"`
	LocationDescription  string   `json:"locationDescription"`
	LowRate              float64  `json:"lowRate"`
	Metadata             Metadata `json:"metadata"`
	PostalCode           string   `json:"postalCode"`
	PropertyCategory     int      `json:"propertyCategory"`
	ProximityDistance    float64  `json:"proximityDistance"`
	ProximityUnit        string   `json:"proximityUnit"`
	RateCurrencyCode     string   `json:"rateCurrencyCode"`
	ShortDescription     string   `json:"shortDescription"`
	StateProvinceCode    string   `json:"stateProvinceCode"`
	ThumbNailUrl         string   `json:"thumbNailUrl"`
	TripAdvisorRating    float64  `json:"tripAdvisorRating"`
	TripAdvisorRatingUrl string   `json:"tripAdvisorRatingUrl"`
}

// Location represents geographic coordinates
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Metadata contains additional hotel information
type Metadata struct {
	Path string `json:"path"`
}

// HotelResponse represents the response format for hotel data
type HotelResponse struct {
	Hotels []Hotel `json:"hotels"`
}

// ErrorResponse represents the error response format
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SearchParams represents the search parameters for filtering hotels
type SearchParams struct {
	Name        string  `json:"name"`
	City        string  `json:"city"`
	CountryCode string  `json:"countryCode"`
	MinRate     float64 `json:"minRate"`
	MaxRate     float64 `json:"maxRate"`
	MinRating   float64 `json:"minRating"`
	MaxRating   float64 `json:"maxRating"`
	AmenityMask int     `json:"amenityMask"`
	Limit       int     `json:"limit"`
	Offset      int     `json:"offset"`
}
