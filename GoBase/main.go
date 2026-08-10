package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const version = "1.0.1"

type (
	Logger interface { //
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}
	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex //
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir) //
	opts := Options{}
	if options != nil {
		opts = *options
	}
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO)) //
	}
	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}
	if _, err := os.Stat(dir); err != nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir) //
		return &driver, nil
	}
	opts.Logger.Debug("Creating the databse at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing Collection - no place to save record!")
	}
	if resource == "" {
		return fmt.Errorf("Missing Someone - unable to save record(no name)!")
	}
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()
	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tempPath := fnlPath + ".tmp"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	b = append(b, byte('\n'))
	if err := ioutil.WriteFile(tempPath, b, 0644); err != nil {
		return err
	}
	return os.Rename(tempPath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - unable tp read!")
	}
	if resource == "" {
		return fmt.Errorf("Missing resource - unable to read record(no name)!")
	}
	record := filepath.Join(d.dir, collection, resource)
	if _, err := stat(record); err != nil {
		return err
	}
	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("missing Collection - unable to read")
	}
	dir := filepath.Join(d.dir, collection)
	if _, err := stat(dir); err != nil {
		return nil, err
	}
	files, _ := ioutil.ReadDir(dir)
	var records []string
	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(b))
	}
	return records, nil
}

func (d *Driver) Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()
	dir := filepath.Join(d.dir, path)
	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		fmt.Errorf("unable to find file or directory named %v\n", path)
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

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
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil { //
			fmt.Println("Error: ", err)
		}
		allUsers = append(allUsers, employeeFound)
	}
	fmt.Println(allUsers)
	if err := db.Delete("users", "John Smith"); err != nil {
		fmt.Println("Error: ", err)
	}
}
