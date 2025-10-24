package model

import "fmt"

type Room struct {
	ID            string  `json:"id"`
	Number        int     `json:"number"`
	Type          string  `json:"type"`
	Capacity      int     `json:"capacity"`
	PricePerNight float64 `json:"price_per_night"`
	Status        string  `json:"status"`
}

type RoomRequest struct {
	ID            string  `json:"id"`
	Number        int     `json:"number" binding:"required"`
	Type          string  `json:"type" binding:"required"`
	Capacity      int     `json:"capacity" binding:"required"`
	PricePerNight float64 `json:"price_per_night" binding:"required,gt=0"`
	Status        string  `json:"status" binding:"required"`
}

func (r *RoomRequest) Room() *Room {
	return &Room{
		ID:            r.ID,
		Number:        r.Number,
		Type:          r.Type,
		Capacity:      r.Capacity,
		PricePerNight: r.PricePerNight,
		Status:        r.Status,
	}
}

func (r *Room) Validate() error {

	var errs []error
	if r.Number <= 0 {
		errs = append(errs, fmt.Errorf("invalid number must be greater than 0"))
	}
	if r.Capacity <= 0 {
		errs = append(errs, fmt.Errorf("invalid capacity must be greater than 0"))
	}
	if r.PricePerNight <= 0 {
		errs = append(errs, fmt.Errorf("invalid price must be greater than 0"))
	}

	switch r.Type {
	case "STANDARD", "DELUXE", "SUITE":
		break
	default:
		errs = append(errs, fmt.Errorf("invalid type field, must be one of: STANDARD, DELUXE, SUITE"))
	}

	switch r.Status {
	case "ATIVO", "INATIVO":
		break
	default:
		errs = append(errs, fmt.Errorf("invalid status field, must be one of: ATIVO, INATIVO"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors: %v", errs)
	}
	return nil
}
