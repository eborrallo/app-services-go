package mysql

import (
	"app-services-go/configs"
	course "app-services-go/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

// CourseRepository is a MySQL course.CourseRepository implementation.
type CourseRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}

// NewCourseRepository initializes a MySQL-based implementation of course.CourseRepository.
func NewCourseRepository(db *sql.DB, configs configs.DatabaseConfig) *CourseRepository {
	return &CourseRepository{
		db:        db,
		dbTimeout: configs.DbTimeout,
	}
}

// Save implements the course.CourseRepository interface.
func (r *CourseRepository) Save(ctx context.Context, course course.Course) error {
	courseSQLStruct := sqlbuilder.NewStruct(new(sqlCourse))
	query, args := courseSQLStruct.InsertInto(sqlCourseTable, sqlCourse{
		ID:       course.ID,
		Name:     course.Name,
		Duration: course.Duration,
	}).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist course on database: %v", err)
	}

	return nil
}

func (r *CourseRepository) FetchById(ctx context.Context, id string) (course.Course, error) {

	var entity course.Course

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("courses")
	sb.Where(sb.Equal("Id", id))

	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctxTimeout, query, args...).Scan(&entity.ID, &entity.Name, &entity.Duration)
	if err != nil {
		return course.Course{}, fmt.Errorf("error trying to fetch course %s on database: %v", id, err)
	}

	return entity, nil
}
