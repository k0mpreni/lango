package handler

import (
	"lango/cmd/web/view/home"
	"lango/cmd/web/view/pricing"
	"net/http"
)

// func HelloWebHandler(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "Bad Request", http.StatusBadRequest)
// 	}

// 	name := r.FormValue("name")
// 	component := HelloPost(name)
// 	err = component.Render(r.Context(), w)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		log.Fatalf("Error rendering in HelloWebHandler: %e", err)
// 	}
// }

func HomeHandler(w http.ResponseWriter, r *http.Request) error {
	plans := []pricing.Plan{
		{Id: "1", Name: "Personal", Price: "39", Popular: false},
		{Id: "2", Name: "Professional", Price: "79", Popular: true},
	}
	return render(r, w, home.Home(plans))
}
