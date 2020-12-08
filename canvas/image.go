package canvas

import (
	"fmt"
	"image"
	"path"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

var fontCacheMap = map[string]*font.Face{}
var imageCacheMap = map[string]*image.Image{}

// LoadLocalImage 本地图片缓存在内存中
func LoadLocalImage(imagePath string) (image.Image, error) {
	if img, cached := imageCacheMap[imagePath]; cached {
		return *img, nil
	}

	img, err := gg.LoadImage(imagePath)
	if err != nil {
		return nil, err
	}

	imageCacheMap[imagePath] = &img
	return img, nil
}

// LoadLocalFont 字体缓存在内存中
func LoadLocalFont(fontPath string, fontSize float64) (font.Face, error) {
	fontKey := fmt.Sprintf("%s:%f", fontPath, fontSize)
	if font, cached := fontCacheMap[fontKey]; cached {
		return *font, nil
	}

	fontPath = path.Join(fontPath)
	font, err := gg.LoadFontFace(fontPath, fontSize)
	if err != nil {
		return nil, err
	}

	fontCacheMap[fontKey] = &font
	return font, nil
}

// ImageResize 图片缩放
func ImageResize(input image.Image, width uint, height uint) image.Image {
	_ = input.Bounds()
	// TODO

	return input
}

// ImageRound 图片变圆
func ImageRound(input image.Image) image.Image {
	size := (input.Bounds().Dx() + input.Bounds().Dy()) / 2
	ctx := gg.NewContext(size, size)
	ctx.DrawRoundedRectangle(0, 0, float64(size), float64(size), float64(size/2))
	ctx.Clip()
	ctx.DrawImage(input, 0, 0)
	return ctx.Image()
}
