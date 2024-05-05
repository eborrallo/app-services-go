package increasing

import (
	"errors"
	"log"
)

type CourseCounterService struct{}

func NewCourseCounterService() CourseCounterService {
	return CourseCounterService{}
}

func (s CourseCounterService) Increase(id string) error {
	log.Println("Increasing course counter for course", id)
	return errors.New("Drama error")
}
