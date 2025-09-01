package myfunc

import (
	"fmt"
	"os/exec"
)

// OpenApps mở ứng dụng theo đường dẫn/tên, trả về error nếu có lỗi
func OpenApps(SelectApp string) error {
	cmd := exec.Command(SelectApp)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Không thể mở ứng dụng %s: %v", SelectApp, err)
	}
	fmt.Printf("Đang mở ứng dụng %s...\n", SelectApp)
	return nil
}
