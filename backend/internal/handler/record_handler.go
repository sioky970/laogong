package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"laogong-visa/internal/model"
	"laogong-visa/internal/service"
	"laogong-visa/middleware"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// RecordHandler 记录处理器
type RecordHandler struct {
	service *service.RecordService
}

// NewRecordHandler 创建记录处理器
func NewRecordHandler() *RecordHandler {
	return &RecordHandler{
		service: service.NewRecordService(),
	}
}

// Create 创建记录（需要认证）
func (h *RecordHandler) Create(c *gin.Context) {
	var req model.CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	record, err := h.service.Create(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// Update 更新记录（需要认证）
func (h *RecordHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req model.UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, err := h.service.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// Delete 删除记录（需要认证）
func (h *RecordHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// GetByID 根据ID获取记录（公开访问）
func (h *RecordHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	record, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// GetAll 获取所有记录（公开访问，前台使用）
func (h *RecordHandler) GetAll(c *gin.Context) {
	records, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// GetAllForAdmin 获取所有记录（需要认证，后台使用，带分页）
func (h *RecordHandler) GetAllForAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	records, total, err := h.service.GetAllForAdmin(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     records,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetRecordImage 获取记录渲染后的图片（公开访问）
func (h *RecordHandler) GetRecordImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	imageURL, err := h.service.RenderRecordImage(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image_url": imageURL})
}

// UploadImage 上传图片（需要认证）
func (h *RecordHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未上传文件"})
		return
	}
	defer file.Close()

	// 文件大小限制：最大 5MB
	const maxSize = 5 * 1024 * 1024 // 5MB
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制，最大允许 5MB"})
		return
	}

	// 格式验证：仅允许 image/jpeg、image/png、image/gif
	contentType := header.Header.Get("Content-Type")
	allowedTypes := map[string]string{
		"image/jpeg": "jpg",
		"image/png":  "png",
		"image/gif":  "gif",
	}
	ext, ok := allowedTypes[contentType]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件格式，仅允许 JPEG、PNG、GIF 格式"})
		return
	}

	// 读取文件内容
	buffer, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
		return
	}

	imageService := service.NewImageService()
	imageURL, err := imageService.SaveMultipartImage(buffer, ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": imageURL})
}

// GetQRCode 生成记录的二维码（公开访问）
func (h *RecordHandler) GetQRCode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	// 验证记录是否存在
	_, err = h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	// 获取前端 URL
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:9528"
	}

	// 生成二维码内容：前端门户页面 URL
	qrContent := fmt.Sprintf("%s/portal/%d", frontendURL, id)

	// 生成二维码 PNG 图片（256x256 像素）
	png, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成二维码失败"})
		return
	}

	// 返回 PNG 图片
	c.Data(http.StatusOK, "image/png", png)
}
