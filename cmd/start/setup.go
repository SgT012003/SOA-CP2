package main

import (
	"fmt"
	"hotel-soa/db"
	"hotel-soa/model"
	"time"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("Starting setup...")
	createTables()
	seedTables()
	fmt.Println("Setup completed.")
}

func createTables() {
	fmt.Println("Creating tables...")
	createRoomTable()
	createReservationTable()
}

func createRoomTable() {
	fmt.Println("Creating room table...")
	query := `CREATE TABLE IF NOT EXISTS rooms (
		id CHAR(36) PRIMARY KEY,
		number INT NOT NULL UNIQUE,
		type VARCHAR(20) NOT NULL,
		capacity INT NOT NULL,
		price_per_night DECIMAL(10,2) NOT NULL,
		status VARCHAR(20) NOT NULL
	);`
	_, err := db.GetDB().Exec(query)
	if err != nil {
		fmt.Println("Error creating room table:", err)
	}
}

func createReservationTable() {
	fmt.Println("Creating reservation table...")
	query := `CREATE TABLE IF NOT EXISTS reservations (
		id CHAR(36) PRIMARY KEY,
		room_id CHAR(36) NOT NULL,
		guest_name VARCHAR(120) NOT NULL,
		checkin_expected DATE NOT NULL,
		checkout_expected DATE NOT NULL,
		status VARCHAR(20) NOT NULL,
		total_amount DECIMAL(10,2),
		CONSTRAINT fk_reservation_room FOREIGN KEY (room_id) REFERENCES rooms(id)
	);`
	_, err := db.GetDB().Exec(query)
	if err != nil {
		fmt.Println("Error creating reservation table:", err)
	}
}

func seedTables() {
	fmt.Println("Seeding tables...")
	seedRoomTable()
	seedReservationTable()
}

// Seeder de quartos
func seedRoomTable() {
	fmt.Println("Seeding room table...")

	rooms := []model.Room{
		{ID: uuid.NewString(), Number: 101, Type: "STANDARD", Capacity: 1, PricePerNight: 120.50, Status: "ATIVO"},
		{ID: uuid.NewString(), Number: 102, Type: "STANDARD", Capacity: 2, PricePerNight: 180.00, Status: "ATIVO"},
		{ID: uuid.NewString(), Number: 201, Type: "DELUXE", Capacity: 3, PricePerNight: 250.00, Status: "ATIVO"},
		{ID: uuid.NewString(), Number: 202, Type: "DELUXE", Capacity: 2, PricePerNight: 300.00, Status: "INATIVO"},
		{ID: uuid.NewString(), Number: 301, Type: "SUITE", Capacity: 4, PricePerNight: 500.00, Status: "ATIVO"},
	}

	for _, r := range rooms {
		_, err := db.GetDB().Exec(`
			INSERT INTO rooms (id, number, type, capacity, price_per_night, status)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (number) DO NOTHING;`,
			r.ID, r.Number, r.Type, r.Capacity, r.PricePerNight, r.Status)
		if err != nil {
			fmt.Println("Error seeding room:", err)
		}
	}

	fmt.Println("Rooms seeded successfully.")
}

// Seeder de reservas
func seedReservationTable() {
	fmt.Println("Seeding reservation table...")

	rows, err := db.GetDB().Query("SELECT id FROM rooms LIMIT 5;")
	if err != nil {
		fmt.Println("Error fetching rooms for reservation:", err)
		return
	}
	defer rows.Close()

	var roomIDs []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		roomIDs = append(roomIDs, id)
	}

	if len(roomIDs) == 0 {
		fmt.Println("No rooms found, skipping reservation seeding.")
		return
	}

	today := time.Now().Format("2006-01-02")
	twoDays := time.Now().AddDate(0, 0, 2).Format("2006-01-02")
	threeDays := time.Now().AddDate(0, 0, 3).Format("2006-01-02")

	reservations := []model.Reservation{
		{ID: uuid.NewString(), RoomID: roomIDs[0], GuestName: "Alice Silva", CheckinExpected: today, CheckoutExpected: twoDays, Status: "CREATED", TotalAmount: 240.00},
		{ID: uuid.NewString(), RoomID: roomIDs[1], GuestName: "Bruno Lima", CheckinExpected: today, CheckoutExpected: threeDays, Status: "CHECKED_IN", TotalAmount: 540.00},
		{ID: uuid.NewString(), RoomID: roomIDs[2], GuestName: "Carla Souza", CheckinExpected: today, CheckoutExpected: twoDays, Status: "CHECKED_OUT", TotalAmount: 500.00},
		{ID: uuid.NewString(), RoomID: roomIDs[3], GuestName: "Daniel Rocha", CheckinExpected: today, CheckoutExpected: threeDays, Status: "CREATED", TotalAmount: 900.00},
		{ID: uuid.NewString(), RoomID: roomIDs[4], GuestName: "Elisa Costa", CheckinExpected: today, CheckoutExpected: twoDays, Status: "CANCELED", TotalAmount: 0.00},
	}

	for _, r := range reservations {
		_, err := db.GetDB().Exec(`
			INSERT INTO reservations (id, room_id, guest_name, checkin_expected, checkout_expected, status, total_amount)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (id) DO NOTHING;`,
			r.ID, r.RoomID, r.GuestName, r.CheckinExpected, r.CheckoutExpected, r.Status, r.TotalAmount)
		if err != nil {
			fmt.Println("Error seeding reservation:", err)
		}
	}

	fmt.Println("Reservations seeded successfully.")
}
