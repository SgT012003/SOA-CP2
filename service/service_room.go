package service

import (
	"hotel-soa/dao"
	"hotel-soa/model"
)

type RoomService interface {
	Create(room model.Room) (string, error)
	Update(room model.Room) error
	Delete(id string) error
	GetByID(id string) (model.Room, error)
	GetAll() ([]model.Room, error)
}

type roomService struct{}

func NewRoomService() RoomService {
	return &roomService{}
}

func (s *roomService) Create(room model.Room) (string, error) {
	return dao.InsertRoom(room)
}

func (s *roomService) Update(room model.Room) error {
	return dao.UpdateRoom(room)
}

func (s *roomService) Delete(id string) error {
	return dao.DeleteRoom(id)
}

func (s *roomService) GetByID(id string) (model.Room, error) {
	return dao.GetRoomByID(id)
}

func (s *roomService) GetAll() ([]model.Room, error) {
	return dao.GetAllRooms()
}
