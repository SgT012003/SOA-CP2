package model

import (
	"errors"
	"net/http"
)

type Reservation struct {
	ID               string  `json:"id"`
	RoomID           string  `json:"room_id"`
	GuestName        string  `json:"guest_name"`
	CheckinExpected  string  `json:"checkin_expected"`
	CheckoutExpected string  `json:"checkout_expected"`
	Status           string  `json:"status"`
	TotalAmount      float64 `json:"total_amount"`
}

type ReservationResponse struct {
	ID               string  `json:"id"`
	RoomID           string  `json:"room_id" binding:"required"`
	GuestName        string  `json:"guest_name" binding:"required"`
	CheckinExpected  string  `json:"checkin_expected" binding:"required"`
	CheckoutExpected string  `json:"checkout_expected" binding:"required"`
	Status           string  `json:"status" binding:"required"`
	TotalAmount      float64 `json:"total_amount" binding:"required,gt=0"`
}

func (r *ReservationResponse) Reservation() *Reservation {
	return &Reservation{
		ID:               r.ID,
		RoomID:           r.RoomID,
		GuestName:        r.GuestName,
		CheckinExpected:  r.CheckinExpected,
		CheckoutExpected: r.CheckoutExpected,
		Status:           r.Status,
		TotalAmount:      r.TotalAmount,
	}
}

func (r *Reservation) Validate() (int, error) {
	if r.RoomID == "" {
		return http.StatusBadRequest, errors.New("room_id is required")
	}
	if r.GuestName == "" {
		return http.StatusBadRequest, errors.New("guest_name is required")
	}
	if r.CheckinExpected == "" {
		return http.StatusBadRequest, errors.New("checkin_expected is required")
	}
	if r.CheckoutExpected == "" {
		return http.StatusBadRequest, errors.New("checkout_expected is required")
	}
	if r.Status == "" {
		return http.StatusBadRequest, errors.New("status is required")
	}
	if r.TotalAmount <= 0 {
		return http.StatusBadRequest, errors.New("total_amount must be greater than 0")
	}
	return http.StatusOK, nil
}
