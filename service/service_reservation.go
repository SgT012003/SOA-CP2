package service

import (
	"errors"
	"fmt"
	"hotel-soa/dao"
	"hotel-soa/model"
	"net/http"
	"time"
)

type ReservationService interface {
	Create(res model.Reservation) (string, int, error)
	Update(res model.Reservation) (int, error)
	Delete(id string) error
	GetByID(id string) (model.Reservation, error)
	GetAll() ([]model.Reservation, error)
}

type reservationService struct{}

func NewReservationService() ReservationService {
	return &reservationService{}
}

// regras de transição de status válidas
var validTransitions = map[string][]string{
	"CREATED":    {"CHECKED_IN", "CANCELED"},
	"CHECKED_IN": {"CHECKED_OUT"},
}

// ---------------- CREATE ----------------
func (s *reservationService) Create(res model.Reservation) (string, int, error) {
	// 1. Validação de datas
	checkin, checkout, err := parseDates(res.CheckinExpected, res.CheckoutExpected)
	if err != nil {
		return "", http.StatusConflict, err
	}
	if !checkout.After(checkin) {
		return "", http.StatusConflict, errors.New("checkout_expected must be after checkin_expected")
	}

	// 2. Disponibilidade de quarto
	conflict, err := dao.HasReservationConflict(res.RoomID, checkin, checkout, "")
	if err != nil {
		return "", http.StatusConflict, err
	}
	if conflict {
		return "", http.StatusConflict, fmt.Errorf("room %s is not available for the selected dates", res.RoomID)
	}

	// 3. Status inicial
	if res.Status == "" {
		res.Status = "CREATED"
	}

	// 4. Persistência
	id, err := dao.InsertReservation(res)
	return id, http.StatusCreated, err
}

// ---------------- UPDATE ----------------
func (s *reservationService) Update(res model.Reservation) (int, error) {
	// 1. Buscar reserva atual
	current, err := dao.GetReservationByID(res.ID)
	if err != nil {
		return http.StatusNotFound, err
	}
	if current.ID == "" {
		return http.StatusNotFound, errors.New("reservation not found")
	}

	// 2. Validar fluxo de status
	if err := validateStatusTransition(current.Status, res.Status); err != nil {
		return http.StatusBadRequest, err
	}

	// 3. Validar datas se alteradas
	checkin, checkout, err := parseDates(res.CheckinExpected, res.CheckoutExpected)
	if err != nil {
		return http.StatusConflict, err
	}
	if !checkout.After(checkin) {
		return http.StatusConflict, errors.New("checkout_expected must be after checkin_expected")
	}

	// 4. Checar conflitos se mudou datas ou quarto
	if res.RoomID != current.RoomID ||
		res.CheckinExpected != current.CheckinExpected ||
		res.CheckoutExpected != current.CheckoutExpected {

		conflict, err := dao.HasReservationConflict(res.RoomID, checkin, checkout, res.ID)
		if err != nil {
			return http.StatusConflict, err
		}
		if conflict {
			return http.StatusConflict, fmt.Errorf("room %s is not available for the selected dates", res.RoomID)
		}
	}

	err = dao.UpdateReservation(res)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// ---------------- DELETE ----------------
func (s *reservationService) Delete(id string) error {
	return dao.DeleteReservation(id)
}

// ---------------- GET BY ID ----------------
func (s *reservationService) GetByID(id string) (model.Reservation, error) {
	return dao.GetReservationByID(id)
}

// ---------------- GET ALL ----------------
func (s *reservationService) GetAll() ([]model.Reservation, error) {
	return dao.GetAllReservations()
}

// ---------------- HELPERS ----------------

func parseDates(checkinStr, checkoutStr string) (time.Time, time.Time, error) {
	layout := "2006-01-02"
	checkin, err := time.Parse(layout, checkinStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid checkin_expected format (expected YYYY-MM-DD)")
	}
	checkout, err := time.Parse(layout, checkoutStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid checkout_expected format (expected YYYY-MM-DD)")
	}
	return checkin, checkout, nil
}

func validateStatusTransition(current, next string) error {
	if current == next {
		return nil
	}
	validNext, ok := validTransitions[current]
	if !ok {
		return fmt.Errorf("invalid current status: %s", current)
	}
	for _, allowed := range validNext {
		if next == allowed {
			return nil
		}
	}
	return fmt.Errorf("invalid status transition: %s → %s", current, next)
}
