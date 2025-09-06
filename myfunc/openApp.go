package myfunc

import (
	"database/sql"
	"fmt"
	"os/exec"
)

// OpenApps mở ứng dụng theo đường dẫn/tên, trả về error nếu có lỗi
func OpenApps(db *sql.DB, AppID int) error {
	var name, path string
	// Lấy ra name và path từ DB
	err := db.QueryRow("SELECT name, path FROM resources WHERE id = ?", AppID).Scan(&name, &path)
	if err != nil {
		return fmt.Errorf("Không tìm thấy ứng dụng với id %d: %v", AppID, err)
	}
	// Tạo lệnh mở app
	cmd := exec.Command(path)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Không thể mở ứng dụng %s: %v", name, err)
	}
	fmt.Printf("Đang mở ứng dụng %s...\n", name)
	return nil
}

