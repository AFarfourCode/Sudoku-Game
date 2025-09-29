package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const size = 9

var (
	initialBoard  [9][9]byte
	solutionBoard [9][9]byte
	currentBoard  [9][9]byte

	selectedRow, selectedCol int = -1, -1
	gameWindow               fyne.Window
)

func main() {
	a := app.New()
	gameWindow = a.NewWindow("🎮 Sudoku Fun Edition")
	gameWindow.SetIcon(theme.FyneLogo())
	gameWindow.Resize(fyne.NewSize(560, 740))
	gameWindow.SetFixedSize(true)

	rand.Seed(time.Now().UnixNano())
	showDifficultySelection()

	gameWindow.ShowAndRun()
}

// ---------------------------------
// UI: Difficulty selection
// ---------------------------------
func showDifficultySelection() {
	title := canvas.NewText("🎲 Sudoku Fun 🎲", color.NRGBA{R: 255, G: 90, B: 95, A: 255})
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Animate title left-right
	go func() {
		dir := 1
		for {
			time.Sleep(100 * time.Millisecond)
			title.Move(fyne.NewPos(float32(title.Position().X+float32(5*dir)), title.Position().Y))
			if title.Position().X > 200 || title.Position().X < -200 {
				dir *= -1
			}
			title.Refresh()
		}
	}()

	easyBtn := widget.NewButton("😎 Easy", func() { startGame(45) })
	mediumBtn := widget.NewButton("🤔 Medium", func() { startGame(35) })
	hardBtn := widget.NewButton("🔥 Hard", func() { startGame(25) })

	easyBtn.Importance = widget.LowImportance
	mediumBtn.Importance = widget.MediumImportance
	hardBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(title),
		layout.NewSpacer(),
		widget.NewLabel("👉 Choose Your Adventure:"),
		container.NewCenter(container.NewVBox(
			easyBtn,
			mediumBtn,
			hardBtn,
		)),
		layout.NewSpacer(),
	)

	gameWindow.SetContent(container.NewMax(
		canvas.NewLinearGradient(
			color.NRGBA{R: 255, G: 230, B: 230, A: 255},
			color.NRGBA{R: 200, G: 230, B: 255, A: 255},
			0, // horizontal gradient
		),
		content,
	))
}

// ---------------------------------
// Game generator & solver
// ---------------------------------
func startGame(clues int) {
	generateBoard(clues)
	showGameUI()
}

func generateBoard(numbersToShow int) {
	solutionBoard = generateSolvedBoard()
	currentBoard = solutionBoard

	cells := make([][2]int, 0, 81)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			cells = append(cells, [2]int{i, j})
		}
	}

	rand.Shuffle(len(cells), func(i, j int) { cells[i], cells[j] = cells[j], cells[i] })

	// Keep first numbersToShow cells filled, clear the others
	for i := numbersToShow; i < len(cells); i++ {
		row, col := cells[i][0], cells[i][1]
		currentBoard[row][col] = '.'
	}

	initialBoard = currentBoard
}

func generateSolvedBoard() [9][9]byte {
	var board [9][9]byte
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			board[i][j] = '.'
		}
	}
	solve(&board)
	return board
}

func solve(board *[9][9]byte) bool {
	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			if board[r][c] == '.' {
				nums := [9]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}
				rand.Shuffle(9, func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })

				for _, n := range nums {
					if isValid(*board, r, c, n) {
						board[r][c] = n
						if solve(board) {
							return true
						}
						board[r][c] = '.'
					}
				}
				return false
			}
		}
	}
	return true
}

// ---------------------------------
// UI: Game board and controls
// ---------------------------------
func showGameUI() {
	title := canvas.NewText("🎮 Sudoku Playground 🎮", color.NRGBA{R: 255, G: 50, B: 120, A: 255})
	title.TextSize = 22
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Animated title bobbing effect
	go func() {
		dir := 1
		for {
			time.Sleep(150 * time.Millisecond)
			title.Move(fyne.NewPos(title.Position().X, title.Position().Y+float32(3*dir)))
			if title.Position().Y > 10 || title.Position().Y < -10 {
				dir *= -1
			}
			title.Refresh()
		}
	}()

	// Grid
	boardGrid := container.NewGridWithColumns(size)
	cellBackgrounds := make([][]*canvas.Rectangle, size)
	cellTexts := make([][]*canvas.Text, size)

	colors := []color.NRGBA{
		{R: 255, G: 230, B: 230, A: 255},
		{R: 230, G: 255, B: 230, A: 255},
		{R: 230, G: 230, B: 255, A: 255},
		{R: 255, G: 250, B: 200, A: 255},
	}

	for i := 0; i < size; i++ {
		cellBackgrounds[i] = make([]*canvas.Rectangle, size)
		cellTexts[i] = make([]*canvas.Text, size)
	}

	for r := 0; r < size; r++ {
		for c := 0; c < size; c++ {
			bg := canvas.NewRectangle(colors[(r/3+c/3)%len(colors)])
			bg.SetMinSize(fyne.NewSize(55, 55))

			val := currentBoard[r][c]
			if val == '.' {
				text := canvas.NewText("", color.NRGBA{R: 50, G: 50, B: 50, A: 255})
				text.Alignment = fyne.TextAlignCenter
				text.TextSize = 22
				cellTexts[r][c] = text

				rr, cc := r, c
				btn := widget.NewButton("", func() {
					onCellTapped(rr, cc, cellBackgrounds)
				})
				btn.Importance = widget.LowImportance

				cell := container.NewMax(bg, container.NewCenter(text), btn)
				boardGrid.Add(cell)
				cellBackgrounds[r][c] = bg
			} else {
				text := canvas.NewText(string(val), color.Black)
				text.TextStyle = fyne.TextStyle{Bold: true}
				text.Alignment = fyne.TextAlignCenter
				text.TextSize = 22
				cellTexts[r][c] = text

				cell := container.NewMax(bg, container.NewCenter(text))
				boardGrid.Add(cell)
				cellBackgrounds[r][c] = bg
			}
		}
	}

	// Number pad
	numberPad := container.NewGridWithColumns(9)
	for i := 1; i <= 9; i++ {
		n := i
		btn := widget.NewButton(string(rune('0'+n)), func() {
			if selectedRow >= 0 && selectedCol >= 0 && initialBoard[selectedRow][selectedCol] == '.' {
				currentBoard[selectedRow][selectedCol] = byte('0' + n)
				cellTexts[selectedRow][selectedCol].Text = string(rune('0' + n))
				cellTexts[selectedRow][selectedCol].Refresh()
			}
		})
		numberPad.Add(btn)
	}

	clearBtn := widget.NewButton("🧹 Clear", func() {
		if selectedRow >= 0 && selectedCol >= 0 && initialBoard[selectedRow][selectedCol] == '.' {
			currentBoard[selectedRow][selectedCol] = '.'
			cellTexts[selectedRow][selectedCol].Text = ""
			cellTexts[selectedRow][selectedCol].Refresh()
		}
	})

	checkBtn := widget.NewButton("✅ Check", func() {
		if isSolved() {
			dialog.ShowInformation("🎉 Winner!", "You solved the puzzle!", gameWindow)
		} else {
			dialog.ShowInformation("🤔 Not yet", "There are mistakes or missing numbers.", gameWindow)
		}
	})

	hintBtn := widget.NewButton("💡 Hint", func() {
		if selectedRow >= 0 && selectedCol >= 0 && initialBoard[selectedRow][selectedCol] == '.' {
			currentBoard[selectedRow][selectedCol] = solutionBoard[selectedRow][selectedCol]
			cellTexts[selectedRow][selectedCol].Text = string(solutionBoard[selectedRow][selectedCol])
			cellTexts[selectedRow][selectedCol].Refresh()
			cellTexts[selectedRow][selectedCol].TextStyle = fyne.TextStyle{Bold: true}
			cellTexts[selectedRow][selectedCol].Refresh()
		}
	})

	resetBtn := widget.NewButton("🔄 Reset", func() {
		for r := 0; r < size; r++ {
			for c := 0; c < size; c++ {
				if initialBoard[r][c] == '.' {
					currentBoard[r][c] = '.'
					cellTexts[r][c].Text = ""
					cellTexts[r][c].Refresh()
				}
				restoreCellBackground(r, c, cellBackgrounds)
			}
		}
		selectedRow, selectedCol = -1, -1
	})

	newGameBtn := widget.NewButton("🎲 New Game", showDifficultySelection)

	controls := container.NewVBox(
		container.NewHBox(checkBtn, hintBtn, resetBtn, newGameBtn),
		widget.NewLabel("🎯 Pick a number:"),
		numberPad,
		clearBtn,
	)

	content := container.NewVBox(
		container.NewCenter(title),
		container.New(layout.NewCenterLayout(), container.NewMax(boardGrid)),
		layout.NewSpacer(),
		controls,
	)

	gameWindow.SetContent(container.NewMax(
		canvas.NewRadialGradient(
			color.NRGBA{R: 255, G: 245, B: 200, A: 255},
			color.NRGBA{R: 200, G: 220, B: 255, A: 255},
		),
		content,
	))
}

// ---------------------------------
// Helpers
// ---------------------------------
func onCellTapped(r, c int, backgrounds [][]*canvas.Rectangle) {
	if initialBoard[r][c] != '.' {
		return
	}
	if selectedRow >= 0 && selectedCol >= 0 {
		restoreCellBackground(selectedRow, selectedCol, backgrounds)
	}
	selectedRow, selectedCol = r, c
	if backgrounds[r][c] != nil {
		backgrounds[r][c].FillColor = color.NRGBA{R: 255, G: 200, B: 120, A: 255}
		backgrounds[r][c].Refresh()
	}
}

func restoreCellBackground(r, c int, backgrounds [][]*canvas.Rectangle) {
	if backgrounds[r][c] == nil {
		return
	}
	colors := []color.NRGBA{
		{R: 255, G: 230, B: 230, A: 255},
		{R: 230, G: 255, B: 230, A: 255},
		{R: 230, G: 230, B: 255, A: 255},
		{R: 255, G: 250, B: 200, A: 255},
	}
	backgrounds[r][c].FillColor = colors[(r/3+c/3)%len(colors)]
	backgrounds[r][c].Refresh()
}

func isSolved() bool {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if currentBoard[i][j] != solutionBoard[i][j] {
				return false
			}
		}
	}
	return true
}

func isValid(board [9][9]byte, row, col int, char byte) bool {
	for i := 0; i < size; i++ {
		if board[row][i] == char || board[i][col] == char {
			return false
		}
	}
	startRow := row - row%3
	startCol := col - col%3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[startRow+i][startCol+j] == char {
				return false
			}
		}
	}
	return true
}
