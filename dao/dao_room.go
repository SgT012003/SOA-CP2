package dao

import (
	"database/sql"
	"hotel-soa/db"
	"hotel-soa/model"

	"github.com/google/uuid"
)

func InsertRoom(room model.Room) (string, error) {
	query := "INSERT INTO rooms (id, number, type, capacity, price_per_night, status) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (number) DO NOTHING;"
	_, err := db.GetDB().Exec(query, uuid.NewString(), room.Number, room.Type, room.Capacity, room.PricePerNight, room.Status)
	if err != nil {
		return "", err
	}
	return room.ID, nil
}

func UpdateRoom(room model.Room) error {
	query := "UPDATE rooms SET number = $1, type = $2, capacity = $3, price_per_night = $4, status = $5 WHERE id = $6;"
	_, err := db.GetDB().Exec(query, room.Number, room.Type, room.Capacity, room.PricePerNight, room.Status, room.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRoom(id string) error {
	query := "DELETE FROM rooms WHERE id = $1;"
	_, err := db.GetDB().Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func GetAllRooms() ([]model.Room, error) {
	var rooms []model.Room
	query := "SELECT id, number, type, capacity, price_per_night, status FROM rooms;"
	rows, err := db.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room model.Room
		if err := rows.Scan(&room.ID, &room.Number, &room.Type, &room.Capacity, &room.PricePerNight, &room.Status); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

func GetRoomByID(id string) (model.Room, error) {
	query := "SELECT id, number, type, capacity, price_per_night, status FROM rooms WHERE id = $1;"
	row := db.GetDB().QueryRow(query, id)
	var room model.Room
	if err := row.Scan(&room.ID, &room.Number, &room.Type, &room.Capacity, &room.PricePerNight, &room.Status); err != nil {
		if err == sql.ErrNoRows {
			return model.Room{}, nil
		}
		return model.Room{}, err
	}
	return room, nil
}
