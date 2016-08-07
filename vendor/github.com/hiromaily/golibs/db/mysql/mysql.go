package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
	//fmt.Printf("dbInfo.db %+v\n", *dbInfo.DB)
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

// SELECT : Get All field you set(Though you get only record, use it.)
func (ms *MS) SelectSQLAllField(selectSQL string, args ...interface{}) ([]map[string]interface{}, []string, error) {
	defer times.Track(time.Now(), "SelectSQLAllField()")
	//540.417µs

	//If no args, set nil

	//1. create sql and exec
	//rows, err := self.db.Query("SELECT * FROM t_users WHERE delete_flg=?", "0")
	rows, err := ms.DB.Query(selectSQL, args...)
	if err != nil {
		return nil, nil, err
	}

	return ms.convertRowsToMaps(rows)
}

//get Rows and return db instance
func (ms *MS) SelectSQLAllFieldIns(selectSQL string, args ...interface{}) *MS {
	defer times.Track(time.Now(), "SelectSQLAllFieldIns()")
	//SelectSQLAllFieldIns() took 471.577µs

	//If no args, set nil

	//1. create sql and exec
	//rows, err := self.db.Query("SELECT * FROM t_users WHERE delete_flg=?", "0")
	ms.Rows, ms.Err = ms.DB.Query(selectSQL, args...)
	if ms.Err != nil {
		fmt.Println("SelectSQLAllFieldIns():[Error]:", ms.Err)
	}

	return ms
}

//set extracted data into parameter variable
func (ms *MS) ScanOne(x interface{}) {
	defer times.Track(time.Now(), "ScanOne()")
	//ScanOne() took 5.23µs

	if ms.Err != nil {
		fmt.Println("ScanOne():[Error]:", ms.Err)
		return
	}

	//e.g)v = person Person
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		ms.Err = errors.New("parameter is not valid. it sould be pointer and not nil.")
		return
	} else {
		if v.Elem().Kind() == reflect.Struct {

			//values := convertToSlice(v.Elem())
			structType := v.Elem().Type()
			values := make([]interface{}, structType.NumField())
			scanArgs := make([]interface{}, len(values))
			for i := range values {
				scanArgs[i] = &values[i]
			}

			//check len(value) and column
			columns, err := ms.Rows.Columns()
			if err != nil {
				//when Rows are closed, error occur.
				ms.Err = err
				return
			}
			if len(columns) != len(values) {
				ms.Err = fmt.Errorf("number of struct field(%d) doesn't match to columns of sql(%d).", len(values), len(columns))
				return
			}

			ret := ms.Rows.Next()
			if !ret {
				ms.Err = errors.New("Rows.Next(): No data")
			}

			//ms.Err = ms.Rows.Scan(values...)
			ms.Err = ms.Rows.Scan(scanArgs...)
			if ms.Err != nil {
				return
			}
			//ms.Err = ms.Rows.Scan(v)
			scan(values, v.Elem())
		} else {
			ms.Err = errors.New("parameter should be pointer of struct slice or struct")
		}
	}
}

//TODO:work in progress
//func (ms *MS) convertRowsToStruct(data interface{}) error {
func (ms *MS) Scan(v ...interface{}) {
	//use gorm source code as a reference
	// gorm/callback_query.go

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		ms.Err = errors.New("parameter is not valid. it sould be pointer and not nil.")
		return
	}

	ms.Rows.Next()
	ms.Err = ms.Rows.Scan(v...)

	return
}

//Set data
func scan(values []interface{}, v reflect.Value) {
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
				//v.Field(i).Set(reflect.ValueOf(u.ItoT(values[i])))
				//plus format datetime
				v.Field(i).Set(reflect.ValueOf(u.ItoT(values[i]).Format(tFomt)))
			}
		default: // reflect.Array, reflect.Struct, reflect.Interface
			fmt.Println(val.Kind(), val.Type(), ":dafault")
			v.Field(i).Set(reflect.ValueOf(values[i]))
		}
	}
	return
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

	//variable for stored each field data
	//values := make([]sql.RawBytes, len(columns)) //it cause error
	values := make([]interface{}, len(columns))

	// rows.Scan は引数に `[]interface{}`が必要.
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		//I don't know why set address to another variable
		//set address of value to variable for scan
		scanArgs[i] = &values[i]
	}

	//retMaps := []map[string]string{}
	//rowdata := map[string]string{}
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
			//Check type
			//val := reflect.ValueOf(value) // ValueOfでreflect.Value型のオブジェクトを取得
			//fmt.Println("val.Type()", val.Type(), "val.Kind()", val.Kind()) // Typeで変数の型を取得

			if u.CheckInterface(value) == "[]uint8" {
				//[]uint8 to []byte to string
				//if tmp, ok := value.([]byte); ok {
				//	//value = strconv.Itoa(int(tmp))
				//	value = string(tmp)
				//}
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
			//fmt.Println(columns[i], ": ", v)
		}
		retMaps = append(retMaps, rowdata)
	}
	return retMaps, columns, nil
}

// Execution simply
func (ms *MS) ExecSQL(sqlString string, args ...interface{}) error {
	//result, err := self.db.Exec("INSERT t_users SET first_name=?, last_name=?", "Mika", "Haruda")
	_, err := ms.DB.Exec(sqlString, args...)
	return err
}

// INSERT
func (self *MS) InesrtSQL(insertSQL string, args ...interface{}) (int64, error) {
	//1.creates a prepared statement (placeholder)
	//insertSQL := "INSERT t_users SET first_name=?, last_name=?"
	stmt, err := self.DB.Prepare(insertSQL)
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

// UPDATE
func (ms *MS) UpdateSQL(updateSQL string, args ...interface{}) (int64, error) {

	//1.creates a prepared statement (placeholder)
	//updateSQL := "UPDATE t_users SET first_name=? WHERE user_id=?"
	stmt, err := ms.DB.Prepare(updateSQL)
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
