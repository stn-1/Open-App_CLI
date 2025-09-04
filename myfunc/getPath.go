package myfunc

import (
	//"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/windows/registry"
)

// Kết nối DB (tạo file apps.db nếu chưa có)
func InitDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "data/apps.db")
    if err != nil {
        return nil, err
    }

    // Thư mục "data" phải tồn tại; nếu chưa có thì tạo (tùy bạn làm ở ngoài)
    if err := db.Ping(); err != nil {
        return nil, err
    }

    // Tạo bảng nếu chưa có
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS apps (
            id   INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            path TEXT NOT NULL
        );
    `)
    if err != nil {
        return nil, err
    }

    // Đảm bảo name là duy nhất để UPSERT hoạt động
    _, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_apps_name ON apps(name);`)
    if err != nil {
        return nil, err
    }

    return db, nil
}

// Lưu app vào DB
func SaveAppToDB(db *sql.DB, name, path string) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }
    _, err := db.Exec(`
        INSERT INTO apps (name, path)
        VALUES (?, ?)
        ON CONFLICT(name) DO UPDATE SET
            path = excluded.path;
    `, name, path)
    return err
}
// Lấy tất cả app

func GetAllAppsByName(db *sql.DB)  error {
    allApps := make(map[string]string)

    // 1. Registry
    registryApps, err := GetInstalledAppsByName()
    if err != nil {
        return err
    }
    for name, path := range registryApps {
        allApps[name] = path
        if err := SaveAppToDB(db, name, path); err != nil {
            // log thay vì nuốt lỗi
            log.Printf("Save registry app failed (%s): %v", name, err)
        }
    }

    // 2. Running processes
    runningApps, err := GetRunningProcessesByName()
    if err != nil {
        return err
    }
    for name, path := range runningApps {
        if _, exists := allApps[name]; !exists {
            allApps[name] = path
        }
        if err := SaveAppToDB(db, name, path); err != nil {
            log.Printf("Save running app failed (%s): %v", name, err)
        }
    }

    return nil
}
func GetAppsFromDB(db *sql.DB) (map[string]string, error) {
    apps := make(map[string]string)

    rows, err := db.Query("SELECT name, path FROM apps")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var name, path string
        if err := rows.Scan(&name, &path); err != nil {
            return nil, err
        }
        apps[name] = path
    }

    // kiểm tra lỗi khi duyệt rows
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return apps, nil
}

// Lấy app từ registry
func GetInstalledAppsByName() (map[string]string, error) {
	apps := make(map[string]string)
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
				apps[strings.ToLower(displayName)] = exePath
			}
			sub.Close()
		}
	}
	return apps, nil
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
