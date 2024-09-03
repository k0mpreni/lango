package handler

import (
	"context"
	"errors"
	"fmt"
	"lango/cmd/web/view/courses"
	"lango/internal/database"
	"lango/internal/database/domain"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func CoursesHandler(w http.ResponseWriter, r *http.Request) error {
	var c []domain.Course

	u := getAuthenticatedUser(r)
	ctx := context.Background()
	user, err := database.DB.GetUserByEmail(ctx, u.Email)
	if err != nil {
		fmt.Println("error get user courses", err)
		hxRedirect(w, r, "/login")
	}

	t, err := database.DB.GetTeacherByUserId(ctx, user.ID)
	if err != nil {
		fmt.Println("error getting teacher", err)
		return render(r, w, courses.CoursesList(c))
	}

	c, err = database.DB.ListCoursesByTeacherId(ctx, t.ID)
	if err != nil {
		fmt.Println("error courses", err)
		return err
	}

	return render(r, w, courses.CoursesList(c))
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
	fmt.Println("COURSE FORM", c)

	ctx := context.Background()

	u := getAuthenticatedUser(r)
	user, err := database.DB.GetUserByEmail(ctx, u.Email)
	if err != nil {
		e := courses.CourseFormErrors{
			Create: "Error finding user",
		}
		return render(r, w, courses.CreateCourse(c, e))
	}
	fmt.Println(u)

	teacher := domain.CreateTeacherParams{UserID: user.ID}
	// TODO: Do a Get or create
	t, err := database.DB.GetTeacherByUserId(ctx, user.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		t, err = database.DB.CreateTeacher(ctx, teacher)
		if err != nil {
			fmt.Println("error creating teacher account", err)
			e := courses.CourseFormErrors{
				Create: "Error creating teacher account",
			}
			return render(r, w, courses.CreateCourse(c, e))
		}
	}

	student, err := database.DB.GetUserByEmail(ctx, c.StudentEmail)
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

	course := domain.CreateCourseParams{
		TeacherID:   t.ID,
		StudentID:   student.ID,
		Title:       pgtype.Text{String: c.Title, Valid: true},
		Description: pgtype.Text{String: c.Description, Valid: true},
		Date: pgtype.Timestamptz{
			Valid: true,
			Time:  courseDate,
		},
	}

	createdCourse, err := database.DB.CreateCourse(ctx, course)
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

	fmt.Println("COURSE CREATED:", createdCourse)
	return render(r, w, courses.CreateCourse(successForm, e))
}

func CoursesJoinHandler(w http.ResponseWriter, r *http.Request) error {
	courseIdParam := chi.URLParam(r, "courseId")
	fmt.Println("Course ID PARAM", courseIdParam)
	c := domain.Course{
		Title:       pgtype.Text{String: "Mathematics 101"},
		Description: pgtype.Text{String: "You will learn the fundamentals of mathematics"},
		Completed:   false,
		Canceled:    false,
		Date: pgtype.Timestamptz{
			Time: time.Now(),
		},
	}

	return render(r, w, courses.JoinCourse(c))
}
