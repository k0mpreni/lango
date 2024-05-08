package handler

import (
	"net/http"

	"lango/cmd/web/view/pricing"
)

func HandlePricingIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, pricing.Pricing())
}
