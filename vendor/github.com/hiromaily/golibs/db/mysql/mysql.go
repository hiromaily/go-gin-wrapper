package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	lg "github.com/hiromaily/golibs/log"
	"github.com/hiromaily/golibs/times"
	u "github.com/hiromaily/golibs/utils"
	"reflect"
	"time"
)

//TODO:トランザクションの機能もあるので、どこかに追加しておく
//TODO:異なるlibraryを使っているが、各funcのInterfaceを統一すればよいのでは？
//http://qiita.com/tenntenn/items/dddb13c15643454a7c3b
//http://go-database-sql.org/

const tFomt = "2006-01-02 15:04:05"

type MS struct {
	DB         *sql.DB
	Rows       *sql.Rows
	Err        error
	ServerInfo //embeded
}

type ServerInfo struct {
	host   string
	port   uint16
	dbname string
	user   string
	pass   string
}

var dbInfo MS

//-----------------------------------------------------------------------------
// Basic
//-----------------------------------------------------------------------------
func New(host, dbname, user, pass string, port uint16) {
	var err error
	if dbInfo.DB == nil {
		dbInfo.host = host
		dbInfo.port = port
		dbInfo.dbname = dbname
		dbInfo.user = user
		dbInfo.pass = pass

		dbInfo.DB, err = dbInfo.Connection()
	}
	//lg.Debugf("dbInfo.db %+v\n", *dbInfo.DB)
	if err != nil {
		panic(err.Error())
	}
	return
}

// singleton architecture
func GetDBInstance() *MS {
	var err error
	if dbInfo.DB == nil {
		//TODO: it may be better to call New()
		dbInfo.DB, err = dbInfo.Connection()
	}
	if err != nil {
		panic(err.Error())
	}
	return &dbInfo
}

func (ms *MS) getDsn() string {
	//If use nil on Date column, set *time.Time
	//Be careful when parsing is required on Date type
	// e.g. db, err := sql.Open("mysql", "root:@/?parseTime=true")
	param := "?charset=utf8&parseTime=True&loc=Local"
	//user:password@tcp(localhost:3306)/dbname?tls=skip-verify&autocommit=true
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s",
		ms.user, ms.pass, ms.host, ms.port, ms.dbname, param)
}

// Connection
// Be careful, sql.Open() doesn't return err. Use db.Ping() to check DB condition.
func (ms *MS) Connection() (*sql.DB, error) {
	//return sql.Open("mysql", getDsn())
	db, _ := sql.Open("mysql", ms.getDsn())
	return db, db.Ping()
}

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
func (ms *MS) SetMaxIdleConns(n int) {
	ms.DB.SetMaxIdleConns(n)
}

//SetMaxOpenConns sets the maximum number of open connections to the database.
func (ms *MS) SetMaxOpenConns(n int) {
	ms.DB.SetMaxOpenConns(n)
}

// Close
func (ms *MS) Close() {
	ms.DB.Close()
}

//-----------------------------------------------------------------------------
// Select
//-----------------------------------------------------------------------------
// SELECT Count: Get number of rows
func (ms *MS) SelectCount(countSql string, args ...interface{}) (int, error) {
	//field on table
	var count int

	//1. create sql and exec
	//err := self.db.QueryRow("SELECT count(user_id) FROM t_users WHERE delete_flg=?", "0").Scan(&count)
	err := ms.DB.QueryRow(countSql, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

//get Rows and return db instance
func (ms *MS) SelectIns(selectSQL string, args ...interface{}) *MS {
	defer times.Track(time.Now(), "SelectIns()")
	//SelectSQLAllFieldIns() took 471.577µs

	//If no args, set nil

	//1. create sql and exec
	//rows, err := self.db.Query("SELECT * FROM t_users WHERE delete_flg=?", "0")
	ms.Rows, ms.Err = ms.DB.Query(selectSQL, args...)
	if ms.Err != nil {
		lg.Errorf("SelectSQLAllFieldIns()->ms.DB.Query():error is %s, \n %s", ms.Err, selectSQL)
	}

	return ms
}

//set extracted data into parameter variable
func (ms *MS) ScanOne(x interface{}) bool {
	//defer times.Track(time.Now(), "ScanOne()")
	//ScanOne() took 5.23µs

	if ms.Err != nil {
		lg.Errorf("ScanOne(): ms.Err has error: %s", ms.Err)
		return false
	}

	//e.g)v = person Person
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		ms.Err = errors.New("parameter is not valid. it sould be pointer and not nil.")
		return false
	} else {
		if v.Elem().Kind() == reflect.Struct {

			//create container to set scaned record on database
			values, scanArgs := makeScanArgs(v.Elem().Type())

			//check len(value) and column
			validateStructAndColumns(ms, values)
			if ms.Err != nil {
				return false
			}

			// rows.Next()
			ret := ms.Rows.Next()
			if !ret {
				//ms.Err = errors.New("nodata")
				return false
			}

			// rows.Scan()
			ms.Err = ms.Rows.Scan(scanArgs...)
			if ms.Err != nil {
				return false
			}

			//ms.Err = ms.Rows.Scan(v)
			scanStruct(values, v.Elem())
		} else {
			ms.Err = errors.New("parameter should be pointer of struct slice or struct")
			return false
		}
	}
	return true
}

func (ms *MS) Scan(x interface{}) bool {
	//defer times.Track(time.Now(), "Scan()")
	//Scan() took 465.971µs

	if ms.Err != nil {
		lg.Errorf("Scan(): ms.Err has error: %s", ms.Err)
		return false
	}

	//e.g)v = persons []Person
	v := reflect.ValueOf(x)

	if v.Kind() != reflect.Ptr || v.IsNil() {
		ms.Err = errors.New("parameter is not valid. it sould be pointer and not nil.")
		return false
	} else {
		if v.Elem().Kind() == reflect.Slice || v.Elem().Kind() == reflect.Array {
			elemType := v.Elem().Type().Elem() //reflects_test.TeacherInfo
			newElem := reflect.New(elemType).Elem()

			//create container to set scaned record on database
			values, scanArgs := makeScanArgs(newElem.Type())

			//check len(value) and column
			validateStructAndColumns(ms, values)
			if ms.Err != nil {
				return false
			}

			// rows.Next()
			cnt := 0
			for ms.Rows.Next() {
				ms.Err = ms.Rows.Scan(scanArgs...)
				if ms.Err != nil {
					return false
				}

				scanStruct(values, newElem)
				v.Elem().Set(reflect.Append(v.Elem(), newElem))
				cnt++
			}
			if cnt == 0 {
				return false
			}
		} else {
			ms.Err = errors.New("parameter is not valid. it sould be pointer and not nil.")
			return false
		}
	}

	return true
}

func makeScanArgs(structType reflect.Type) ([]interface{}, []interface{}) {
	values := make([]interface{}, structType.NumField())
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	return values, scanArgs
}

func validateStructAndColumns(ms *MS, values []interface{}) error {
	columns, err := ms.Rows.Columns()
	if err != nil {
		//when Rows are closed, error occur.
		ms.Err = err
		return ms.Err
	}
	if len(columns) != len(values) {
		ms.Err = fmt.Errorf("number of struct field(%d) doesn't match to columns of sql(%d).", len(values), len(columns))
		return ms.Err
	}
	return nil
}

//Set data
func scanStruct(values []interface{}, v reflect.Value) {
	structType := v.Type()
	for i := 0; i < structType.NumField(); i++ {
		val := reflect.ValueOf(values[i])
		switch val.Kind() {
		case reflect.Invalid:
			//nil: for now, it skips.
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.Field(i).Set(reflect.ValueOf(u.Itoi(values[i])))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v.Field(i).Set(reflect.ValueOf(u.ItoUi(values[i])))
		case reflect.Bool:
			v.Field(i).Set(reflect.ValueOf(u.Itob(values[i])))
		case reflect.String:
			v.Field(i).Set(reflect.ValueOf(u.Itos(values[i])))
		case reflect.Slice:
			if u.CheckInterface(values[i]) == "[]uint8" {
				v.Field(i).Set(reflect.ValueOf(u.ItoBS(values[i])))
			}
		//case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Map:
		case reflect.Struct:
			//time.Time
			if u.CheckInterface(values[i]) == "time.Time" {
				v.Field(i).Set(reflect.ValueOf(u.ItoT(values[i]).Format(tFomt)))
			}
		default: // reflect.Array, reflect.Struct, reflect.Interface
			v.Field(i).Set(reflect.ValueOf(values[i]))
		}
	}
	return
}

// SELECT : Get All field you set(Though you get only record, use it.)
func (ms *MS) Select(selectSQL string, args ...interface{}) ([]map[string]interface{}, []string, error) {
	defer times.Track(time.Now(), "SelectSQLAllField()")
	//540.417µs

	//1. create sql and exec
	//rows, err := self.db.Query("SELECT * FROM t_users WHERE delete_flg=?", "0")
	rows, err := ms.DB.Query(selectSQL, args...)
	if err != nil {
		return nil, nil, err
	}

	return ms.convertRowsToMaps(rows)
}

// Convert result of select into Map[] type. Return multiple array map and interface(plural lines)
func (ms *MS) convertRowsToMaps(rows *sql.Rows) ([]map[string]interface{}, []string, error) {
	defer times.Track(time.Now(), "convertRowsToMaps()")
	//convertRowsToMaps() took 85.191µs

	// Get column name
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	values := make([]interface{}, len(columns))

	// rows.Scan は引数に `[]interface{}`が必要.
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	retMaps := []map[string]interface{}{}
	//
	for rows.Next() { //true or false
		//Get data into scanArgs
		err = rows.Scan(scanArgs...)

		if err != nil {
			return nil, columns, err
		}

		rowdata := map[string]interface{}{}

		//var v string
		for i, value := range values {
			if u.CheckInterface(value) == "[]uint8" {
				value = u.ItoBS(value)
			} else if u.CheckInterface(value) == "time.Time" {
				value = u.ItoT(value).Format(tFomt)
			}

			// Here we can check if the value is nil (NULL value)
			//if value == nil {
			//	v = "NULL"
			//} else {
			//	v = string(value)
			//}

			//if b, ok := value.([]byte); ok{
			//	v = string(b)
			//} else {
			//	v = "NULL"
			//}

			//rowdata[columns[i]] = v
			rowdata[columns[i]] = value
		}
		retMaps = append(retMaps, rowdata)
	}
	return retMaps, columns, nil
}

//-----------------------------------------------------------------------------
// Insert
//-----------------------------------------------------------------------------
func (self *MS) Insert(sql string, args ...interface{}) (int64, error) {
	//1.creates a prepared statement (placeholder)
	//insertSQL := "INSERT t_users SET first_name=?, last_name=?"
	stmt, err := self.DB.Prepare(sql)
	if err != nil {
		return 0, err
	}

	//2.set parameter to prepared statement
	//res, err := stmt.Exec("mitsuo", "fujita")
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	defer stmt.Close() //statementもcloseする必要がある

	//3.Get id from response
	//id, err := res.LastInsertId()
	return res.LastInsertId()
}

//-----------------------------------------------------------------------------
// UPDATE / DELETE
//-----------------------------------------------------------------------------
func (ms *MS) Exec(sql string, args ...interface{}) (int64, error) {

	//1.creates a prepared statement (placeholder)
	//updateSQL := "UPDATE t_users SET first_name=? WHERE user_id=?"
	stmt, err := ms.DB.Prepare(sql)
	if err != nil {
		return 0, err
	}

	//2.set parameter to prepared statement
	//res, err := stmt.Exec("genjiro", 3)
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	defer stmt.Close() //statementもcloseする必要がある

	//3.Get number of changed rows
	//rows, err := res.RowsAffected()
	return res.RowsAffected()
}

// Execution simply
func (ms *MS) Exec2(sql string, args ...interface{}) error {
	//result, err := self.db.Exec("INSERT t_users SET first_name=?, last_name=?", "Mika", "Haruda")
	_, err := ms.DB.Exec(sql, args...)
	return err
}
