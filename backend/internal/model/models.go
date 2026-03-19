package model

import (
	"time"
	"gorm.io/gorm"
)

// Admin 管理员表
type Admin struct {
	ID        uint64         `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"` // 不返回密码
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Record 劳工签证记录表
type Record struct {
	ID                   uint64         `gorm:"primarykey" json:"id"`
	Name                 string         `gorm:"size:100" json:"name"`                           // 姓名
	Position             string         `gorm:"size:200" json:"position"`                       // 职位
	Gender               string         `gorm:"size:20" json:"gender"`                          // 性别
	Title                string         `gorm:"size:200" json:"title"`
	Content              string         `gorm:"type:text" json:"content"`                       // 文本内容
	Images               string         `gorm:"type:text" json:"images"`                        // JSON 数组存储图片URL
	// Foreigner Profile
	DateOfBirth          string         `gorm:"size:50" json:"date_of_birth"`                   // 出生日期
	CountryOfOrigin      string         `gorm:"size:100" json:"country_of_origin"`              // 国籍
	PassportNo           string         `gorm:"size:50" json:"passport_no"`                     // 护照号
	PassportIssuedDate   string         `gorm:"size:50" json:"passport_issued_date"`            // 护照签发日期
	PassportExpiredDate  string         `gorm:"size:50" json:"passport_expired_date"`           // 护照过期日期
	EducationBackground  string         `gorm:"size:100" json:"education_background"`           // 教育背景
	VisaEntryDate        string         `gorm:"size:50" json:"visa_entry_date"`                 // 签证入境日期
	CardIssuedDate       string         `gorm:"size:50" json:"card_issued_date"`                // 最新卡片签发日期
	CardExpiredDate      string         `gorm:"size:50" json:"card_expired_date"`               // 最新卡片过期日期
	// Working History
	WorkingSession       string         `gorm:"size:50" json:"working_session"`                 // 工作年度
	CompanyName          string         `gorm:"size:200" json:"company_name"`                   // 公司名称
	StartWorkingDate     string         `gorm:"size:50" json:"start_working_date"`              // 开始工作日期
	StopWorkingDate      string         `gorm:"size:50" json:"stop_working_date"`               // 停止工作日期
	CreatedBy            uint64         `json:"created_by"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	Admin Admin  `json:"admin"`
}

// CreateRecordRequest 创建记录请求
type CreateRecordRequest struct {
	Name                string   `json:"name"`
	Position            string   `json:"position"`
	Gender              string   `json:"gender"`
	Title               string   `json:"title"`
	Content             string   `json:"content"`
	Images              []string `json:"images"`
	// Foreigner Profile
	DateOfBirth         string   `json:"date_of_birth"`
	CountryOfOrigin     string   `json:"country_of_origin"`
	PassportNo          string   `json:"passport_no"`
	PassportIssuedDate  string   `json:"passport_issued_date"`
	PassportExpiredDate string   `json:"passport_expired_date"`
	EducationBackground string   `json:"education_background"`
	VisaEntryDate       string   `json:"visa_entry_date"`
	CardIssuedDate      string   `json:"card_issued_date"`
	CardExpiredDate     string   `json:"card_expired_date"`
	// Working History
	WorkingSession      string   `json:"working_session"`
	CompanyName         string   `json:"company_name"`
	StartWorkingDate    string   `json:"start_working_date"`
	StopWorkingDate     string   `json:"stop_working_date"`
}

// UpdateRecordRequest 更新记录请求
type UpdateRecordRequest struct {
	Name                string   `json:"name"`
	Position            string   `json:"position"`
	Gender              string   `json:"gender"`
	Title               string   `json:"title"`
	Content             string   `json:"content"`
	Images              []string `json:"images"`
	// Foreigner Profile
	DateOfBirth         string   `json:"date_of_birth"`
	CountryOfOrigin     string   `json:"country_of_origin"`
	PassportNo          string   `json:"passport_no"`
	PassportIssuedDate  string   `json:"passport_issued_date"`
	PassportExpiredDate string   `json:"passport_expired_date"`
	EducationBackground string   `json:"education_background"`
	VisaEntryDate       string   `json:"visa_entry_date"`
	CardIssuedDate      string   `json:"card_issued_date"`
	CardExpiredDate     string   `json:"card_expired_date"`
	// Working History
	WorkingSession      string   `json:"working_session"`
	CompanyName         string   `json:"company_name"`
	StartWorkingDate    string   `json:"start_working_date"`
	StopWorkingDate     string   `json:"stop_working_date"`
}

// RecordResponse 记录响应（包含图片渲染URL）
type RecordResponse struct {
	ID                   uint64    `json:"id"`
	Name                 string    `json:"name"`
	Position             string    `json:"position"`
	Gender               string    `json:"gender"`
	Title                string    `json:"title"`
	Content              string    `json:"content"`
	Images               []string  `json:"images"`
	ImageURL             string    `json:"image_url"`  // 渲染后的图片URL
	// Foreigner Profile
	DateOfBirth          string    `json:"date_of_birth"`
	CountryOfOrigin      string    `json:"country_of_origin"`
	PassportNo           string    `json:"passport_no"`
	PassportIssuedDate   string    `json:"passport_issued_date"`
	PassportExpiredDate  string    `json:"passport_expired_date"`
	EducationBackground  string    `json:"education_background"`
	VisaEntryDate        string    `json:"visa_entry_date"`
	CardIssuedDate       string    `json:"card_issued_date"`
	CardExpiredDate      string    `json:"card_expired_date"`
	// Working History
	WorkingSession       string    `json:"working_session"`
	CompanyName          string    `json:"company_name"`
	StartWorkingDate     string    `json:"start_working_date"`
	StopWorkingDate      string    `json:"stop_working_date"`
	CreatedAt            time.Time `json:"created_at"`
}
