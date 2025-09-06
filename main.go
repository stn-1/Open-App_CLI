package main

import (
	//"fmt"
	"fmt"
	"myapp/myfunc"
)
func main(){
	db,err:=myfunc.InitDB()
	if(err!=nil){
		fmt.Println(err)
	} 
	err=myfunc.GetAllResourcesByName(db)
	if(err!=nil){
		fmt.Println(err)
	}
	myfunc.ShowDB(db)
	apps:=[]string{"zalo","google"}
	myfunc.FindRes(db,apps)
	
}
