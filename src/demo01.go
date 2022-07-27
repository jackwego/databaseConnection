package main

import (
	"database/sql"
	"fmt"
)

//針對在database的存取上，golang設計了一個sql抽象介面叫做 database/sql，
//工程師可以遵照這個 interface 分開去實作 mysql、sqlite ... 等等，
//如果未來有測試或是抽換需求，只需要更新 driver ，完全不需要更動其他相關程式

// go-sql-driver/mysql

func main() {
	conn := getConnection()
	defer closeConnection(conn)
	idForNewRecord := insertData(conn, "Dr.Lin", 40)
	updateData(conn, 20, idForNewRecord)
	doctor := querySingleRowData(conn, idForNewRecord)
	fmt.Println("Doctor age updated, Age = ", doctor.Age)
	insertData(conn, "Dr.Wu", 20)
	doctorList := queryMultipleRowData(conn, 18)
	fmt.Println("Return size = ", len(*doctorList))
	deleteData(conn, idForNewRecord)

	//Prepare data for transaction
	insertData(conn, "Dr.Huang", 40)
	insertData(conn, "Dr.Chen", 20)
	transaction(conn, "Dr.Huang", "Dr.Chen")
}

func getConnection() *sql.DB {
	dbConfig := GetMysqlConfig()
	//完整的資料格式連線如下
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	fmt.Println(dbConfig.FormatDSN())

	conn, err := sql.Open("mysql", dbConfig.FormatDSN())

	if err != nil {
		panic("[database] cannot open connection" + err.Error())
	} else {
		fmt.Println("[database] successfully connected")
	}

	err = conn.Ping()
	if err != nil {
		panic("[database] ping failed" + err.Error())
	} else {
		fmt.Println("[database] ping successfully")
	}

	return conn
}

func insertData(conn *sql.DB, name string, age int64) int64 {
	result, err := conn.Exec("insert into doctor_tb(name,age,sex,addTime) values(?,?,?,Now())", name, age, 2)
	if err != nil {
		panic("[Insert] insert data error " + err.Error())
	}
	newID, err := result.LastInsertId() //ID for new record
	if err != nil {
		panic("[Insert] get last insert id error " + err.Error())
	}
	i, err := result.RowsAffected()
	if err != nil {
		panic("[Insert] get number of row affected error " + err.Error())
	}

	fmt.Printf("[Insert] ID for new record：%d , Rows afected：%d \n", newID, i)

	return newID
}

func deleteData(conn *sql.DB, id int64) {
	result, err := conn.Exec("delete from doctor_tb where id = ?", id)
	if err != nil {
		panic("[Delete] delete data error " + err.Error())
	}
	i, _ := result.RowsAffected()
	fmt.Printf("[Delete] Rows Affected：%d \n", i)
}

func querySingleRowData(conn *sql.DB, id int64) Doctor {
	var doctor Doctor
	row := conn.QueryRow("select * from doctor_tb where id = ?", id)

	//If QueryRow return more than 1 result, Scan method discards all but the first.
	if err := row.Scan(&doctor.ID, &doctor.Name, &doctor.Age, &doctor.Sex, &doctor.AddTime); err != nil {
		panic("[QueryRow] Query row error " + err.Error())
	} else {
		fmt.Println("[QueryRow] Query result：", doctor)
	}

	return doctor
}

func queryMultipleRowData(conn *sql.DB, age int64) *[]Doctor {
	rows, err := conn.Query("select * from doctor_tb where age > ?", age)
	if err != nil {
		panic("[QueryMultipleRow] Query Multiple Row error " + err.Error())
	}

	var docList []Doctor
	for rows.Next() {
		var doc2 Doctor
		rows.Scan(&doc2.ID, &doc2.Name, &doc2.Age, &doc2.Sex, &doc2.AddTime)
		//加入数组
		docList = append(docList, doc2)
	}
	fmt.Println("Query Multiple Row", docList)

	return &docList
}

func transaction(conn *sql.DB, nameForFirstDoctor string, nameForSecondDoctor string) {
	tx, _ := conn.Begin()
	firstResult, _ := tx.Exec("update doctor_tb set age = age + 1 where name = ?", nameForFirstDoctor)
	secondResult, _ := tx.Exec("update doctor_tb set age = age + 1 where name = ?", nameForSecondDoctor)

	//Rows Affected，if RowsAffected=0, means update failed
	firstRowsAffected, _ := firstResult.RowsAffected()
	secondRowsAffected, _ := secondResult.RowsAffected()
	if firstRowsAffected > 0 && secondRowsAffected > 0 {
		//Commit if both update successfully
		err := tx.Commit()
		if err != nil {
			panic("[Transaction] Commit failed" + err.Error())
		}
		fmt.Println("[Transaction] Commit success")
	} else {
		//Rollback if fail updated
		err := tx.Rollback()
		if err != nil {
			panic("[Transaction] Rollback failed " + err.Error())
		}
		fmt.Println("[Transaction] Rollback success")
	}
}

func updateData(conn *sql.DB, age int64, id int64) {
	result, err := conn.Exec("update doctor_tb set age=? where id = ?", age, id)
	if err != nil {
		panic("[Update] Update data error " + err.Error())
	}

	i, _ := result.RowsAffected()
	if err != nil {
		panic("[Update] get number of row affected error " + err.Error())
	}

	fmt.Printf("[Update] Rows afected：%d \n", i)
}

func closeConnection(conn *sql.DB) {

	if conn != nil {
		conn.Close()
		fmt.Println("[database] successfully closed")
	} else {
		fmt.Println("[database] connection == nil")
	}

}
