package handler

import (
	"fmt"
	"lango/cmd/web/view/courses"
	"lango/internal/database/domain"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func CoursesHandler(w http.ResponseWriter, r *http.Request) error {
	c := []domain.Course{
		{
			Title:       "Mathematics 101",
			Description: "You will learn the fundamentals of mathematics",
			Completed:   false,
			Canceled:    false,
			Date:        time.Now(),
		},
		{
			Title:       "Computer Science 101",
			Description: "You will learn the fundamentals of computer science",
			Completed:   true,
			Canceled:    false,
			Date:        time.Now(),
		},
		{
			Title:       "Geology 101",
			Description: "Canceled because who wants to learn about rocks",
			Completed:   false,
			Canceled:    true,
			Date:        time.Now(),
		},
	}

	return render(r, w, courses.CoursesList(c))
}

func CoursesCreateHandler(w http.ResponseWriter, r *http.Request) error {
	c := courses.CourseForm{}
	e := courses.CourseFormErrors{}

	return render(r, w, courses.CreateCourse(c, e))
}

func CoursesJoinHandler(w http.ResponseWriter, r *http.Request) error {
	courseIdParam := chi.URLParam(r, "courseId")
	fmt.Println("Course ID PARAM", courseIdParam)
	c := domain.Course{
		Title:       "Mathematics 101",
		Description: "You will learn the fundamentals of mathematics",
		Completed:   false,
		Canceled:    false,
		Date:        time.Now(),
	}

	return render(r, w, courses.JoinCourse(c))
}
