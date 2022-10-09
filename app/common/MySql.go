package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MySql struct {
	source       string  // 数据库源
	driver       string  // 数据库驱动
	fields       string  // 字段
	tableName    string  // 表名
	whereStr     string  // where语句
	limitNumber  string  // 限制条数
	offsetNumber string  // 限制条数
	orderBy      string  // 排序条件
	execStr      string  // 执行sql语句
	conn         *sql.DB // 数据库连接
}

var DbPool *sync.Pool
var GlobalDb *sql.DB

// 初始化连接池
func init() {
	MySql := MySql{}
	source, _ := ConfigObj.GetString("db", "source")
	driver, _ := ConfigObj.GetString("db", "driver")
	db, err := sql.Open(driver, source)
	db.SetMaxOpenConns(2000)             // 最大链接
	db.SetMaxIdleConns(1000)             // 空闲连接，也就是连接池里面的数量
	db.SetConnMaxLifetime(7 * time.Hour) // 设置最大生成周期是7个小时
	MySql.checkErr(err)
	GlobalDb = db
}

func (MySql MySql) GetConn() *MySql {
	MySql.conn = GlobalDb
	return &MySql
}

func (MySql *MySql) Close() error {
	err := MySql.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (MySql *MySql) Select(tableName string, field []string) *MySql {
	var allField string
	allField = strings.Join(field, ",")
	MySql.fields = "select " + allField + " from " + tableName
	MySql.tableName = tableName
	return MySql
}

func (MySql *MySql) Where(cond map[string]string) *MySql {
	var whereStr = ""
	if len(cond) != 0 {
		whereStr = " where "
		for key, value := range cond {
			if strings.Contains(key, "like") || strings.Contains(key, "IN") {
				whereStr += key + value + " AND "
				continue
			}
			if !strings.Contains(key, "=") && !strings.Contains(key, ">") && !strings.Contains(key, "<") {
				key += "="
			}
			whereStr += key + "'" + value + "'" + " AND "
		}
	}
	// 删除所有字段最后一个,
	whereStr = strings.TrimSuffix(whereStr, "AND ")
	MySql.whereStr = whereStr
	return MySql
}

func (MySql *MySql) Limit(number int) *MySql {
	MySql.limitNumber = " limit " + strconv.Itoa(number)
	return MySql
}
func (MySql *MySql) Offset(number int) *MySql {
	MySql.offsetNumber = " OFFSET " + strconv.Itoa(number)
	return MySql
}

func (MySql *MySql) OrderByString(orderString ...string) *MySql {
	if len(orderString) > 2 || len(orderString) <= 0 {
		log.Fatal("传入参数错误")
	} else if len(orderString) == 1 {
		MySql.orderBy = " ORDER BY " + orderString[0] + " ASC"
	} else {
		MySql.orderBy = " ORDER BY " + orderString[0] + " " + orderString[1]
	}
	return MySql
}

func (MySql MySql) Insert(tableName string, data map[string]string) int64 {
	var allField = ""
	var allValue = ""
	var allTrueValue []interface{}
	if len(data) != 0 {
		for key, value := range data {
			allField += key + ","
			allValue += "?" + ","
			allTrueValue = append(allTrueValue, value)
		}
	}
	allValue = strings.TrimSuffix(allValue, ",")
	allField = strings.TrimSuffix(allField, ",")
	allValue = "(" + allValue + ")"
	allField = "(" + allField + ")"
	var theStr = "insert into " + tableName + " " + allField + " values " + allValue
	MySql.printSql(theStr)
	stmt, err := MySql.conn.Prepare(theStr)
	MySql.checkErr(err)
	res, err := stmt.Exec(allTrueValue...)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	MySql.checkErr(err)
	id, err := res.LastInsertId()
	return id
}

func (MySql MySql) Update(tableName string, str map[string]string) int64 {
	var tempStr = ""
	var allValue []interface{}
	for key, value := range str {
		tempStr += key + "=" + "?" + ","
		allValue = append(allValue, value)
	}
	tempStr = strings.TrimSuffix(tempStr, ",")
	MySql.execStr = "update " + tableName + " set " + tempStr
	var allStr = MySql.execStr + MySql.whereStr
	MySql.printSql(allStr)
	stmt, err := MySql.conn.Prepare(allStr)
	MySql.checkErr(err)
	res, err := stmt.Exec(allValue...)
	MySql.checkErr(err)
	rows, err := res.RowsAffected()
	return rows

}

func (MySql MySql) Delete(tableName string) int64 {
	var tempStr = ""
	tempStr = "delete from " + tableName + MySql.whereStr
	MySql.printSql(tempStr)
	stmt, err := MySql.conn.Prepare(tempStr)
	MySql.checkErr(err)
	res, err := stmt.Exec()
	MySql.checkErr(err)
	rows, err := res.RowsAffected()
	return rows
}

func (MySql MySql) QueryAll() []map[string]string {
	var queryStr = MySql.fields + MySql.whereStr + MySql.orderBy + MySql.limitNumber + MySql.offsetNumber
	return MySql.ExecSql(queryStr)
	/*
		rows, err := MySql.conn.Query(queryStr)
		defer rows.Close()
		MySql.checkErr(err)
		Column, err := rows.Columns()
		MySql.checkErr(err)
		// 创建一个查询字段类型的slice
		values := make([]sql.RawBytes, len(Column))
		// 创建一个任意字段类型的slice
		scanArgs := make([]interface{}, len(values))
		// 创建一个slice保存所以的字段
		var allRows []map[string]string
		for i := range values {
			// 把values每个参数的地址存入scanArgs
			scanArgs[i] = &values[i]
		}
		for rows.Next() {
			// 把存放字段的元素批量放进去
			err = rows.Scan(scanArgs...)
			MySql.checkErr(err)
			tempRow := make(map[string]string, len(Column))
			for i, col := range values {
				var key = Column[i]
				tempRow[key] = string(col)
			}
			allRows = append(allRows, tempRow)
		}
		return allRows
	*/
}

func (MySql MySql) ExecSql(queryStr string) []map[string]string {
	MySql.printSql(queryStr)
	rows, err := MySql.conn.Query(queryStr)
	defer rows.Close()
	MySql.checkErr(err)
	Column, err := rows.Columns()
	MySql.checkErr(err)
	// 创建一个查询字段类型的slice
	values := make([]sql.RawBytes, len(Column))
	// 创建一个任意字段类型的slice
	scanArgs := make([]interface{}, len(values))
	// 创建一个slice保存所以的字段
	var allRows []map[string]string
	for i := range values {
		// 把values每个参数的地址存入scanArgs
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		// 把存放字段的元素批量放进去
		err = rows.Scan(scanArgs...)
		MySql.checkErr(err)
		tempRow := make(map[string]string, len(Column))
		for i, col := range values {
			var key = Column[i]
			tempRow[key] = string(col)
		}
		allRows = append(allRows, tempRow)
	}
	return allRows
}

func (MySql MySql) QueryRow() map[string]string {
	var queryStr = MySql.fields + MySql.whereStr + MySql.orderBy + MySql.limitNumber + MySql.offsetNumber
	MySql.printSql(queryStr)
	result, err := MySql.conn.Query(queryStr)
	defer result.Close()
	MySql.checkErr(err)
	Column, err := result.Columns()
	// 创建一个查询字段类型的slice的键值对
	values := make([]sql.RawBytes, len(Column))
	// 创建一个任意字段类型的slice的键值对
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		// 把values每个参数的地址存入scanArgs
		scanArgs[i] = &values[i]
	}

	for result.Next() {
		err = result.Scan(scanArgs...)
		MySql.checkErr(err)
	}
	tempRow := make(map[string]string, len(Column))
	for i, col := range values {
		var key = Column[i]
		tempRow[key] = string(col)
	}
	return tempRow

}

func (MySql MySql) printSql(s string) {
	fmt.Println(s)
}

/**
检查错误
*/
func (MySql MySql) checkErr(err error) {
	if err != nil {
		log.Fatal("错误：", err)
	}
}
