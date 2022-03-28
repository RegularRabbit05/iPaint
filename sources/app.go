package sources

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Size = [2]int32{1000, 700}
var Terminate = false
var mouse rl.Vector2
var exitIcons = [3]rl.Vector2{{20, 12}, {40, 12}, {60, 12}}
var Freeze = false
var Frame = 0
var isMinimizing = false
var shouldTerminate = false
var TitleBarSize = [2]int32{Size[0], 24}
var AreaSize = [2]int32{Size[0], Size[1] - TitleBarSize[1]}
var Icons [8]rl.Texture2D

func Setup() {
	rl.SetConfigFlags(rl.FlagWindowAlwaysRun)
	rl.SetConfigFlags(rl.FlagWindowTransparent)
	rl.SetConfigFlags(rl.FlagWindowUndecorated)
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(Size[0], Size[1], "iPaint")
	var icon = rl.LoadImage("../assets/icon.png")
	rl.SetWindowIcon(*icon)
	rl.UnloadImage(icon)
	rl.InitAudioDevice()

	rl.HideCursor()
	loadCursors()

	pushIcon(0, "save.png")
	pushIcon(1, "load.png")
	MakeCanvas()
}

func loadCursors() {
	var loader = rl.LoadImage("../assets/pen.png")
	Cursors[0] = rl.LoadTextureFromImage(loader)
	rl.UnloadImage(loader)
	loader = rl.LoadImage("../assets/erase.png")
	Cursors[1] = rl.LoadTextureFromImage(loader)
	rl.UnloadImage(loader)
	loader = rl.LoadImage("../assets/bucket.png")
	Cursors[2] = rl.LoadTextureFromImage(loader)
	rl.UnloadImage(loader)
	loader = rl.LoadImage("../assets/cur.png")
	Cursors[3] = rl.LoadTextureFromImage(loader)
	rl.UnloadImage(loader)
}

func pushIcon(i int, name string) {
	var loader = rl.LoadImage(fmt.Sprintf("../assets/%s", name))
	Icons[i] = rl.LoadTextureFromImage(loader)
	rl.SetTextureFilter(Icons[i], rl.FilterPoint)
	rl.UnloadImage(loader)
}

func unloadData() {
	for i := 0; i < len(Icons); i++ {
		rl.UnloadTexture(Icons[i])
	}
	for i := 0; i < len(Cursors); i++ {
		rl.UnloadTexture(Cursors[i])
	}
}

func Exit() {
	unloadData()
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

const stopLoadTextureProc = false

func beforeWindow() {
	Frame++
	if Frame >= 60 {
		Frame = 0
	}
	if !Freeze {
		mouse = rl.GetMousePosition()
	}
	if !stopLoadTextureProc {
		CanvasTexture = rl.LoadTextureFromImage(Canvas)
	}
	rl.BeginDrawing()
}

func afterWindow() {
	rl.EndDrawing()
	if !stopLoadTextureProc {
		rl.UnloadTexture(CanvasTexture)
	}
	if !Freeze {
		handleMenuBar()
	}
}

func drawWindow() {
	rl.ClearBackground(rl.Color{})
	drawDecorations()
	drawBlankArea()
}

func drawDecorations() {
	rl.DrawRectangle(10, 0, Size[0]-20, 24, rl.Color{R: 205, G: 205, B: 205, A: 255})
	rl.DrawCircle(10, 12, 12, rl.Color{R: 205, G: 205, B: 205, A: 255})
	rl.DrawCircle(Size[0]-10, 12, 12, rl.Color{R: 205, G: 205, B: 205, A: 255})
	rl.DrawRectangle(0, 12, Size[0], 12, rl.Color{R: 205, G: 205, B: 205, A: 255})
	rl.DrawCircle(int32(exitIcons[0].X), int32(exitIcons[0].Y), 6, rl.Color{R: 255, G: 95, B: 87, A: 255})
	rl.DrawCircle(int32(exitIcons[1].X), int32(exitIcons[1].Y), 6, rl.Color{R: 255, G: 190, B: 47, A: 255})
	rl.DrawCircle(int32(exitIcons[2].X), int32(exitIcons[2].Y), 6, rl.Color{R: 41, G: 204, B: 65, A: 255})
}

func drawBlankArea() {
	rl.DrawRectangle(0, TitleBarSize[1], AreaSize[0], AreaSize[1], rl.RayWhite)
}

var tmpFrame = 0
var winMinPos rl.Vector2
var winMinSize = rl.Vector2{X: float32(Size[0]), Y: float32(Size[1])}
var swapMinPos = false

func handleMenuBar() {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		if rl.CheckCollisionPointCircle(mouse, exitIcons[0], 6) {
			Freeze = true
			isMinimizing = true
			tmpFrame = 0
			swapMinPos = false
			winMinPos = rl.GetWindowPosition()
			shouldTerminate = true
		}
		if rl.CheckCollisionPointCircle(mouse, exitIcons[1], 6) {
			Freeze = true
			isMinimizing = true
			tmpFrame = 0
			swapMinPos = false
			winMinPos = rl.GetWindowPosition()
			shouldTerminate = false
			rl.SetTargetFPS(60)
		}
		if rl.CheckCollisionPointCircle(mouse, exitIcons[2], 6) {
			rl.ToggleFullscreen()
		}
	}
	HandleMovement()
}

func Minimize() {
	tmpFrame++
	if swapMinPos {
		swapMinPos = false
		rl.MinimizeWindow()
		rl.SetTargetFPS(-1)
		if shouldTerminate {
			Terminate = true
		}
		return
	}
	rl.SetWindowPosition(int(winMinPos.X)+tmpFrame*50, int(winMinPos.Y)+tmpFrame*100)
	rl.SetWindowSize(int(winMinSize.X/(float32(tmpFrame/(15/10)))), int(winMinSize.Y/(float32(tmpFrame/(15/10)))))
	if tmpFrame >= 60 {
		isMinimizing = false
		swapMinPos = true
		Freeze = false
		rl.SetWindowSize(int(winMinSize.X), int(winMinSize.Y))
		rl.SetWindowPosition(int(winMinPos.X), int(winMinPos.Y))
	}
}

var lastClickPos = [2]float32{-1, -1}

func HandleMovement() {
	if rl.IsMouseButtonDown(rl.MouseLeftButton) && rl.CheckCollisionPointRec(mouse, rl.Rectangle{
		Width: float32(TitleBarSize[0]), Height: float32(TitleBarSize[1])}) {
		var MouseX = float32(rl.GetMouseX())
		var MouseY = float32(rl.GetMouseY())
		if lastClickPos[0] != -1 && lastClickPos[1] != -1 {
			if MouseX-lastClickPos[0] == 0 && MouseY-lastClickPos[1] == 0 {
				goto skip
			}
			var posX = rl.GetWindowPosition().X
			var posY = rl.GetWindowPosition().Y
			posX += MouseX - float32(TitleBarSize[0]/2)
			posY += MouseY - 12
			rl.SetWindowPosition(
				int(posX),
				int(posY))
		}
	skip:
		lastClickPos[0] = MouseX
		lastClickPos[1] = MouseY
	} else {
		lastClickPos = [2]float32{-1, -1}
	}
}
