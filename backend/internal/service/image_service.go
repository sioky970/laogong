package service

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// ImageService 图片服务
type ImageService struct {
	UploadDir string
}

// NewImageService 创建图片服务
func NewImageService() *ImageService {
	uploadDir := "/app/uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}
	return &ImageService{UploadDir: uploadDir}
}

// RenderRecordToImage 将记录渲染为图片
func (s *ImageService) RenderRecordToImage(title, content string, imageURLs []string) (string, error) {
	// 创建画布 800x1200
	width := 800
	height := 1200
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 白色背景
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 绘制标题
	titleY := 60
	s.drawText(img, title, 40, titleY, width-80, true, color.Black)

	// 绘制分隔线
	for x := 40; x < width-40; x++ {
		img.Set(x, titleY+40, color.Black)
	}

	// 绘制内容
	contentY := titleY + 80
	s.drawText(img, content, 40, contentY, width-80, false, color.Black)

	// 绘制上传的图片
	if len(imageURLs) > 0 {
		currentY := contentY + 200
		for _, url := range imageURLs {
			if currentY > height-150 {
				break // 避免超出画布
			}
			imgPath := filepath.Join(s.UploadDir, filepath.Base(url))
			if _, err := os.Stat(imgPath); err == nil {
				s.drawUploadedImage(img, imgPath, 40, currentY, width-80)
				currentY += 220
			}
		}
	}

	// 添加水印
	s.drawText(img, "劳工签证管理系统", width-250, height-40, 200, false, color.Gray{128})

	// 保存图片
	filename := fmt.Sprintf("record_%s.jpg", uuid.New().String())
	filepath := filepath.Join(s.UploadDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := jpeg.Encode(file, img, &jpeg.Options{Quality: 90}); err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}

// drawText 在图片上绘制文本（自动换行）
func (s *ImageService) drawText(img *image.RGBA, text string, x, y, maxWidth int, isBold bool, c color.Color) {
	face := basicfont.Face7x13
	if isBold {
		face = basicfont.Face7x13 // 基础字体，可以替换为更粗的字体
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: face,
	}

	lines := s.wrapText(text, maxWidth, face)
	lineHeight := 20

	for i, line := range lines {
		d.Dot = fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(y + i*lineHeight),
		}
		d.DrawString(line)
	}
}

// wrapText 自动换行
func (s *ImageService) wrapText(text string, maxWidth int, face font.Face) []string {
	var lines []string
	var currentLine string

	for _, word := range strings.Split(text, "") {
		testLine := currentLine + word
		bounds, _ := font.BoundString(face, testLine)
		width := (bounds.Max.X - bounds.Min.X).Ceil()

		if width > maxWidth && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// drawUploadedImage 绘制上传的图片到画布
func (s *ImageService) drawUploadedImage(dst *image.RGBA, imgPath string, x, y, maxWidth int) {
	file, err := os.Open(imgPath)
	if err != nil {
		return
	}
	defer file.Close()

	var src image.Image
	if strings.HasSuffix(strings.ToLower(imgPath), ".png") {
		src, _ = png.Decode(file)
	} else {
		src, _ = jpeg.Decode(file)
	}

	if src == nil {
		return
	}

	// 计算缩放比例
	bounds := src.Bounds()
	ratio := float64(maxWidth) / float64(bounds.Dx())
	if ratio > 1 {
		ratio = 1
	}

	newWidth := int(float64(bounds.Dx()) * ratio)
	newHeight := int(float64(bounds.Dy()) * ratio)

	// 绘制缩放后的图片
	for py := 0; py < newHeight; py++ {
		for px := 0; px < newWidth; px++ {
			sx := int(float64(px) / ratio)
			sy := int(float64(py) / ratio)
			if sx < bounds.Dx() && sy < bounds.Dy() {
				dst.Set(x+px, y+py, src.At(bounds.Min.X+sx, bounds.Min.Y+sy))
			}
		}
	}
}

// SaveBase64Image 保存 base64 编码的图片
func (s *ImageService) SaveBase64Image(base64Data string) (string, error) {
	// 解析 base64
	parts := strings.Split(base64Data, ",")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 data")
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	// 生成文件名
	ext := "jpg"
	if strings.Contains(parts[0], "png") {
		ext = "png"
	}
	filename := fmt.Sprintf("%s_%d.%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(s.UploadDir, filename)

	// 保存文件
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}

// SaveMultipartImage 保存上传的图片文件
func (s *ImageService) SaveMultipartImage(data []byte, ext string) (string, error) {
	filename := fmt.Sprintf("%s_%d.%s", uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(s.UploadDir, filename)

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}

// GetImageData 获取图片数据
func (s *ImageService) GetImageData(filename string) ([]byte, string, error) {
	filepath := filepath.Join(s.UploadDir, filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, "", err
	}

	contentType := "image/jpeg"
	if strings.HasSuffix(strings.ToLower(filename), ".png") {
		contentType = "image/png"
	}

	return data, contentType, nil
}
