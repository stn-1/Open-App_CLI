package main

import (
	"bufio"
	"fmt"
	"myapp/myfunc"
	"os"
	"strconv"
	"strings"
)

func main() {
    db, err := myfunc.InitDB()
    if err != nil {
        fmt.Println(err)
        return
    }
    defer db.Close()

    err = myfunc.GetAllResourcesByName(db)
    if err != nil {
        fmt.Println(err)
    }

    reader := bufio.NewReader(os.Stdin)

    fmt.Println("Enter what you want")
    for {
        fmt.Print(">>> ")
        line, _ := reader.ReadString('\n')
        line = strings.TrimSpace(line)

        if line == "quit" {
            fmt.Println("Bye!")
            break
        }

        parts := strings.Fields(line) // tách theo dấu cách
        if len(parts) == 0 {
            continue
        }

        switch parts[0] {
        case "reload":
            err = myfunc.GetAllResourcesByName(db)
            if err != nil {
                fmt.Println(err)
            }
        case "showdb":
            myfunc.ShowDB(db)
        case "find":
            if len(parts) < 2 {
                fmt.Println("Bạn cần nhập ít nhất một tên app sau 'find'")
                continue
            }
            myfunc.FindRes(db, parts[1:]) // truyền danh sách app
		case "makegroup":
			listID:=[]int{}
			if len(parts) < 3 {
			fmt.Println("Bạn cần nhập ít nhất một tên app sau 'find'")
				continue
			}
			for _,id := range parts[2:]{
				idres,err:=strconv.Atoi(id)
				if(err!=nil){
					fmt.Println("❌ Lỗi chuyển ID:", err)
           			continue
				}
				listID = append(listID,idres)
			}
			err:=myfunc.CreateGroup(db,parts[1],listID)
			if err!=nil{
				fmt.Println("lỗi ở hàm CreateGroup ",err)
			}
		case "del":
			if len(parts) < 2 {
                fmt.Println("Bạn cần nhập ít nhất một tên group")
                continue
            }
			myfunc.DeleteGroupByName(db,parts[1])
		case "addweb":
			if len(parts) < 3 {
                fmt.Println("Bạn cần nhập theo format addweb+<name>+<path>")
                continue
            }
			myfunc.SaveWebToDB(db,parts[1],parts[2])
			fmt.Println("web đã được lưu vào database")
		case "showgroup":
			myfunc.ShowGroups(db)
		case "rungroup":
			myfunc.RunGroup(db,parts[1])
		case "delres":
			if len(parts) < 3 {
				fmt.Println("Bạn cần nhập theo format: delres <id|name> <giá trị>")
				continue
			}
			mode := parts[1]
			value := parts[2]

			if mode == "id" {
				id, err := strconv.Atoi(value)
				if err != nil {
					fmt.Println("❌ ID phải là số:", err)
					continue
				}
				err = myfunc.DeleteResourceByID(db, id)
				if err != nil {
					fmt.Println(err)
				}
			} else if mode == "name" {
				err = myfunc.DeleteResourceByName(db, value)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("⚠️ Cách dùng: delres <id|name> <giá trị>")
			}
		case "help":
			fmt.Println("📌 Danh sách lệnh có sẵn:")
			fmt.Println("  reload                - Tải lại dữ liệu resource từ DB")
			fmt.Println("  showdb                - Hiển thị toàn bộ resource trong DB")
			fmt.Println("  find <name...>        - Tìm resource theo tên")
			fmt.Println("  makegroup <gname> <ids...> - Tạo group mới với danh sách resource ID")
			fmt.Println("  del <gname>           - Xóa group theo tên")
			fmt.Println("  addweb <name> <url>   - Lưu một web app vào DB")
			fmt.Println("  showgroup             - Hiển thị toàn bộ group trong DB")
			fmt.Println("  rungroup <gname>      - Chạy tất cả resource trong group")
			fmt.Println("  quit                  - Thoát chương trình")
			fmt.Println("  help                  - Hiển thị danh sách lệnh")
			fmt.Println("  delres id <id>        - Xóa resource theo id (ví dụ: delres id 3)")
    		fmt.Println("  delres name <name>    - Xóa resource theo tên (ví dụ: delres name chrome)")
        default:
            fmt.Println("Lệnh không hợp lệ:", parts[0])
        }
    }
}
