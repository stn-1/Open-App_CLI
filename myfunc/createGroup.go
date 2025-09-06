package myfunc

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CreateGroup(db *sql.DB, name string) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	_, err := db.Exec(`INSERT INTO groups (nameG) VALUES (?)`, name)
	if err != nil {
		return fmt.Errorf("failed to insert group: %w", err)
	}

	return nil
}
func DeleteGroup(db *sql.DB, name string) error{
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	_,err:=db.Exec(`DELETE FROM groups WHERE nameG=?`,name)
	if(err!=nil){
		return fmt.Errorf("failed to delete group %s: %w", name, err)
	}
	return nil
}
func RunGroup(db *sql.DB, ids []int) error {
	for _, id := range ids {
		var isWeb bool
		// Truy vấn cờ phân biệt
		err := db.QueryRow("SELECT is_web FROM resources WHERE id = ?", id).Scan(&isWeb)
		if err != nil {
			fmt.Printf("❌ Không tìm thấy resource với id %d: %v\n", id, err)
			continue // bỏ qua cái lỗi này, chạy tiếp cái khác
		}

		if isWeb {
			if err := openURL(db, id); err != nil {
				fmt.Printf("⚠️ Lỗi mở web id=%d: %v\n", id, err)
			}
		} else {
			if err := OpenApps(db, id); err != nil {
				fmt.Printf("⚠️ Lỗi mở app id=%d: %v\n", id, err)
			}
		}
	}
	return nil
}