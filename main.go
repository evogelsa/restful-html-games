package main

import (
	"fmt"
	"log"
	"net/http"

	"github.iu.edu/ise-engr-e222/sp20-id-0007/rest-ttt/src/checkers"
	"github.iu.edu/ise-engr-e222/sp20-id-0007/rest-ttt/src/tictactoe"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Enpoint hit: /")
	fmt.Fprintf(w, "<div>Hello, World!</div>")
	fmt.Fprintf(w, "<div><a href=\"/ttt\">click to play tictactoe!</a></div>")
	fmt.Fprintf(w, "<div><a href=\"/checkers\">click to play terrible checkers</a></div>")
}

func newRouter() *mux.Router {
	// create new router
	router := mux.NewRouter()

	// add paths, using explicit methods here technically unnecessary
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/tttTurn/{row}/{col}", tictactoe.Turn).Methods("GET")
	router.HandleFunc("/tttReset", tictactoe.Reset).Methods("GET")
	router.HandleFunc("/ttt", tictactoe.Board).Methods("GET")
	router.HandleFunc("/checkersTurn/{row}/{col}", checkers.Turn).Methods("GET")
	router.HandleFunc("/checkersReset", checkers.Reset).Methods("GET")
	router.HandleFunc("/checkersUndo", checkers.Undo).Methods("GET")
	router.HandleFunc("/checkersPlayer", checkers.PlayerSwitch).Methods("GET")
	router.HandleFunc("/checkers", checkers.Board).Methods("GET")

	return router
}

func main() {
	router := newRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
