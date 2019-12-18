package db

import (
	"fmt"
	mydb "github.com/yuzhikuan/filestore-server/db/mysql"
)

type User struct {
	Username string
	Email string
	Phone string
	SignupAt string
	LastActiveAt string
	Status int
}

// 通过用户名及密码完成user表的注册操作
func UserSignup(username string, passwd string) bool {
	sql := "insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// 判断密码是否一致
func UserSignin(username string, encpwd string) bool {
	sql := "select * from tbl_user where user_name=? limit 1"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// 更新用户登录的token
func UpdateToken(username string, token string) bool {
	sql := "replace into tbl_user_token (`user_name`,`user_token`) values (?,?)"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// 查询用户信息
func GetUserInfo(username string) (User, error) {
	user := User{}
	sql := "select user_name,signup_at from tbl_user where user_name=? limit 1"
	stmt, err := mydb.DBConn().Prepare(sql)
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		fmt.Println(err.Error())
	}
	return user, err
}