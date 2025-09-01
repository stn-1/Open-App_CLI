package myfunc

import (
	"fmt"
)
func GetApp(){
	db,err:=InitDB()
	if(err!=nil){
		fmt.Println("error")
	}
	app,err:=GetAllAppsByName(db)
	for name,_:= range(app){
		fmt.Println(name)
	}
}