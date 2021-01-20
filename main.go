package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lesson25/middleware"
	"net/http"
	"os"
)

// 小项目 todo

/*
1. GET
2. POST 增
3. PUT  改
3. DELETE 删
 */

// 数据库模型
type TODO struct {
	Id uint `json:"id"`
	Title string `json:"title" binding:"required"`
	Status bool `json:"status"`
}

// 数据库链接初始化 全局设置db
var  db *gorm.DB
func initMYSQL() (err error){
	// 数据库mysql
	dsn := "root:p@ss1234@anji@tcp(10.108.26.60:3307)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		os.Exit(0)
		return
	} else {
		// 设置参数
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		return nil
	}

}

func todoGet(c *gin.Context) {
	var todoList []TODO
	if err := db.Find(&todoList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		// 自动结构化成 json对对象
		c.JSON(http.StatusOK, todoList)
	}
}

// 增加待办项
func todoPost(c *gin.Context) {
	var todo TODO
	if err :=c.ShouldBindJSON(&todo); err != nil {
		fmt.Printf("todo Post 方法, 参数有误, error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg": "参数绑定错误",
		})
		return
	}

	if err := db.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg": todo,
		})
	}

}

func todoPut(c *gin.Context) {
	id := c.Param("id")
	var todo TODO
	// 么有找到会有error: record not found
	err := db.Where("id = ?", id).First(&todo).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// 将传递的参数进行保存
	c.BindJSON(&todo)
	if err := db.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": todo})
	}

}


// 删除TODO LIST
func todoDelete(c *gin.Context) {
	id := c.Param("id")
	err := db.Where("id = ?", id).Delete(&TODO{}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}else{
		c.JSON(http.StatusOK, gin.H{id:"deleted"})
	}

}


func main() {
	r := gin.Default()
	r.Use(middleware.Cors())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "基于 gin Gorm的小项目 TODO待办事项")
	})
	// 定义路由组
	v1Group := r.Group("/v1")
	{
		v1Group.GET("/todo", todoGet)
		v1Group.POST("/todo", todoPost)
		v1Group.PUT("/todo/:id", todoPut)
		v1Group.DELETE("/todo/:id", todoDelete)
	}

	initMYSQL()

	// 表迁移
	db.AutoMigrate(&TODO{})


	// http server启动
	r.Run(":8080")
}
