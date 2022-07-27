package main

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

//根據不同的情境去設定SetMaxOpenConns，SetMaxIdleConns和SetConnMaxLifetime，可以有效提高性能和降低系统資源消耗

//SetMaxOpenConns 設定最大的連接數量，Default = 0，表示不限制。
//SetMaxIdleConns 設定connection pool最大閒置的連接數量，Default = 2。
//SetConnMaxLifetime 設定每個connection可以重複使用的最大時效，Default = 0,表示不限制

//Result: 一秒內運行次數，平均每次的執行時間

func insertRecord(b *testing.B, db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO isbns VALUES ('978-3-598-21500-1')")
	if err != nil {
		b.Fatal(err)
	}
}

////default runtime.GOMAXPROCS(0), GOMAXPROCS sets the maximum number of CPUs that can be executing.
//// 總核心數目：	10
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

////Performance for MaxOpenConns with different setting (1,2,5,10,20,Unlimited)
////If only open 5 connection and all connection is in used, the request will need to wait util connection release
func BenchmarkMaxOpenConns1(b *testing.B) {

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
	if err != nil {
		b.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			insertRecord(b, db)
		}
	})
}

func BenchmarkMaxOpenConns2(b *testing.B) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
	if err != nil {
		b.Fatal(err)
	}
	db.SetMaxOpenConns(2)
	defer db.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			insertRecord(b, db)
		}
	})
}

func BenchmarkMaxOpenConns5(b *testing.B) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
	if err != nil {
		b.Fatal(err)
	}
	db.SetMaxOpenConns(5)
	defer db.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			insertRecord(b, db)
		}
	})
}

func BenchmarkMaxOpenConns10(b *testing.B) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
	if err != nil {
		b.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	defer db.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			insertRecord(b, db)
		}
	})
}

func BenchmarkMaxOpenConnsUnlimited(b *testing.B) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			insertRecord(b, db)
		}
	})
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//// if set to none = every time need to create connection after insert data
//// 代價: 佔用系統memory
//// MySQL的wait_timeout在default的設定下會自動關閉8小時内未使用的connection 在這個情況下 sql.DB會自動重試兩次，若還是無法連線，將會從connection pool中刪除並且create新的connection
//// MaxIdleConns應該小於或等於MaxOpenConns。Go會檢查並且自動減少MaxIdleConns
//
//func BenchmarkMaxIdleConnsNone(b *testing.B) {
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
//	if err != nil {
//		b.Fatal(err)
//	}
//	db.SetMaxIdleConns(0)
//	defer db.Close()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			insertRecord(b, db)
//		}
//	})
//}
//
//func BenchmarkMaxIdleConns1(b *testing.B) {
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
//	if err != nil {
//		b.Fatal(err)
//	}
//	db.SetMaxIdleConns(1)
//	defer db.Close()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			insertRecord(b, db)
//		}
//	})
//}
//
//func BenchmarkMaxIdleConns2(b *testing.B) {
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
//	if err != nil {
//		b.Fatal(err)
//	}
//	db.SetMaxIdleConns(2)
//	defer db.Close()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			insertRecord(b, db)
//		}
//	})
//}
//
//func BenchmarkMaxIdleConns5(b *testing.B) {
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
//	if err != nil {
//		b.Fatal(err)
//	}
//	db.SetMaxIdleConns(5)
//	defer db.Close()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			insertRecord(b, db)
//		}
//	})
//}
//
//func BenchmarkMaxIdleConns10(b *testing.B) {
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
//	if err != nil {
//		b.Fatal(err)
//	}
//	db.SetMaxIdleConns(10)
//	defer db.Close()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			insertRecord(b, db)
//		}
//	})
//}
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
////The connection will expire 1 hour after it was first created — not 1 hour after it last became idle.
//
//func BenchmarkConnMaxLifetimeUnlimited(b *testing.B) {
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
//	if err != nil {
//		b.Fatal(err)
//	}
//	db.SetConnMaxLifetime(0)
//	defer db.Close()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			insertRecord(b, db)
//		}
//	})
//}
//
//func BenchmarkConnMaxLifetime1000(b *testing.B) {
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
//	if err != nil {
//		b.Fatal(err)
//	}
//	db.SetConnMaxLifetime(1000 * time.Millisecond)
//	defer db.Close()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			insertRecord(b, db)
//		}
//	})
//}
//
//func BenchmarkConnMaxLifetime500(b *testing.B) {
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
//	if err != nil {
//		b.Fatal(err)
//	}
//	db.SetConnMaxLifetime(500 * time.Millisecond)
//	defer db.Close()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			insertRecord(b, db)
//		}
//	})
//}
//
//func BenchmarkConnMaxLifetime50(b *testing.B) {
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
//	if err != nil {
//		b.Fatal(err)
//	}
//	db.SetConnMaxLifetime(50 * time.Millisecond)
//	defer db.Close()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			insertRecord(b, db)
//		}
//	})
//}
//
//func BenchmarkConnMaxLifetime10(b *testing.B) {
//	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
//	if err != nil {
//		b.Fatal(err)
//	}
//	db.SetConnMaxLifetime(10 * time.Millisecond)
//	defer db.Close()
//
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			insertRecord(b, db)
//		}
//	})
//}
