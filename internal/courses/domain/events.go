package domain

import (
	"app-services-go/kit/event"
	"encoding/json"
)

const CourseCreatedEventType event.Type = "events.course.created"

type CourseCreatedEvent struct {
	event.BaseEvent
	Id       string `json:"id"`
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

func NewCourseCreatedEvent(id, name, duration string) CourseCreatedEvent {
	return CourseCreatedEvent{
		Id:       id,
		Name:     name,
		Duration: duration,

		BaseEvent: event.NewBaseEvent(id),
	}
}
func (e CourseCreatedEvent) ID() string {
	return e.EventID
}
func (e CourseCreatedEvent) Type() event.Type {
	return CourseCreatedEventType
}

func (e CourseCreatedEvent) CourseID() string {
	return e.Id
}

func (e CourseCreatedEvent) CourseName() string {
	return e.Name
}

func (e CourseCreatedEvent) CourseDuration() string {
	return e.Duration
}
func (e CourseCreatedEvent) FromJSON(data []byte) (event.Event, error) {
	course := &CourseCreatedEvent{}
	err := json.Unmarshal(data, course)
	if err != nil {
		return CourseCreatedEvent{}, err
	}
	return *course, nil
}
