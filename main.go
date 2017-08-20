package main

import (
	"net/http"

	"github.com/ttimasdf/getcet/handler"
	"github.com/urfave/negroni"
)

func main() {
	middle := negroni.New()
	middle.Use(negroni.NewLogger())
	middle.Use(negroni.NewStatic(http.Dir("public")))

	router := http.NewServeMux()
	router.HandleFunc("/ticket", handler.TicketHandler)
	router.HandleFunc("/score", handler.ScoreHandler)
	middle.UseHandler(router)
	http.ListenAndServe(":8000", middle)
}
