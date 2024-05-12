package handler

import (
	"lango/cmd/web/view/pricing"
	"net/http"
)

func HandlePricingIndex(w http.ResponseWriter, r *http.Request) error {
	plans := []pricing.Plan{
		{Id: "1", Name: "Personal", Price: "39", Popular: false},
		{Id: "2", Name: "Professional", Price: "79", Popular: true},
	}
	return render(r, w, pricing.Pricing(plans))
}
