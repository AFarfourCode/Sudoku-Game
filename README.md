# Sudoku-Game
This project started as a Sudoku Solver for the 01 Talent Academy Piscine stage. The original task was to implement an algorithm that generates and solves Sudoku boards.
I decided to improve the challenge by turning it into a real interactive Sudoku game with a modern UI.


✨ Features

✅ Automatic Sudoku Solver – backtracking algorithm generates valid full solutions.

🎲 Difficulty Levels – Easy, Medium, Hard (different number of clues).

🖥 Graphical Interface (Fyne) – clean and colorful board, built with Fyne
.

🎯 Interactive Gameplay

Tap cells to select

Number pad to input values

Clear button to reset a cell

💡 Hint System – reveal the correct number for a selected cell.

🔄 Reset & New Game – restart current puzzle or generate a new one.

🎨 Visual Design Improvements

Highlight selected cell

Different colors for 3×3 sub-grids

Bold fixed numbers

Gradient background & animated title for fun UI.


🧠 How it Works

Uses a backtracking algorithm to generate a fully solved Sudoku board.

Randomly removes numbers depending on difficulty (while keeping a valid solution).

User interacts with the board using Fyne widgets.

The Check button compares the user’s progress with the solved board.


🚀 Technologies

Language: Go (Golang)

GUI Framework: Fyne

Algorithms: Backtracking Sudoku solver
