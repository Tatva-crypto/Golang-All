package main

//run this go file as this is the final code for rest api with db gorm mux and cors is coming soon
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// type Customer struct {
// 	gorm.Model
// 	Name    string `json:"name"`
// 	OrderId int    `json:"order_id"`
// }
// type Order struct {
// 	gorm.Model
// 	Number   uint     `json:"number"`
// 	Customer Customer `gorm:"ForeignKey:OrderId"`
// }
type Department struct {
	gorm.Model
	Depno    uint     `json:"dep_no"`
	DepName  string   `json:"dep_name"`
	Employee Employee `gorm:"ForeignKey:DepartmentId"`
	Project  Project  `gorm:"ForeignKey:DepartmentId"`
}

type Employee struct {
	gorm.Model
	Name         string `json:"name"`
	Idno         int    `json:"idno"`
	Age          int    `json:"age"`
	DepartmentId int    `json:"dep_id"`
}
type Project struct {
	gorm.Model
	Name         string `json:"name"`
	DepartmentId int    `json:"dep_id"`
}

var db *gorm.DB
var e error

func main() {
	db, e = gorm.Open("postgres", "user=postgres password=tatva972000 dbname=postgres sslmode=disable")
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println("Connection Established")
	}
	defer db.Close()
	db.SingularTable(true)
	db.AutoMigrate(&Employee{}, &Department{}, &Project{})
	r := mux.NewRouter()
	r.HandleFunc("/", getnames).Methods("GET")
	r.HandleFunc("/", addcustomer).Methods("POST")
	//var d Department
	//fmt.Println("%T", d)
	http.ListenAndServe(":1991", r)
}

//this is join query
//select dep_name,e.name,e.age,p.name from employee e inner join department d on
//e.department_id=d.id inner join project p on p.department_id=d.id;

func getnames(w http.ResponseWriter, r *http.Request) {
	// var customers []Customer
	type Result struct {
		Depno    uint   `json:"dep_no"`
		Dep_Name string `json:"dep_name"`
		// Project_Name    string `json:"project_name"`
		Age  int    `json:"age"`
		Name string `json:"name"`
	}
	var result Result

	if e := db.Table("employee").Select("department.depno,department.dep_name,employee.age,project.name").Joins("JOIN department on department.id = employee.department_id").Joins("JOIN project on project.department_id=department.id").Find(&result).Error; e != nil { //there has to be some changes regarding the joins bcs only find is not getting all the data only the data that i am asking for which is dep
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Response-Code", "06")
		w.Header().Set("Response-Desc", "Data Not Found")
		w.WriteHeader(404)
		w.Write([]byte(`{"message":"data not found"}`))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Response-Code", "00")
		w.Header().Set("Response-Desc", "Success")
		json.NewEncoder(w).Encode(result)
	}
	// db.Table("employee").Select("dep_name,employee.name,employee.age").Joins("full join department on employee.department_id = department.id").Scan(&result)

	// db.Table("employee").Select("department.depno,department.dep_name,employee.age,project.name").Joins("JOIN department on department.id = employee.department_id").Joins("JOIN project on project.department_id=department.id").Find(&result)
	// json.NewEncoder(w).Encode(result)
}

// db.Model(&Employee{}).Select("dep_name, employees.name,employees.age,p.name").Joins("inner join project p on p.department_id=d.id").Scan(&result{})
// db.Joins("JOIN employee ON employee.department_id = department.id").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Find(&result)``
func addcustomer(w http.ResponseWriter, r *http.Request) {
	var department Department //why i initiated Department variable here?

	var _ = json.NewDecoder(r.Body).Decode(&department)

	db.Create(&department)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Response-Code", "00")
	w.Header().Set("Response-Desc", "Success")
	json.NewEncoder(w).Encode(department)

}
