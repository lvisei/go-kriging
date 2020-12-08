package canvas

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/fogleman/gg"
)

// Canvas 画布
type Canvas struct {
	context *gg.Context
	Width   uint
	Height  uint
}

// TextConfig 绘制文字配置
type TextConfig struct {
	Text     string
	FontName string
	FontSize float64
	Color    color.Color
	OffsetX  float64
	OffsetY  float64
	AlignX   float64
	AlignY   float64
	MaxWidth float64
}

// NewCanvas 创建新的画布
func NewCanvas(width, height int) *Canvas {
	canvas := &Canvas{
		context: gg.NewContext(width, height),
		Width:   uint(width),
		Height:  uint(height),
	}

	return canvas
}

// NewCanvasWithImage 从图片创建新的画布
func NewCanvasWithImage(background image.Image) *Canvas {
	canvas := &Canvas{
		context: gg.NewContext(background.Bounds().Dx(), background.Bounds().Dy()),
		Width:   uint(background.Bounds().Dx()),
		Height:  uint(background.Bounds().Dy()),
	}

	canvas.context.DrawImage(background, 0, 0)
	return canvas
}

// NewCanvasWithLocal 从本地图片创建新的画布
func NewCanvasWithLocalImagePath(imagePath string) (*Canvas, error) {
	img, err := LoadLocalImage(imagePath)
	if err != nil {
		return nil, err
	}

	return NewCanvasWithImage(img), nil
}

// DrawImage 绘制图片
func (canvas *Canvas) DrawImage(img image.Image, offsetX int, offsetY int) {
	canvas.context.DrawImage(img, offsetX, offsetY)
}

// DrawLine 绘制横线
func (canvas *Canvas) DrawLine(color color.Color, offsetX, offsetY, length float64) {
	canvas.context.DrawLine(offsetX, offsetY, offsetX+length, offsetY)
	canvas.context.ClosePath()
	canvas.context.SetLineWidth(1)
	canvas.context.SetColor(color)
	canvas.context.StrokePreserve()
	canvas.context.Stroke()
}

// DrawRect 绘制矩形
func (canvas *Canvas) DrawRect(x, y, w, h float64, c color.Color) {
	canvas.context.SetColor(c)
	canvas.context.DrawRectangle(x, y, w, h)
	canvas.context.Fill()
}

// MeasureString 测量文字
func (canvas *Canvas) MeasureString(opt *TextConfig) float64 {
	font, err := LoadLocalFont(opt.FontName, opt.FontSize)
	if err != nil {
		return 0
	}

	canvas.context.SetFontFace(font)
	width, _ := canvas.context.MeasureString(opt.Text)
	return width
}

// DrawText 绘制图片
func (canvas *Canvas) DrawText(opt *TextConfig) error {
	font, err := LoadLocalFont(opt.FontName, opt.FontSize)
	if err != nil {
		return err
	}

	canvas.context.SetColor(opt.Color)

	canvas.context.SetFontFace(font)
	width := canvas.MeasureString(opt)
	if opt.MaxWidth > 0 && width > opt.MaxWidth {
		font, err = LoadLocalFont(opt.FontName, opt.FontSize*opt.MaxWidth/width)
		if err != nil {
			return err
		}
	}

	canvas.context.SetFontFace(font)
	canvas.context.DrawStringAnchored(opt.Text, opt.OffsetX, opt.OffsetY, opt.AlignX, opt.AlignY)
	return nil
}

// Output 输出 PNG 图片
func (canvas *Canvas) Output() ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, canvas.context.Image()); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// SavePNG
func (canvas *Canvas) SavePNG(path string) error {
	if err := canvas.context.SavePNG(path); err != nil {
		return err
	}

	return nil
}
