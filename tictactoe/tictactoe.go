package tictactoe

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var board = [][]string{
	{
		"<a href=\"/tttTurn/0/0\">&nbsp</a> ",
		"<a href=\"/tttTurn/0/1\">&nbsp</a> ",
		"<a href=\"/tttTurn/0/2\">&nbsp</a> ",
	},
	{
		"<a href=\"/tttTurn/1/0\">&nbsp</a> ",
		"<a href=\"/tttTurn/1/1\">&nbsp</a> ",
		"<a href=\"/tttTurn/1/2\">&nbsp</a> ",
	},
	{
		"<a href=\"/tttTurn/2/0\">&nbsp</a> ",
		"<a href=\"/tttTurn/2/1\">&nbsp</a> ",
		"<a href=\"/tttTurn/2/2\">&nbsp</a> ",
	},
}

var player rune

var mutex sync.Mutex

func Turn(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Println("Endpoint hit: /tttTurn")
	vars := mux.Vars(r)
	row, err := strconv.Atoi(vars["row"])
	if err != nil {
		fmt.Println("Bro you messed something up wtf")
		panic(err)
	}
	col, err := strconv.Atoi(vars["col"])
	if err != nil {
		fmt.Println("Bro you messed something up wtf")
		panic(err)
	}
	if player == 0 {
		board[row][col] = "X&nbsp"
		player = 1
	} else {
		board[row][col] = "O&nbsp"
		player = 0
	}

	http.Redirect(w, r, "/ttt", http.StatusFound)
}

func Board(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: /ttt")
	players := []rune{'X', 'O'}
	fmt.Fprintf(
		w,
		`
			<head>
				<meta http-equiv="refresh" content="2" />
			</head>
			<div><h2>%cs Turn</h2></div>
			<font face="monospace" size="14">
				<div>%s|&nbsp%s|&nbsp%s</div>
				<div>--|---|--</div>
				<div>%s|&nbsp%s|&nbsp%s</div>
				<div>--|---|--</div>
				<div>%s|&nbsp%s|&nbsp%s</div>
			</font>
		`,
		players[player],
		board[0][0], board[0][1], board[0][2],
		board[1][0], board[1][1], board[1][2],
		board[2][0], board[2][1], board[2][2],
	)
	checkBoard(w)
	fmt.Fprintf(
		w,
		`
			<div><h4><a href="/tttReset">(click to reset)</a></h4></div>
			<div><h4><a href="/">(click to return to home)</a></h4></div>
		`,
	)
}

func Reset(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Println("Endpoint hit: /tttReset")
	board = [][]string{
		{
			"<a href=\"/tttTurn/0/0\">&nbsp</a> ",
			"<a href=\"/tttTurn/0/1\">&nbsp</a> ",
			"<a href=\"/tttTurn/0/2\">&nbsp</a> ",
		},
		{
			"<a href=\"/tttTurn/1/0\">&nbsp</a> ",
			"<a href=\"/tttTurn/1/1\">&nbsp</a> ",
			"<a href=\"/tttTurn/1/2\">&nbsp</a> ",
		},
		{
			"<a href=\"/tttTurn/2/0\">&nbsp</a> ",
			"<a href=\"/tttTurn/2/1\">&nbsp</a> ",
			"<a href=\"/tttTurn/2/2\">&nbsp</a> ",
		},
	}

	http.Redirect(w, r, "/ttt", http.StatusFound)
}

func checkBoard(w http.ResponseWriter) {
	if checkWin() {
		// player 0 went last, Xs win
		if player == 1 {
			fmt.Fprintf(w, "<div><h2>Xs won!</h2></div>")
		} else {
			fmt.Fprintf(w, "<div><h2>Os won!</h2></div>")
		}
	} else if checkDraw() {
		fmt.Fprintf(w, "<div><h2>It's a draw, everyone loses</h2></div>")
	}
}

func checkWin() bool {
	tokens := []byte{'X', 'O'}
	var count int
	var t byte

	for n := 0; n < 2; n++ {
		t = tokens[n]

		// check horizontals
		for i := 0; i < 3; i++ {
			count = 0
			for j := 0; j < 3; j++ {
				if board[i][j][0] == t {
					count++
				}
			}
			if count == 3 {
				return true
			}
		}

		// check verticals
		for j := 0; j < 3; j++ {
			count = 0
			for i := 0; i < 3; i++ {
				if board[i][j][0] == t {
					count++
				}
			}
			if count == 3 {
				return true
			}
		}

		// check top left to bottom right
		count = 0
		for i := 0; i < 3; i++ {
			if board[i][i][0] == t {
				count++
			}
		}
		if count == 3 {
			return true
		}

		// check bottom right to top left
		count = 0
		for i, j := 2, 0; i > -1; i, j = i-1, j+1 {
			if board[i][j][0] == t {
				count++
			}
		}
		if count == 3 {
			return true
		}

	}

	return false
}

func checkDraw() bool {
	count := 0
	// check board full
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j][0] == '<' {
				count++
			}
		}
	}
	return count == 0
}
