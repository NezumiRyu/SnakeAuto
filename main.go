package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type BoardDesc struct {
	boxX     int
	boxY     int
	spaceI   int
	spaceO   int
	boxesX   int
	boxesY   int
	headPosX int
	headPosY int
}

func (boardDesc *BoardDesc) windowX() float64 {
	return float64(2*boardDesc.spaceO + boardDesc.boxesX*boardDesc.boxX + (boardDesc.boxesX-1)*boardDesc.spaceI)
}

func (boardDesc *BoardDesc) windowY() float64 {
	return float64(2*boardDesc.spaceO + boardDesc.boxesY*boardDesc.boxY + (boardDesc.boxesY-1)*boardDesc.spaceI)
}

func (boardDesc *BoardDesc) Draw(win *pixelgl.Window, board *[][]box) {
	for i := 0; i < boardDesc.boxesX; i++ {
		for j := 0; j < boardDesc.boxesY; j++ {
			(*board)[j][i].getImd(boardDesc).Draw(win)
		}
	}
}

type box struct {
	x, y         int
	prevX, prevY int
	t            int //Type
	/*
		0. empty
		1. goal
		2. head
		3. body
		4. tail
	*/
}

func Box(x, y int) box {
	var box box
	box.x = x
	box.y = y
	return box
}

func (box *box) getColor() pixel.RGBA {
	switch box.t {
	case 1:
		return pixel.RGB(1, 0, 0) //RED
	case 2:
		return pixel.RGB(1, 1, 0) //Yellow
	case 3:
		return pixel.RGB(1, 1, 1) //White
	case 4:
		return pixel.RGB(1, 0.5, 0) //Brown
	default:
		return pixel.RGB(0, 0, 0)
	}
}

func (box *box) getImd(bd *BoardDesc) *imdraw.IMDraw {
	imd := imdraw.New(nil)

	imd.Color = box.getColor()
	imd.Push(
		pixel.V(
			float64(bd.spaceO+box.x*(bd.boxX+bd.spaceI)),
			float64(bd.spaceO+box.y*(bd.boxX+bd.spaceI))),
		pixel.V(
			float64(bd.spaceO+box.x*(bd.boxX+bd.spaceI)+bd.boxX),
			float64(bd.spaceO+box.y*(bd.boxX+bd.spaceI))),
		pixel.V(
			float64(bd.spaceO+box.x*(bd.boxX+bd.spaceI)+bd.boxX),
			float64(bd.spaceO+box.y*(bd.boxX+bd.spaceI)+bd.boxY)),
		pixel.V(
			float64(bd.spaceO+box.x*(bd.boxX+bd.spaceI)),
			float64(bd.spaceO+box.y*(bd.boxX+bd.spaceI)+bd.boxY)))
	imd.Polygon(0)

	return imd
}

func run() {

	var boardDesc BoardDesc
	boardDesc.boxX = 50
	boardDesc.boxY = 50
	boardDesc.spaceI = 5
	boardDesc.spaceO = 5
	boardDesc.boxesX = 10
	boardDesc.boxesY = 10

	board := [][]box{}
	for i := 0; i < boardDesc.boxesX; i++ {
		row := []box{}
		for j := 0; j < boardDesc.boxesY; j++ {
			row = append(row, Box(i, j))
		}
		board = append(board, row)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "SnakeAuto",
		Bounds: pixel.R(0, 0, boardDesc.windowX(), boardDesc.windowY()),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//START
	//HEAD
	boardDesc.headPosX = boardDesc.boxesX / 2
	boardDesc.headPosY = boardDesc.boxesY / 2
	board[boardDesc.headPosX][boardDesc.headPosY].t = 2
	board[boardDesc.headPosX][boardDesc.headPosY].prevX = boardDesc.headPosX - 1
	board[boardDesc.headPosX][boardDesc.headPosY].prevY = boardDesc.headPosY
	//BODY
	board[boardDesc.headPosX-1][boardDesc.headPosY].t = 3
	board[boardDesc.headPosX-1][boardDesc.headPosY].prevX = boardDesc.headPosX - 2
	board[boardDesc.headPosX-1][boardDesc.headPosY].prevY = boardDesc.headPosY
	//TAIL
	board[boardDesc.headPosX-2][boardDesc.headPosY].t = 4
	//GOAL
	board[boardDesc.headPosX+2][boardDesc.headPosY].t = 1

	for !win.Closed() {
		win.Clear(colornames.Black)
		boardDesc.Draw(win, &board)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
