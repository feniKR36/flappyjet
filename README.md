## Ready...jet set go!

## Project Overview:
 This is a copy of the OG Flappy Bird Game which is originally written in Python using the Pygame library and then later translated into Go using the Ebiten game engine to explore and compare the language. This project demonstrates:
 -Basic game dev concepts
 -Object movement and gravity simulation
 -Collision detection(without effects)
 -File based score saving

## Features:
 -User Login: Players can input their username before playing that allows the game to store individual highscores.
 -Score tracking: Can record the player's highest score and store it in a text file(score.txt)
 -Player can view their personal best score.

## Code and Design Structures:
 The Python version is written using pygame that handles graphics, input events, and game loop. The program is structures susing multiple functions to manage different parts of the syste, including login screen, main menu, profile screen, leaderboard, and again- game loop. run_game() is the function that manages the gameplay mechanics like jet movement, tower generation, collision detection, and score updates that is stored in (scores.txt).
 
 The Go version is implemented using Ebiten game engine that provides the needed tools for this project. Logic is organized in "Game" and stores the jet's position, velocity, etc,. It follows Ebiten's required interface like Update() that processes player input, applies grvity, and updates tower positions, Draw() renders the bg, jet, towers, and score, and Layout(). Program behins in the main() function that loads the game resources. ebiten.RunGame() starts the game loop.

 P.S. The structure still needed to be updated coz it is chaotic though it works

# Screenshots
<img src="mainwindow.png" alt="Main Window Screenshot" width="1100" height="700" />
<img src="ingamess.png" alt="In-game Screenshot" width="1100" height="700" />
<img src="leaderboard.png" alt="Report Screenshot" width="1100" height="700" />

## How to run the program
-Python Version:

1. Install Python
Download and install Python from: https://python.org
2. Install Pygame
Open terminal and run: pip install pygame
3. Run the game: python main.py

-Go Version:
1. Install Go
Download Go from: https://go.dev
2. Install Ebiten
go get github.com/hajimehoshi/ebiten/v2
3. Run the game: 
go tidy
go run main.go