package store

import (
	"fmt"
	"log-record/utils"
	"time"
)

func init() {
	hasTable := DB.Migrator().HasTable(&RemoteLog{})
	if !hasTable {
		err := DB.Migrator().CreateTable(&RemoteLog{})
		if err != nil {
			fmt.Println(err)
		}
	}
	hasTable = DB.Migrator().HasTable(&User{})

	if !hasTable {
		err := DB.Migrator().CreateTable(&User{})
		if err != nil {
			fmt.Println(err)
		}
	}
}

type RemoteLog struct {
	// ID 自增ID
	Id uint64 `gorm:"primaryKey;type:int(20) AUTO_INCREMENT;comment:主键" json:"id"`
	// Application 应用标识
	Application string `gorm:"index;type:varchar(20);comment:应用标识" json:"application"`
	// Timestamp 时间戳，日志创建时间
	Timestamp int64 `gorm:"type:bigint(22);comment:时间戳，日志创建时间" json:"timestamp"`
	// LoggerName 日志名
	LoggerName string `gorm:"column:logger_name;type:varchar(100);comment:日志名" json:"loggerName"`
	// Level 日志级别
	Level string `gorm:"index;type:varchar(10);comment:日志级别" json:"level"`
	// Message 日志消息
	Message string `gorm:"type:text;comment:日志消息" json:"message"`
	//Address 客户端ip
	Address string `gorm:"type:varchar(20);comment:客户端ip" json:"address"`
	// Thread 线程
	Thread string `gorm:"type:varchar(40);comment:线程" json:"thread"`
	// LinkId 链路id，有则存在
	LinkId string `gorm:"column:link_id;index;type:varchar(40);default:'';comment:链路id，有则存在" json:"linkId"`
	// LinkCounter 链路自增器，保证链路有序
	LinkCounter string `gorm:"column:link_counter;type:int(5);default:0;comment:链路自增器，保证链路有序" json:"linkCounter"`
	// Metadata 其他元数据信息
	Metadata string `gorm:"type:text;comment:其他元数据信息" json:"metadata"`
}

func (RemoteLog) TableName() string {
	return "t_remote_log"
}

type User struct {
	// ID 自增ID
	Id         uint64    `gorm:"primaryKey;type:int(20) AUTO_INCREMENT;comment:'主键'" json:"id"`
	Username   string    `gorm:"type:varchar(20) NOT NULL;comment:用户名" json:"username"`
	Password   string    `gorm:"type:varchar(40) NOT NULL;comment:密码" json:"password"`
	CreateTime time.Time `gorm:"type:datetime;comment:创建时间" json:"create_time"`
}

func (User) TableName() string {
	return "t_user"
}

func (u *User) Auth(username, password string) bool {
	var count int64
	DB.Model(u).Where("username=? AND password=?", username, utils.Md5(password)).Count(&count)
	return count > 0
}
