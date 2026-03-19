package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"laogong-visa/internal/model"
	"laogong-visa/internal/repository"
)

// RecordService 记录服务
type RecordService struct {
	imageService *ImageService
}

// NewRecordService 创建记录服务
func NewRecordService() *RecordService {
	return &RecordService{
		imageService: NewImageService(),
	}
}

// Create 创建记录
func (s *RecordService) Create(req *model.CreateRecordRequest, userID uint64) (*model.Record, error) {
	imagesJSON, _ := json.Marshal(req.Images)

	record := &model.Record{
		Name:                req.Name,
		Position:            req.Position,
		Gender:              req.Gender,
		Title:               req.Title,
		Content:             req.Content,
		Images:              string(imagesJSON),
		CreatedBy:           userID,
		// Foreigner Profile
		DateOfBirth:         req.DateOfBirth,
		CountryOfOrigin:     req.CountryOfOrigin,
		PassportNo:          req.PassportNo,
		PassportIssuedDate:  req.PassportIssuedDate,
		PassportExpiredDate: req.PassportExpiredDate,
		EducationBackground: req.EducationBackground,
		VisaEntryDate:       req.VisaEntryDate,
		CardIssuedDate:      req.CardIssuedDate,
		CardExpiredDate:     req.CardExpiredDate,
		// Working History
		WorkingSession:      req.WorkingSession,
		CompanyName:         req.CompanyName,
		StartWorkingDate:    req.StartWorkingDate,
		StopWorkingDate:     req.StopWorkingDate,
	}

	if err := repository.DB.Create(record).Error; err != nil {
		return nil, err
	}

	return record, nil
}

// Update 更新记录
func (s *RecordService) Update(id uint64, req *model.UpdateRecordRequest) (*model.Record, error) {
	var record model.Record
	if err := repository.DB.First(&record, id).Error; err != nil {
		return nil, errors.New("record not found")
	}

	if req.Name != "" {
		record.Name = req.Name
	}
	if req.Position != "" {
		record.Position = req.Position
	}
	if req.Gender != "" {
		record.Gender = req.Gender
	}
	if req.Title != "" {
		record.Title = req.Title
	}
	if req.Content != "" {
		record.Content = req.Content
	}
	if req.Images != nil {
		imagesJSON, _ := json.Marshal(req.Images)
		record.Images = string(imagesJSON)
	}
	// Foreigner Profile
	if req.DateOfBirth != "" {
		record.DateOfBirth = req.DateOfBirth
	}
	if req.CountryOfOrigin != "" {
		record.CountryOfOrigin = req.CountryOfOrigin
	}
	if req.PassportNo != "" {
		record.PassportNo = req.PassportNo
	}
	if req.PassportIssuedDate != "" {
		record.PassportIssuedDate = req.PassportIssuedDate
	}
	if req.PassportExpiredDate != "" {
		record.PassportExpiredDate = req.PassportExpiredDate
	}
	if req.EducationBackground != "" {
		record.EducationBackground = req.EducationBackground
	}
	if req.VisaEntryDate != "" {
		record.VisaEntryDate = req.VisaEntryDate
	}
	if req.CardIssuedDate != "" {
		record.CardIssuedDate = req.CardIssuedDate
	}
	if req.CardExpiredDate != "" {
		record.CardExpiredDate = req.CardExpiredDate
	}
	// Working History
	if req.WorkingSession != "" {
		record.WorkingSession = req.WorkingSession
	}
	if req.CompanyName != "" {
		record.CompanyName = req.CompanyName
	}
	if req.StartWorkingDate != "" {
		record.StartWorkingDate = req.StartWorkingDate
	}
	if req.StopWorkingDate != "" {
		record.StopWorkingDate = req.StopWorkingDate
	}

	if err := repository.DB.Save(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

// Delete 删除记录
func (s *RecordService) Delete(id uint64) error {
	return repository.DB.Delete(&model.Record{}, id).Error
}

// GetByID 根据ID获取记录
func (s *RecordService) GetByID(id uint64) (*model.RecordResponse, error) {
	var record model.Record
	if err := repository.DB.Preload("CreatedByAdmin").First(&record, id).Error; err != nil {
		return nil, errors.New("record not found")
	}

	return s.toResponse(&record), nil
}

// GetAll 获取所有记录（前台使用）
func (s *RecordService) GetAll() ([]model.RecordResponse, error) {
	var records []model.Record
	if err := repository.DB.Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}

	var responses []model.RecordResponse
	for _, record := range records {
		responses = append(responses, *s.toResponse(&record))
	}

	return responses, nil
}

// GetAllForAdmin 获取所有记录（后台使用，带分页）
func (s *RecordService) GetAllForAdmin(page, pageSize int) ([]model.RecordResponse, int64, error) {
	var records []model.Record
	var total int64

	repository.DB.Model(&model.Record{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := repository.DB.Preload("CreatedByAdmin").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	var responses []model.RecordResponse
	for _, record := range records {
		responses = append(responses, *s.toResponse(&record))
	}

	return responses, total, nil
}

// RenderRecordImage 渲染记录为图片
func (s *RecordService) RenderRecordImage(id uint64) (string, error) {
	var record model.Record
	if err := repository.DB.First(&record, id).Error; err != nil {
		return "", errors.New("record not found")
	}

	var images []string
	json.Unmarshal([]byte(record.Images), &images)

	imageURL, err := s.imageService.RenderRecordToImage(record.Title, record.Content, images)
	if err != nil {
		return "", err
	}

	return imageURL, nil
}

// toResponse 转换为响应格式
func (s *RecordService) toResponse(record *model.Record) *model.RecordResponse {
	var images []string
	json.Unmarshal([]byte(record.Images), &images)

	return &model.RecordResponse{
		ID:                   record.ID,
		Name:                 record.Name,
		Position:             record.Position,
		Gender:               record.Gender,
		Title:                record.Title,
		Content:              record.Content,
		Images:               images,
		ImageURL:             fmt.Sprintf("/api/records/%d/image", record.ID),
		// Foreigner Profile
		DateOfBirth:          record.DateOfBirth,
		CountryOfOrigin:      record.CountryOfOrigin,
		PassportNo:           record.PassportNo,
		PassportIssuedDate:   record.PassportIssuedDate,
		PassportExpiredDate:  record.PassportExpiredDate,
		EducationBackground:  record.EducationBackground,
		VisaEntryDate:        record.VisaEntryDate,
		CardIssuedDate:       record.CardIssuedDate,
		CardExpiredDate:      record.CardExpiredDate,
		// Working History
		WorkingSession:       record.WorkingSession,
		CompanyName:          record.CompanyName,
		StartWorkingDate:     record.StartWorkingDate,
		StopWorkingDate:      record.StopWorkingDate,
		CreatedAt:            record.CreatedAt,
	}
}
