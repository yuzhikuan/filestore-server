package db

import (
	"fmt"
	mydb "github.com/yuzhikuan/filestore-server/db/mysql"
	"time"
)

// UserFile 用户文件表结构体
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

// OnUserFileUploadFinished 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	sql := "insert ignore tbl_user_file (`user_name`,`file_sha1`,`file_name`,`file_size`,`upload_at`) values (?,?,?,?,?)"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		return false
	}
	return true
}

// QueryUserFileMetas 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	sql := "select file_sha1,file_name,file_size,upload_at,last_update from tbl_user_file where user_name=? limit ?"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}
