package main

import (
	"MySQL_Job/Repository/DAO"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := "user:password@tcp(127.0.0.13316)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移模式
	err = db.AutoMigrate(&DAO.Department{}, &DAO.User{}, &DAO.Permission{}, &DAO.Job{}, &DAO.JobExecutionHistory{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database migration completed.")

}
