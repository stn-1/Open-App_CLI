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
	apps:=[]string{"zalo","google"}
	app,_:=myfunc.ReturnAppList(db,apps)
	for _,a := range app{
		fmt.Println(a)
	}
}
