package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/ovear/go-tumblr-crawler/config/global"
	"fmt"
	"database/sql"
	"github.com/ovear/go-tumblr-crawler/downloader/response"
	"github.com/ovear/go-tumblr-crawler/utils"
	"encoding/json"
)

var dsn string

func Init() bool {
	settings := global.Settings
	if (!settings.EnableDB) {
		return false
	}
	//{username}:{password}@{address}:{port}/{dbname}?charset=utf8mb4,utf8
	//id:password@tcp(your-amazonaws-uri.com:3306)/dbname
	dsn_format := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8"
	dsn = fmt.Sprintf(dsn_format, settings.DBUsername, settings.DBPassword, settings.DBHost, settings.DBPort,settings.DBDatabase)
	fmt.Printf("初始化数据库[%s]\n", dsn)

	return true
}

func Test() bool {
	if dsn == "" {
		Init()
	}
	db, err := sql.Open("mysql", dsn)
	_, err = db.Exec("select id from " + global.Settings.DBTable + " limit 1")
	if err != nil {
		fmt.Printf("测试数据库失败 [%s]\n", err)
	}
	defer db.Close()
	return err == nil
}

var dbc *sql.DB
func Connection() *sql.DB {
	if dsn == "" {
		Init()
	}
	if dbc == nil {
		d, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("获取数据库连接出错[%s]\n", err)
			return nil
		}
		dbc = d
	}
	return dbc
}

func Exist(rawURL string) (exist bool, error bool) {
	conn := Connection()
	sql := "select id from %s where url_raw = ? limit 1"
	sql = fmt.Sprintf(sql, global.Settings.DBTable)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		fmt.Printf("执行数据库指令出错[%s][%s]\n", sql, err)
		error = true
		return
	}
	res, err := stmt.Query(rawURL)
	if err != nil {
		fmt.Printf("执行数据库指令出错[%s][%s]\n", sql, err)
		error = true
		return
	}
	defer res.Close()
	defer stmt.Close()
	exist = res.Next()
	return
}

func Insert(post response.BasePost, u string) bool {
	rawURL := u
	if post.Type == "photo" {
		rawURL = utils.ConvertToRawTumblrAddress(u)
	}
	conn := Connection()
	sql_command := "insert into %s set blog=?, post_id=?, photo_caption=?, tags=?, type=?, reblog_key=?, url=?, url_raw=?"
	sql_command = fmt.Sprintf(sql_command, global.Settings.DBTable)
	exist, err := Exist(rawURL)
	if err {
		return false
	}
	if exist {
		fmt.Printf("[DB]已存在{blog=%s, post_id=%s}\n", post.Tumblelog.Name, post.Id)
		return true
	}
	fmt.Printf("[DB]插入到数据库{blog=%s, post_id=%s}\n", post.Tumblelog.Name, post.Id)

	stmt, e := conn.Prepare(sql_command)
	defer stmt.Close()
	if e != nil {
		fmt.Printf("执行数据库指令出错[%s][%s]\n", sql_command, e)
		return false
	}
	t := 1
	if post.Type == "video" {
		t = 2
	}
	tagsJson, _ := json.Marshal(post.Tag)
	tagJsonString := string(tagsJson)
	if tagJsonString == "null" {
		tagJsonString = "[]"
	}
	res, e := stmt.Exec(post.Tumblelog.Name, post.Id, post.PhotoCaption, tagJsonString, t, post.ReblogKey, u, rawURL)
	if e != nil {
		fmt.Printf("执行数据库指令出错[%s][%s]\n", sql_command, e)
		return false
	}
	_, e = res.LastInsertId()
	return e == nil
}
