package dao

import (
	"database/sql"
	"hotel-soa/db"
	"hotel-soa/model"
	"time"

	"github.com/google/uuid"
)

func InsertReservation(res model.Reservation) (string, error) {
	id := uuid.NewString()
	query := `INSERT INTO reservations 
		(id, room_id, guest_name, checkin_expected, checkout_expected, status, total_amount)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err := db.GetDB().Exec(query,
		id,
		res.RoomID,
		res.GuestName,
		res.CheckinExpected,
		res.CheckoutExpected,
		res.Status,
		res.TotalAmount,
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

func UpdateReservation(res model.Reservation) error {
	query := `UPDATE reservations 
		SET room_id = $1, guest_name = $2, checkin_expected = $3, 
		    checkout_expected = $4, status = $5, total_amount = $6
		WHERE id = $7;`
	_, err := db.GetDB().Exec(query,
		res.RoomID,
		res.GuestName,
		res.CheckinExpected,
		res.CheckoutExpected,
		res.Status,
		res.TotalAmount,
		res.ID,
	)
	return err
}

func DeleteReservation(id string) error {
	query := "DELETE FROM reservations WHERE id = $1;"
	_, err := db.GetDB().Exec(query, id)
	return err
}

func GetAllReservations() ([]model.Reservation, error) {
	var reservations []model.Reservation
	query := `SELECT id, room_id, guest_name, checkin_expected, 
		checkout_expected, status, total_amount FROM reservations;`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r model.Reservation
		if err := rows.Scan(
			&r.ID,
			&r.RoomID,
			&r.GuestName,
			&r.CheckinExpected,
			&r.CheckoutExpected,
			&r.Status,
			&r.TotalAmount,
		); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reservations, nil
}

func GetReservationByID(id string) (model.Reservation, error) {
	query := `SELECT id, room_id, guest_name, checkin_expected, 
		checkout_expected, status, total_amount 
		FROM reservations WHERE id = $1;`
	row := db.GetDB().QueryRow(query, id)

	var r model.Reservation
	if err := row.Scan(
		&r.ID,
		&r.RoomID,
		&r.GuestName,
		&r.CheckinExpected,
		&r.CheckoutExpected,
		&r.Status,
		&r.TotalAmount,
	); err != nil {
		if err == sql.ErrNoRows {
			return model.Reservation{}, nil
		}
		return model.Reservation{}, err
	}
	return r, nil
}

func HasReservationConflict(roomID string, checkin, checkout time.Time, excludeID string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM reservations 
		WHERE room_id = $1 
		  AND id != $2
		  AND status != 'CANCELED'
		  AND (
			(checkin_expected, checkout_expected) OVERLAPS ($3::date, $4::date)
		  );`
	var count int
	err := db.GetDB().QueryRow(query, roomID, excludeID, checkin, checkout).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
