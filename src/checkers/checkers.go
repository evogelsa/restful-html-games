package checkers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var ROOT_PATH string

var mutex sync.Mutex

var board = [][]byte{
	{'-', 'X', '-', 'X', '-', 'X', '-', 'X'},
	{'X', '-', 'X', '-', 'X', '-', 'X', '-'},
	{'-', 'X', '-', 'X', '-', 'X', '-', 'X'},
	{' ', '-', ' ', '-', ' ', '-', ' ', '-'},
	{'-', ' ', '-', ' ', '-', ' ', '-', ' '},
	{'O', '-', 'O', '-', 'O', '-', 'O', '-'},
	{'-', 'O', '-', 'O', '-', 'O', '-', 'O'},
	{'O', '-', 'O', '-', 'O', '-', 'O', '-'},
}

var boardPrev [8][8]byte

var boardFormat []interface{}

var lastRow int = -1
var lastCol int = -1

var player int
var playerColor = []string{"Red", "White"}

func tokenReplace(t byte) []interface{} {
	switch t {
	case 'X':
		return []interface{}{
			`<font style="color:white"><b>(</b></font>`,
			`<font style="color:white"><b>)</b></font>`,
		}
	case 'x':
		return []interface{}{
			`<font style="color:white"><b>{</b></font>`,
			`<font style="color:white"><b>}</b></font>`,
		}
	case 'O':
		return []interface{}{
			`<font style="color:red"><b>(</b></font>`,
			`<font style="color:red"><b>)</b></font>`,
		}
	case 'o':
		return []interface{}{
			`<font style="color:red"><b>{</b></font>`,
			`<font style="color:red"><b>}</b></font>`,
		}
	case ' ':
		return []interface{}{"&nbsp", "&nbsp"}
	case '-':
		return []interface{}{"BLANK", "BLANK"}
	case '#':
		return []interface{}{
			`<font style="color:blue"><b>(</b></font>`,
			`<font style="color:blue"><b>)</b></font>`,
		}
	case '$':
		return []interface{}{
			`<font style="color:blue"><b>(</b></font>`,
			`<font style="color:blue"><b>)</b></font>`,
		}
	case '3':
		return []interface{}{
			`<font style="color:blue"><b>{</b></font>`,
			`<font style="color:blue"><b>}</b></font>`,
		}
	case '4':
		return []interface{}{
			`<font style="color:blue"><b>{</b></font>`,
			`<font style="color:blue"><b>}</b></font>`,
		}
	}
	return []interface{}{"ERROR"}
}

func getBoardFormat() {
	boardFormat = append(boardFormat, playerColor[player])
	var winString string
	playerWin := checkWin()
	if playerWin == 0 {
		winString = `Red wins!`
	} else if playerWin == 1 {
		winString = `White wins!`
	} else {
		winString = `&nbsp`
	}
	boardFormat = append(boardFormat, winString)
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			replaced := tokenReplace(board[i][j])
			if replaced[0] != "BLANK" {
				tokens := []interface{}{
					fmt.Sprintf(
						`<a href="`+ROOT_PATH+`/checkersTurn/%s/%s" style="text-decoration:none">%s</a>`,
						strconv.Itoa(i),
						strconv.Itoa(j),
						replaced[0],
					),
					fmt.Sprintf(
						`<a href="`+ROOT_PATH+`/checkersTurn/%s/%s" style="text-decoration:none">%s</a>`,
						strconv.Itoa(i),
						strconv.Itoa(j),
						replaced[1],
					),
				}
				boardFormat = append(boardFormat, tokens...)
			} else {
				boardFormat = append(boardFormat, "&nbsp", "&nbsp")
			}
		}
	}
}

func setBoardPrev() {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			t := board[i][j]
			if t != '#' && t != '$' && t != '3' && t != '4' {
				boardPrev[i][j] = board[i][j]
			} else if t == '#' {
				boardPrev[i][j] = 'X'
			} else if t == '$' {
				boardPrev[i][j] = 'O'
			} else if t == '3' {
				boardPrev[i][j] = 'x'
			} else if t == '4' {
				boardPrev[i][j] = 'o'
			}
		}
	}
}

func move(row, col int) {
	switchPlayer()
	t := board[lastRow][lastCol]
	if t == '#' {
		t = 'X'
	} else if t == '$' {
		t = 'O'
	} else if t == '3' {
		t = 'x'
	} else if t == '4' {
		t = 'o'
	}
	board[lastRow][lastCol] = ' '
	if row == 0 && t == 'O' {
		t = 'o'
	} else if row == 7 && t == 'X' {
		t = 'x'
	}
	board[row][col] = t
	lastRow, lastCol = -1, -1
}

func switchPlayer() {
	player = int(math.Abs(float64(player - 1)))
}

func checkWin() int {
	countX, countO := 0, 0
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			if board[i][j] == 'X' || board[i][j] == 'x' {
				countX++
			} else if board[i][j] == 'O' || board[i][j] == 'o' {
				countO++
			}
		}
	}
	if countO == 0 {
		// no os so x wins
		return 1
	} else if countX == 0 {
		// no xs so o wins
		return 0
	}
	return -1
}

func Turn(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	// fmt.Println("Endpoint hit: /checkersTurn")
	vars := mux.Vars(r)
	row, err := strconv.Atoi(vars["row"])
	if err != nil {
		panic(err)
	}
	col, err := strconv.Atoi(vars["col"])
	if err != nil {
		panic(err)
	}
	if lastRow > -1 && lastCol > -1 {
		// moving to empty space
		if board[row][col] == ' ' {
			// check distance valid
			disRow := math.Abs(float64(row - lastRow))
			disCol := math.Abs(float64(col - lastCol))
			if disRow == 1 && disCol == 1 {
				setBoardPrev()
				move(row, col)
			} else if disRow == 2 && disCol == 2 {
				difRow := row - lastRow
				difCol := col - lastCol
				opponent := int(math.Abs(float64(player - 1)))
				tokens := []byte{'O', 'X', 'o', 'x'}
				// up to left
				if difRow < 0 && difCol < 0 {
					if board[lastRow-1][lastCol-1] == tokens[opponent] ||
						board[lastRow-1][lastCol-1] == tokens[opponent+2] {
						setBoardPrev()
						board[lastRow-1][lastCol-1] = ' '
						move(row, col)
					}
					// up to right
				} else if difRow < 0 && difCol > 0 {
					if board[lastRow-1][lastCol+1] == tokens[opponent] ||
						board[lastRow-1][lastCol+1] == tokens[opponent+2] {
						setBoardPrev()
						board[lastRow-1][lastCol+1] = ' '
						move(row, col)
					}
					// down to right
				} else if difRow > 0 && difCol > 0 {
					if board[lastRow+1][lastCol+1] == tokens[opponent] ||
						board[lastRow+1][lastCol+1] == tokens[opponent+2] {
						setBoardPrev()
						board[lastRow+1][lastCol+1] = ' '
						move(row, col)
					}
					// down to left
				} else if difRow > 0 && difCol < 0 {
					if board[lastRow+1][lastCol-1] == tokens[opponent] ||
						board[lastRow+1][lastCol-1] == tokens[opponent+2] {
						setBoardPrev()
						board[lastRow+1][lastCol-1] = ' '
						move(row, col)
					}
				}
			}
		}
	} else if board[row][col] != ' ' {
		if board[row][col] == 'O' && player == 0 ||
			board[row][col] == 'X' && player == 1 ||
			board[row][col] == 'o' && player == 0 ||
			board[row][col] == 'x' && player == 1 {
			lastRow, lastCol = row, col
			t := board[lastRow][lastCol]
			if t == 'X' {
				board[lastRow][lastCol] = '#'
			} else if t == 'O' {
				board[lastRow][lastCol] = '$'
			} else if t == 'x' {
				board[lastRow][lastCol] = '3'
			} else if t == 'o' {
				board[lastRow][lastCol] = '4'
			}
		}
	}

	mutex.Unlock()
	http.Redirect(w, r, ROOT_PATH+"/checkers", http.StatusFound)
}

func Reset(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	board = [][]byte{
		{'-', 'X', '-', 'X', '-', 'X', '-', 'X'},
		{'X', '-', 'X', '-', 'X', '-', 'X', '-'},
		{'-', 'X', '-', 'X', '-', 'X', '-', 'X'},
		{' ', '-', ' ', '-', ' ', '-', ' ', '-'},
		{'-', ' ', '-', ' ', '-', ' ', '-', ' '},
		{'O', '-', 'O', '-', 'O', '-', 'O', '-'},
		{'-', 'O', '-', 'O', '-', 'O', '-', 'O'},
		{'O', '-', 'O', '-', 'O', '-', 'O', '-'},
	}

	mutex.Unlock()
	http.Redirect(w, r, ROOT_PATH+"/checkers", http.StatusFound)
}

func Undo(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	// fmt.Println("Endpoint hit: /checkersUndo")
	if lastRow != -1 && lastCol != -1 {
		t := board[lastRow][lastCol]
		if t == '#' {
			board[lastRow][lastCol] = 'X'
		} else if t == '$' {
			board[lastRow][lastCol] = 'O'
		} else if t == '3' {
			board[lastRow][lastCol] = 'x'
		} else if t == '4' {
			board[lastRow][lastCol] = 'o'
		}
		lastRow, lastCol = -1, -1
	} else {
		for i := 0; i < len(boardPrev); i++ {
			for j := 0; j < len(boardPrev[0]); j++ {
				board[i][j] = boardPrev[i][j]
			}
		}
		switchPlayer()
	}
	mutex.Unlock()
	http.Redirect(w, r, ROOT_PATH+"/checkers", http.StatusFound)
}

func PlayerSwitch(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	// fmt.Println("Endpoint hit: /checkersPlayer")
	switchPlayer()
	mutex.Unlock()
	http.Redirect(w, r, ROOT_PATH+"/checkers", http.StatusFound)
}

func Board(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Endpoint hit: /checkers")
	boardFormat = []interface{}{}
	getBoardFormat()
	fmt.Fprintf(
		w,
		`
		<head>
			<meta http-equiv="refresh" content="2" />
		</head>
		<div><h2>%s player turn</h2></div>
		<div><h3>%s</h2></div>
		<div><h4><a href="`+ROOT_PATH+`/checkersPlayer">(switch player / double jump)</a></h4></div>
		<font face="monospace" size="14">
		<div>
			<span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span>
		</div>
		<div>
			<span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span>
		</div>
		<div>
			<span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span>
		</div>
		<div>
			<span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span>
		</div>
		<div>
			<span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span>
		</div>
		<div>
			<span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span>
		</div>
		<div>
			<span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span>
		</div>
		<div>
			<span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#000000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span><!--
			--><span style="background-color:#FF0000">%s</span>
		</div></font>
		<div><h4><a href="`+ROOT_PATH+`/checkersUndo">(undo last move)</a></h4></div>
		<div><h4><a href="`+ROOT_PATH+`/checkersReset">(reset game)</a></h4></div>
		<div><h4><a href="/">(return to home)</a></h4></div>
		`,
		boardFormat...,
	)
}
