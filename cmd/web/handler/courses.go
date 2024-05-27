package handler

import (
	"lango/cmd/web/view/courses"
	"net/http"
)

func CoursesHandler(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, courses.Index())
}
