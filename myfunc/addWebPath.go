package myfunc

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)
func InitWebDB() (*sql.DB, error) {
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
        CREATE TABLE IF NOT EXISTS webs (
            id   INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            path TEXT NOT NULL
        );
    `)
    if err != nil {
        return nil, err
    }

    // Đảm bảo name là duy nhất để UPSERT hoạt động
    _, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_webs_name ON webs(name);`)
    if err != nil {
        return nil, err
    }

    return db, nil
}
func SaveWebToDB(db *sql.DB, name, path string) error {
    if db == nil {
        return fmt.Errorf("db is nil")
    }
    _, err := db.Exec(`
        INSERT INTO webs (name, path)
        VALUES (?, ?)
        ON CONFLICT(name) DO UPDATE SET
            path = excluded.path;
    `, name, path)
    return err
}
func GetAllURLsByInput(db *sql.DB, urls map[string]string)  error {
    for name, path := range urls {
        if err := SaveWebToDB(db, name, path); err != nil {
            // log thay vì nuốt lỗi
            log.Printf("Save registry web failed (%s): %v", name, err)
        }
    }
    return  nil
}
func GetURLsFromDB(db *sql.DB) (map[string]string, error) {
    urls := make(map[string]string)

    rows, err := db.Query("SELECT name, path FROM webs")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var name, path string
        if err := rows.Scan(&name, &path); err != nil {
            return nil, err
        }
        urls[name] = path
    }

    // kiểm tra lỗi khi duyệt rows
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return urls, nil
}
