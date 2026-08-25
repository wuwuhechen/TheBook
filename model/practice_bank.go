package model

import "fmt"

type PracticeBank struct {
	practices map[int]*Practice
}

func NewPracticeBank() *PracticeBank {
	return &PracticeBank{
		practices: make(map[int]*Practice),
	}
}

func (pb *PracticeBank) Create(practice *Practice) error {
	if _, exists := pb.practices[practice.ID]; exists {
		return fmt.Errorf("practice with ID %d already exists", practice.ID)
	}
	pb.practices[practice.ID] = practice
	return nil
}

func (pb *PracticeBank) FindByID(id int) (*Practice, error) {
	practice, exists := pb.practices[id]
	if !exists {
		return nil, fmt.Errorf("practice with ID %d not found", id)
	}
	return practice, nil
}

func (pb *PracticeBank) Save(practice *Practice) error {
	if _, exists := pb.practices[practice.ID]; !exists {
		return fmt.Errorf("practice with ID %d does not exist", practice.ID)
	}
	pb.practices[practice.ID] = practice
	return nil
}

func (pb *PracticeBank) Delete(practice *Practice) error {
	if _, exists := pb.practices[practice.ID]; !exists {
		return fmt.Errorf("practice with ID %d does not exist", practice.ID)
	}
	delete(pb.practices, practice.ID)
	return nil
}
