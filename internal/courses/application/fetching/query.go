package fetching

import (
	course "app-services-go/internal/courses/domain"
	"app-services-go/kit/query"
	"context"
	"errors"
)

const CourseQueryType query.Type = "query.fetching.course"

// CourseQuery is the query dispatched to create a new course.
type CourseQuery struct {
	id string
}

// NewCourseQuery creates a new CourseQuery.
func NewCourseQuery(id string) CourseQuery {
	return CourseQuery{
		id: id,
	}
}

func (c CourseQuery) Type() query.Type {
	return CourseQueryType
}

// CourseQueryHandler is the query controllers responsible for creating courses.
type CourseQueryHandler struct {
	courseRepository course.CourseRepository
}

// NewCourseQueryHandler initializes a new CourseQueryHandler.
func NewCourseQueryHandler(courseRepository course.CourseRepository) CourseQueryHandler {
	return CourseQueryHandler{
		courseRepository: courseRepository,
	}
}

// Handle implements the query.Handler interface.
func (h CourseQueryHandler) Handle(ctx context.Context, qry query.Query) (interface{}, error) {
	courseQry, ok := qry.(CourseQuery)
	if !ok {
		return nil, errors.New("unexpected query")
	}

	idVO, err := course.NewCourseID(courseQry.id)
	if err != nil {
		return nil, err
	}

	course, err := h.courseRepository.FetchById(ctx, idVO)
	if err != nil {
		return nil, err
	}

	return course, nil
}
