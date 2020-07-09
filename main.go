package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/evogelsa/restful-html-games/checkers"
	"github.com/evogelsa/restful-html-games/tictactoe"

	"github.com/gorilla/mux"
)

const (
	ROOT_PATH = "/go"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<div>Hello, World!</div>")
	fmt.Fprintf(w, "<div><a href=\""+ROOT_PATH+"/ttt\">click to play tictactoe!</a></div>")
	fmt.Fprintf(w, "<div><a href=\""+ROOT_PATH+"/checkers\">click to play terrible checkers</a></div>")
}

func newRouter() *mux.Router {
	//define root paths
	checkers.ROOT_PATH = ROOT_PATH
	tictactoe.ROOT_PATH = ROOT_PATH

	// create new router
	router := mux.NewRouter()

	// add paths, using explicit methods here technically unnecessary
	router.HandleFunc(ROOT_PATH, homePage).Methods("GET")
	router.HandleFunc(ROOT_PATH+"/tttTurn/{row}/{col}", tictactoe.Turn).Methods("GET")
	router.HandleFunc(ROOT_PATH+"/tttReset", tictactoe.Reset).Methods("GET")
	router.HandleFunc(ROOT_PATH+"/ttt", tictactoe.Board).Methods("GET")
	router.HandleFunc(ROOT_PATH+"/checkersTurn/{row}/{col}", checkers.Turn).Methods("GET")
	router.HandleFunc(ROOT_PATH+"/checkersReset", checkers.Reset).Methods("GET")
	router.HandleFunc(ROOT_PATH+"/checkersUndo", checkers.Undo).Methods("GET")
	router.HandleFunc(ROOT_PATH+"/checkersPlayer", checkers.PlayerSwitch).Methods("GET")
	router.HandleFunc(ROOT_PATH+"/checkers", checkers.Board).Methods("GET")

	return router
}

func main() {
	router := newRouter()
	log.Fatal(http.ListenAndServeTLS(
		":8080",
		"/etc/letsencrypt/live/ethanvogelsang.xyz/cert.pem",
		"/etc/letsencrypt/live/ethanvogelsang.xyz/privkey.pem",
		router,
	))
}
