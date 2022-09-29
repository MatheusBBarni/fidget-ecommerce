package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MatheusBBarni/fidget-ecommerce/internal/cards"
	"github.com/go-chi/chi/v5"
)

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"string,omitempty"`
	Id      int    `json:"id,omitempty"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	amount, err := strconv.Atoi(payload.Amount)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	ok := true

	paymentIntent, message, err := card.Charge(payload.Currency, amount)

	if err != nil {
		ok = false
	}

	if ok {
		out, err := json.MarshalIndent(paymentIntent, "", "  ")

		if err != nil {
			app.errorLog.Println(err)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.Write(out)

		return
	}

	j := jsonResponse{
		Ok:      false,
		Message: message,
		Content: "",
	}

	out, err := json.MarshalIndent(j, "", "  ")

	if err != nil {
		app.errorLog.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(out)
}

func (app *application) GetWidgetById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetId, _ := strconv.Atoi(id)

	widget, err := app.DB.GetWidget(widgetId)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	out, err := json.MarshalIndent(widget, "", "  ")

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(out)
}
