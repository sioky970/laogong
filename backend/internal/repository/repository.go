package repository

import (
	"fmt"
	"os"

	"laogong-visa/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "laogong")
	password := getEnv("DB_PASSWORD", "laogong123")
	dbname := getEnv("DB_NAME", "laogong_visa")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	return DB.AutoMigrate(&model.Admin{}, &model.Record{})
}

// CreateDefaultAdmin 创建默认管理员
func CreateDefaultAdmin() error {
	var count int64
	DB.Model(&model.Admin{}).Count(&count)
	if count == 0 {
		// 默认管理员: admin / admin123
		hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		admin := &model.Admin{
			Username: "admin",
			Password: string(hashed),
		}
		return DB.Create(admin).Error
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
