package main

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MYSQL"
	"github.com/astaxie/beedb-master"
	"log"
	"time"
)

/**
beedb针对驼峰命名会自动帮你转化成下划线字段
例如你定义了Struct名字为UserInfo，
那么转化成 底层实现的时候是user_info，
字段命名也遵循该规则。
*/
type UserInfo struct {
	// 如果表的主键不是id，那么需要加上pk注释，显式的说这个字段是主键
	Uid        int    `PK` //`beedb:"PK" sql:"UID" tname:"user_info"`
	UserName   string `sql:"user_name"`
	DepartName string `sql:"depart_name"`
	Created    string `sql:"created"`
}

func main() {

	db, err := sql.Open("mysql", "root:root@/test?charset=utf8")

	//beedb的New函数实际上应该有两个参数，
	// 第一个参数标准接口的db，
	// 第二个参数是使用的数据库引擎，如果使用的 数据库引擎是MySQL/Sqlite,那么第二个参数都可以省略
	orm := beedb.New(db)
	//orm := beedb.New(db,"mssql")   // SQLServer
	//orm := beedb.New(db,"pg")   	 // PostgreSQL
	//实现调试
	beedb.OnDebug = true

	//新增数据
	var saveUser UserInfo
	saveUser.UserName = "安琪拉"
	saveUser.DepartName = "研发部"
	saveUser.Created = time.Now().Format("2006-01-02 03:04:05")
	err = orm.Save(&saveUser)
	if err != nil {
		log.Printf("Save Failure : %v", err)
	}

	//map数据插入
	addUser := make(map[string]interface{})
	addUser["user_name"] = "胡安琪"
	addUser["depart_name"] = "研发部"
	addUser["created"] = "1993-11-23"
	_, err = orm.SetTable("user_info").Insert(addUser)
	if err != nil {
		log.Printf("Insert Failure : %v", err)
	}

	//批量插入
	addOne := make(map[string]interface{})
	addTwo := make(map[string]interface{})
	addOne["user_name"] = "仙姑"
	addOne["depart_name"] = "云平台开发"
	addOne["created"] = "2012-12-02"
	addTwo["user_name"] = "韩湘子"
	addTwo["depart_name"] = "云平台开发"
	addTwo["created"] = "2012-12-02"

	addSlice := make([]map[string]interface{}, 0)
	addSlice = append(addSlice, addOne, addTwo)

	_, err = orm.SetTable("user_info").InsertBatch(addSlice)
	if err != nil {
		log.Printf("InsertBatch Failure : %v", err)
	}

	//更新数据也支持直接使用map操作
	t := make(map[string]interface{})
	t["user_name"] = "丁凌"
	_, err = orm.SetTable("user_info").SetPK("uid").Where(2).Update(t)
	if err != nil {
		log.Printf("Update Failure : %v", err)
	}
	//数据查询
	var user UserInfo
	err = orm.Where("uid=?", 2).Find(&user)
	if err != nil {
		log.Printf("Find By Primary Key Failure : %v", err)
	}

	var user2 UserInfo
	//可以省略主键,主键必须是 字段  `id`
	//所以这一句一定会执行报错，是一个坑
	err = orm.Where(2).Find(&user2)
	if err != nil {
		log.Printf("Ignore PrimaryKey Find Failure : %v", err)
	}

	var user3 UserInfo
	err = orm.Where("user_name = ?", "丁凌").Find(&user3)
	if err != nil {
		log.Printf("Find By UserName Failure : %v", err)
	}

	var user4 UserInfo
	err = orm.Where("user_name = ? and created >= ?", "丁凌", "1993-11-23").Find(&user4)
	if err != nil {
		log.Printf("Find By UserName And created Failure : %v", err)
	}

	var allUser []UserInfo

	err = orm.Where("uid > ?", 3).FindAll(&allUser)
	if err != nil {
		log.Printf("Find By Uid Failure : %v", err)
	}

	err = orm.Where("uid > ?", 1).Limit(0, 2).FindAll(&allUser)
	if err != nil {
		log.Printf("Find By Uid And Limit 0, 2 Failure : %v", err)
	}

	err = orm.Where("uid > ?", 1).Limit(2).FindAll(&allUser)
	if err != nil {
		log.Printf("Find By Uid And Limit Start 2 Failure : %v", err)
	}

	err = orm.OrderBy("uid desc,user_name desc").FindAll(&allUser)
	if err != nil {
		log.Printf("Find Order By uid desc,user_name desc Failure : %v", err)
	}
	//util.CheckError(err)

	//FindMap()函数返回的是[]map[string][]byte类型，所以你需要自己作类型转换。
	a, err := orm.SetTable("user_info").SetPK("uid").Where(2).Select("user_name,depart_name").FindMap()
	if err != nil {
		log.Printf("Find user_name,depart_name By uid Failure : %v", err)
	}
	for key, value := range a {
		fmt.Println(key)
		fmt.Println(value)
	}

	//查询出来，再删除
	//orm.Delete(&saveUser)
	//全部删除 先查出来
	//orm.DeleteAll(&allUser)
	//根据SQL删除
	//orm.SetTable("user_info").Where("uid>?", 3).DeleteRow()

}
