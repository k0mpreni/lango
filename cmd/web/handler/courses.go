package handler

import (
	"errors"
	"fmt"
	"lango/cmd/web/view/courses"
	"lango/internal/database"
	"lango/internal/database/domain"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func CoursesHandler(w http.ResponseWriter, r *http.Request) error {
	// mockCourses := []domain.Course{
	// 	{
	// 		Title:       "Mathematics 101",
	// 		Description: "You will learn the fundamentals of mathematics",
	// 		Completed:   false,
	// 		Canceled:    false,
	// 		Date:        time.Now(),
	// 	},
	// 	{
	// 		Title:       "Computer Science 101",
	// 		Description: "You will learn the fundamentals of computer science",
	// 		Completed:   true,
	// 		Canceled:    false,
	// 		Date:        time.Now(),
	// 	},
	// 	{
	// 		Title:       "Geology 101",
	// 		Description: "Canceled because who wants to learn about rocks",
	// 		Completed:   false,
	// 		Canceled:    true,
	// 		Date:        time.Now(),
	// 	},
	// }

	u := getAuthenticatedUser(r)
	user, err := database.DB.Users.GetById(u.ID)
	if err != nil {
		fmt.Println("error get user courses", err)
		return err
	}

	t, err := database.DB.Teachers.GetById(user.Id)
	if err != nil {
		fmt.Println("error getting teacher", err)
		return err
	}

	fmt.Println(t)

	c, err := database.DB.Courses.GetAllByTeacher(t.Id)
	if err != nil {
		fmt.Println("error courses", err)
		return err
	}

	return render(r, w, courses.CoursesList(*c))
}

func CoursesCreateHandler(w http.ResponseWriter, r *http.Request) error {
	c := courses.CourseForm{}
	e := courses.CourseFormErrors{}

	return render(r, w, courses.CreateCoursePage(c, e))
}

func CoursesCreatePostHandler(w http.ResponseWriter, r *http.Request) error {
	c := courses.CourseForm{
		Title:        r.FormValue("title"),
		Description:  r.FormValue("description"),
		StudentEmail: r.FormValue("student-email"),
		Date:         r.FormValue("date"),
	}

	u := getAuthenticatedUser(r)
	user, err := database.DB.Users.GetById(u.ID)
	if err != nil {
		e := courses.CourseFormErrors{
			Create: "Error finding user",
		}
		return render(r, w, courses.CreateCourse(c, e))
	}
	fmt.Println(u)

	teacher := domain.Teacher{UserId: user.Id}
	// TODO: Do a Get or create
	t, err := database.DB.Teachers.GetById(user.Id)
	if errors.Is(err, domain.ErrRecordNotFound) {
		err := database.DB.Teachers.Create(&teacher)
		if err != nil {
			fmt.Println("error creating teacher account", err)
			e := courses.CourseFormErrors{
				Create: "Error creating teacher account",
			}
			return render(r, w, courses.CreateCourse(c, e))
		}
	}

	student, err := database.DB.Users.GetByEmail(c.StudentEmail)
	if err != nil {
		e := courses.CourseFormErrors{
			Student: "Student does not exists",
		}
		return render(r, w, courses.CreateCourse(c, e))
	}

	courseDate, err := time.Parse("2006-01-02", c.Date)
	if err != nil {
		e := courses.CourseFormErrors{
			Date: "Date is not valid",
		}
		return render(r, w, courses.CreateCourse(c, e))
	}

	// s, err := database.DB.Users.GetById(u.ID)

	course := domain.Course{
		CreatedAt:   time.Now(),
		TeacherId:   t.Id,
		StudentId:   student.Id,
		Title:       c.Title,
		Description: c.Description,
		Date:        courseDate,
	}

	err = database.DB.Courses.Create(&course)
	if err != nil {
		fmt.Println("error course creation", err)
		// Show error for creating course

		e := courses.CourseFormErrors{
			Create: "An error happened",
		}
		return render(r, w, courses.CreateCourse(c, e))
	}

	successForm := courses.CourseForm{
		Success: "Course created",
	}
	e := courses.CourseFormErrors{}

	return render(r, w, courses.CreateCourse(successForm, e))
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
