package main

import (
	//"fmt"
	"fmt"
	"myapp/myfunc"
)
func main(){
	db,err:=myfunc.InitDB()
	if(err!=nil){
		fmt.Println("error")
	}
	dbweb,err:=myfunc.InitWebDB()
	if(err!=nil){
		fmt.Println("error")
	}
	webs:=map[string]string{"youtube":"https://www.youtube.com/"}
	apps:=[]string{"zalo","google"}
	myfunc.FindApp(db,apps)
	applist:=[]int{59,480}
	weblist:=[]string{"youtube"}
	myfunc.RunListApp(db,applist)
	myfunc.GetAllURLsByInput(dbweb,webs)
	myfunc.FindWeb(dbweb,weblist)
	// for name,a := range apps{
	// 	fmt.Println(name,a)
	// }
}
