package domain

import (
	"context"
	"errors"
	"fmt"

	"app-services-go/kit/event"
	"github.com/google/uuid"
)

var ErrInvalidCourseID = errors.New("invalid Course ID")

// CourseID represents the course unique identifier.

// NewCourseID instantiate the VO for CourseID
func NewCourseID(value string) (string, error) {
	v, err := uuid.Parse(value)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrInvalidCourseID, value)
	}

	return v.String(), nil

}

var ErrEmptyCourseName = errors.New("the field Course Name can not be empty")

// NewCourseName instantiate VO for CourseName
func NewCourseName(value string) (string, error) {
	if value == "" {
		return "", ErrEmptyCourseName
	}

	return value, nil
}

var ErrEmptyDuration = errors.New("the field Duration can not be empty")

// CourseDuration represents the course duration.
type CourseDuration struct {
	value string
}

func NewCourseDuration(value string) (string, error) {

	return value, nil
}

// Course is the data structure that represents a course.
type Course struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Duration string `json:"duration,omitempty"`

	events []event.Event
}

// CourseRepository defines the expected behaviour from a course storage.
type CourseRepository interface {
	Save(ctx context.Context, course Course) error
	FetchById(ctx context.Context, id string) (Course, error)
}

//go:generate mockery --case=snake --outpkg=storagemocks --output=../infrastructure/storage/storagemocks --name=CourseRepository

// NewCourse creates a new course.
func NewCourse(id, name, duration string) (Course, error) {
	idVO, err := NewCourseID(id)
	if err != nil {
		return Course{}, err
	}

	nameVO, err := NewCourseName(name)
	if err != nil {
		return Course{}, err
	}

	durationVO, err := NewCourseDuration(duration)
	if err != nil {
		return Course{}, err
	}

	course := Course{
		ID:       idVO,
		Name:     nameVO,
		Duration: durationVO,
	}
	course.Record(NewCourseCreatedEvent(idVO, nameVO, durationVO))
	return course, nil
}

// Record records a new domain event.
func (c *Course) Record(evt event.Event) {
	c.events = append(c.events, evt)
}

// PullEvents returns all the recorded domain events.
func (c Course) PullEvents() []event.Event {
	evt := c.events
	c.events = []event.Event{}

	return evt
}
