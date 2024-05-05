package creating

import (
	"app-services-go/internal/courses/application/increasing"
	course "app-services-go/internal/courses/domain"
	"app-services-go/kit/event"
	"errors"
)

type IncreaseCoursesCounterOnCourseCreated struct {
	increasingService increasing.CourseCounterService
}

func NewIncreaseCoursesCounterOnCourseCreated(increaserService increasing.CourseCounterService) IncreaseCoursesCounterOnCourseCreated {
	return IncreaseCoursesCounterOnCourseCreated{
		increasingService: increaserService,
	}
}

func (e IncreaseCoursesCounterOnCourseCreated) On(evt event.Event) error {
	courseCreatedEvt, ok := evt.(course.CourseCreatedEvent)
	if !ok {
		return errors.New("unexpected event")
	}
	return e.increasingService.Increase(courseCreatedEvt.ID())
}

func (e IncreaseCoursesCounterOnCourseCreated) SubscribedTo() event.Event {
	return event.Event(course.CourseCreatedEvent{})
}
