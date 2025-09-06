package myfunc

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/windows/registry"
)

// Kết nối DB (tạo file resources.db nếu chưa có)
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "data/resources.db")
	if err != nil {
		return nil, err
	}

	// Thư mục "data" phải tồn tại; nếu chưa có thì tạo (tùy bạn làm ở ngoài)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Tạo bảng nếu chưa có
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS resources (
            id      INTEGER PRIMARY KEY AUTOINCREMENT,
            name    TEXT NOT NULL,
            path    TEXT NOT NULL,
            is_web  BOOLEAN NOT NULL
        );

        CREATE TABLE IF NOT EXISTS groups (
            id      INTEGER PRIMARY KEY AUTOINCREMENT,
            nameG   TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS group_resources (
            group_id    INTEGER NOT NULL,
            resource_id INTEGER NOT NULL,
            PRIMARY KEY (group_id, resource_id),
            FOREIGN KEY (group_id) REFERENCES groups(id),
            FOREIGN KEY (resource_id) REFERENCES resources(id)
        );
    `)
	if err != nil {
		return nil, err
	}

	// Đảm bảo name là duy nhất để UPSERT hoạt động
	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_resources_name ON resources(name);`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Lưu resource vào DB
func SaveResourceToDB(db *sql.DB, name, path string, isWeb bool) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	_, err := db.Exec(`
        INSERT INTO resources (name, path, is_web)
        VALUES (?, ?, ?)
        ON CONFLICT(name) DO UPDATE SET
            path = excluded.path;
    `, name, path, isWeb)
	return err
}

// Lấy tất cả resources (apps + webs)
func GetAllResourcesByName(db *sql.DB) error {
	allResources := make(map[string]string)

	// 1. Registry
	registryResources, err := GetInstalledResourcesByName()
	if err != nil {
		return err
	}
	for name, path := range registryResources {
		allResources[name] = path
		if err := SaveResourceToDB(db, name, path, false); err != nil {
			log.Printf("Save registry resource failed (%s): %v", name, err)
		}
	}

	// 2. Running processes
	runningResources, err := GetRunningProcessesByName()
	if err != nil {
		return err
	}
	for name, path := range runningResources {
		if _, exists := allResources[name]; !exists {
			allResources[name] = path
		}
		if err := SaveResourceToDB(db, name, path, false); err != nil {
			log.Printf("Save running resource failed (%s): %v", name, err)
		}
	}

	return nil
}

// Lấy resources từ DB
func GetResourcesFromDB(db *sql.DB) (map[string]string, error) {
	resources := make(map[string]string)

	rows, err := db.Query("SELECT name, path FROM resources")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name, path string
		if err := rows.Scan(&name, &path); err != nil {
			return nil, err
		}
		resources[name] = path
	}

	// kiểm tra lỗi khi duyệt rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}
//in ra dữ liệu lấy từ GetResourcesFromDB
func ShowDB(db *sql.DB){
    res,err:=GetResourcesFromDB(db)
    if(err!=nil){
        log.Fatal(err)
    }
    for name,path:=range res{
        fmt.Println(name,path)
    }
}

// Lấy resource (ứng dụng) từ registry
func GetInstalledResourcesByName() (map[string]string, error) {
	resources := make(map[string]string)
	registryKeys := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	for _, regKey := range registryKeys {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, regKey, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		defer k.Close()

		subKeys, err := k.ReadSubKeyNames(0)
		if err != nil {
			continue
		}

		for _, subKey := range subKeys {
			sub, err := registry.OpenKey(k, subKey, registry.QUERY_VALUE)
			if err != nil {
				continue
			}

			displayName, _, _ := sub.GetStringValue("DisplayName")
			exePath, _, _ := sub.GetStringValue("DisplayIcon")
			if exePath == "" {
				exePath, _, _ = sub.GetStringValue("InstallLocation")
				if exePath != "" {
					exePath = filepath.Join(exePath, subKey+".exe")
				}
			}

			if exePath != "" && displayName != "" {
				exePath = strings.Split(strings.Trim(exePath, `"`), ",")[0]
				resources[strings.ToLower(displayName)] = exePath
			}
			sub.Close()
		}
	}
	return resources, nil
}

// Lấy process đang chạy
func GetRunningProcessesByName() (map[string]string, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	processMap := make(map[string]string)
	for _, p := range procs {
		path, err := p.Exe()
		if err != nil || path == "" {
			continue
		}

		path = filepath.Clean(strings.ToLower(path))
		if name, err := p.Name(); err == nil {
			processMap[strings.ToLower(name)] = path
		}
	}
	return processMap, nil
}
