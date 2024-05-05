package providers

import (
	"app-services-go/configs"
	"app-services-go/internal/courses/application/creating"
	"app-services-go/internal/courses/application/fetching"
	"app-services-go/internal/courses/application/increasing"
	"app-services-go/internal/courses/infrastructure/storage/mysql"
	"app-services-go/kit/command"
	"app-services-go/kit/event"
	"app-services-go/kit/event/rabbitMQ"
	"app-services-go/kit/query"
	"database/sql"
)

func CoursesContainer(db *sql.DB, commandBus *command.CommandBus, queryBus *query.QueryBus, eventBus *rabbitMQ.EventBus, subscribers *[]event.Subscriber) {
	dbConf, _ := configs.GetDatabaseConfig()

	courseRepository := mysql.NewCourseRepository(db, dbConf)

	creatingCourseService := creating.NewCourseService(courseRepository, eventBus)
	increasingCourseService := increasing.NewCourseCounterService()

	createCourseCommandHandler := creating.NewCourseCommandHandler(creatingCourseService)
	commandBus.Register(creating.CourseCommandType, createCourseCommandHandler)

	fetchCourseQueryHandler := fetching.NewCourseQueryHandler(courseRepository)
	queryBus.Register(fetching.CourseQueryType, fetchCourseQueryHandler)

	*subscribers = append(*subscribers, creating.NewIncreaseCoursesCounterOnCourseCreated(increasingCourseService))

}
