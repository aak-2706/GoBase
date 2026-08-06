package gobase

import (
	"encoding/json"
	"fmt"
	"sync"
)

const version = "1.0.1"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}
	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)
type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./"
	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
	employees := []User{
		{
			Name:    "Aman Anilkumar",
			Age:     json.Number("20"),
			Contact: "9876543210",
			Company: "QuantaDB",
			Address: Address{"Bengaluru", "Karnataka", "India", json.Number("560001")},
		},
		{
			Name:    "John Smith",
			Age:     json.Number("28"),
			Contact: "9876543211",
			Company: "Google",
			Address: Address{"Mountain View", "California", "USA", json.Number("94043")},
		},
		{
			Name:    "Emma Johnson",
			Age:     json.Number("25"),
			Contact: "9876543212",
			Company: "Microsoft",
			Address: Address{"Redmond", "Washington", "USA", json.Number("98052")},
		},
		{
			Name:    "Raj Patel",
			Age:     json.Number("31"),
			Contact: "9876543213",
			Company: "Infosys",
			Address: Address{"Pune", "Maharashtra", "India", json.Number("411001")},
		},
		{
			Name:    "Sophia Williams",
			Age:     json.Number("27"),
			Contact: "9876543214",
			Company: "Amazon",
			Address: Address{"Seattle", "Washington", "USA", json.Number("98101")},
		},
		{
			Name:    "Arjun Mehta",
			Age:     json.Number("29"),
			Contact: "9876543215",
			Company: "Flipkart",
			Address: Address{"Bengaluru", "Karnataka", "India", json.Number("560078")},
		},
		{
			Name:    "Olivia Brown",
			Age:     json.Number("26"),
			Contact: "9876543216",
			Company: "Apple",
			Address: Address{"Cupertino", "California", "USA", json.Number("95014")},
		},
		{
			Name:    "Daniel Garcia",
			Age:     json.Number("34"),
			Contact: "9876543217",
			Company: "Meta",
			Address: Address{"Menlo Park", "California", "USA", json.Number("94025")},
		},
		{
			Name:    "Priya Sharma",
			Age:     json.Number("24"),
			Contact: "9876543218",
			Company: "TCS",
			Address: Address{"Hyderabad", "Telangana", "India", json.Number("500081")},
		},
		{
			Name:    "Liam Wilson",
			Age:     json.Number("30"),
			Contact: "9876543219",
			Company: "Netflix",
			Address: Address{"Los Gatos", "California", "USA", json.Number("95032")},
		},
	}
	for _, value := range employees {
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}
	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(records)
	allUsers := []User{}
	for _, f := range records {
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil {
			fmt.Println("Error: ", err)
		}
		allUsers = append(allUsers, employeeFound)
	}
	fmt.Println(allUsers)
	// if err:=db.Delete("user","john");err!=nil{
	// 	fmt.Println("Error: ",err)
	// }
	// if err:=db.Delete("user","");err!=nil{
	// 	fmt.Println("Error: ",err)
	// }

}
