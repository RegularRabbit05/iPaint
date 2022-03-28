package sources

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

func Loop() {
	for !rl.WindowShouldClose() && !Terminate {
		var skip = false
		beforeWindow()
		Size[0] = int32(rl.GetScreenWidth())
		Size[1] = int32(rl.GetScreenHeight())
		if swapMinPos {
			skip = true
		}
		if !skip {
			ToolPicker()
			switch DrawMode {
			case 0:
				PenToolHandler()
			case 1:
				EraserToolHandler()
			case 2:
				FillToolHandler()
			case 3:
				CursorToolHandler()
			default:
				DrawMode = 3
			}
			drawWindow()
			AppContent()
		}

		CursorHandler()
		afterWindow()

		if !skip {
			if !Freeze {
				AppHandle()
			}
		}

		if isMinimizing || swapMinPos {
			Minimize()
		}
	}
}

var ColorsPalette = [2][10]rl.Color{
	{rl.Purple, rl.Blue, rl.SkyBlue, rl.Green, rl.Yellow, rl.Orange, rl.Red, rl.Brown, rl.LightGray, rl.Black},
	{rl.DarkPurple, rl.DarkBlue, rl.Blue, rl.DarkGreen, rl.Gold, rl.Gold, rl.Color{R: 138, B: 7, A: 255}, rl.DarkBrown, rl.Black, rl.Black},
}

func DrawColorBox(id int32) {
	var mov = id * (10 + 40)
	rl.DrawRectangle(Size[0]-mov-1, TitleBarSize[1]+10-1, 40+2, 40+2, rl.Gray)
	rl.DrawRectangleGradientV(Size[0]-mov, TitleBarSize[1]+10, 40, 40, ColorsPalette[0][id-1], ColorsPalette[1][id-1])
}

func DrawToolsIcons() {
	var initPos = 80 + Icons[0].Width + 10
	var i int32
	for i = 1; i <= 4; i++ {
		rl.DrawTexture(Cursors[i-1], initPos+(Cursors[i-1].Width+8)*i, TitleBarSize[1]+(60-Cursors[i-1].Height)/2, rl.White)
	}
}

func AppContent() {
	rl.DrawRectangle(0, TitleBarSize[1], Size[0], 60, rl.Color{R: 100, G: 100, B: 100, A: 255})
	rl.DrawTextureEx(Icons[0], rl.Vector2{X: 10, Y: float32(TitleBarSize[1] + 5)}, 0, 1, rl.White)
	rl.DrawTextureEx(Icons[1], rl.Vector2{X: float32(10 + Icons[0].Width + 10), Y: float32(TitleBarSize[1] + 5)}, 0, 1, rl.White)

	for i := 1; i <= 10; i++ {
		DrawColorBox(int32(i))
	}
	DrawToolsIcons()

	rl.DrawTexture(CanvasTexture, 0, TitleBarSize[1]+60, rl.White)
}

var CanvasSize = [2]int32{Size[0], Size[1] - (TitleBarSize[1] + 60)}
var Canvas *rl.Image
var CanvasTexture rl.Texture2D
var DrawMode = 0
var DrawColor = rl.Red

func AppHandle() {
	if rl.CheckCollisionPointRec(mouse, rl.Rectangle{
		X:      10,
		Y:      float32(TitleBarSize[1] + 5),
		Width:  float32(Icons[0].Width),
		Height: float32(Icons[0].Height),
	}) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		rl.ExportImage(*Canvas, fmt.Sprintf("%d.png", time.Now().UnixNano()))
	}
}

func MakeCanvas() {
	Canvas = rl.GenImageColor(int(Size[0]), int(Size[1]), rl.White)
}

var Cursors [4]rl.Texture2D

func CursorHandler() {
	if DrawMode >= 4 {
		DrawMode = 3
	}
	if !rl.CheckCollisionPointRec(mouse, rl.Rectangle{Width: float32(TitleBarSize[0]), Height: float32(TitleBarSize[1])}) {
		rl.HideCursor()
		//if (mouse.X != 0 && int32(mouse.X) != Size[0]-1) && (mouse.Y != 0 && int32(mouse.Y) != Size[1]-1) {
		var col = rl.White
		if DrawMode != 3 && DrawMode != 1 {
			col = DrawColor
			col.A = 255
		}
		rl.DrawTexture(Cursors[DrawMode], int32(mouse.X)-Cursors[DrawMode].Width/2, int32(mouse.Y)-Cursors[DrawMode].Height/2, col)
		//}
	} else {
		rl.ShowCursor()
	}
}

func ToolPicker() {
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		var i int32
		for i = 1; i <= 10; i++ {
			var mov = i * (10 + 40)
			if rl.CheckCollisionPointRec(mouse, rl.Rectangle{X: float32(Size[0] - mov), Y: float32(TitleBarSize[1] + 10), Width: 40, Height: 40}) {
				DrawColor = ColorsPalette[0][i-1]
			}
		}
	}
	var initPos = 80 + Icons[0].Width + 10
	var i int32
	for i = 1; i <= 4; i++ {
		if rl.CheckCollisionPointRec(
			mouse,
			rl.Rectangle{
				X:      float32(initPos + (Cursors[i-1].Width+8)*i),
				Y:      float32(TitleBarSize[1] + (60-Cursors[i-1].Height)/2),
				Width:  float32(Cursors[i-1].Width),
				Height: float32(Cursors[i-1].Height),
			},
		) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			DrawMode = int(i - 1)
		}
	}
}

func PenToolHandler() {
	var clickX, clickY = ImagePointsOffset(mouse.X, mouse.Y)
	if rl.CheckCollisionPointRec(
		mouse,
		rl.Rectangle{Y: float32(TitleBarSize[1] + 60), Width: float32(CanvasSize[0]), Height: float32(CanvasSize[1])},
	) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			rl.ImageDrawRectangle(Canvas, clickX-10, clickY-10, 20, 20, DrawColor)
		} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
			rl.ImageDrawRectangle(Canvas, clickX-10, clickY-10, 20, 20, rl.White)
		}
	}
}

func EraserToolHandler() {
	var clickX, clickY = ImagePointsOffset(mouse.X, mouse.Y)
	if rl.CheckCollisionPointRec(
		mouse,
		rl.Rectangle{Y: float32(TitleBarSize[1] + 60), Width: float32(CanvasSize[0]), Height: float32(CanvasSize[1])},
	) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			rl.ImageDrawRectangle(Canvas, clickX-10, clickY-10, 20, 20, rl.White)
		}
	}
}

func FillToolHandler() {
	if rl.CheckCollisionPointRec(
		mouse,
		rl.Rectangle{Y: float32(TitleBarSize[1] + 60), Width: float32(CanvasSize[0]), Height: float32(CanvasSize[1])},
	) {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			rl.ImageDrawRectangle(Canvas, 0, 0, CanvasSize[0], CanvasSize[1], DrawColor)
		} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
			rl.ImageDrawRectangle(Canvas, 0, 0, CanvasSize[0], CanvasSize[1], rl.White)
		}
	}
}

func CursorToolHandler() {
	return
}

func ImagePointsOffset(posX float32, posY float32) (int32, int32) {
	posY -= float32(TitleBarSize[1] + 60)
	return int32(posX), int32(posY)
}
