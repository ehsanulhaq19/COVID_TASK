package main

import (
    "fmt"
	"net"
	"encoding/json"
    "encoding/csv"
    "os"
    "io"
)

type Query struct { //main query object
	Query Fields `json:"query"`
}


type Fields struct { //contains data fields of the query
	Region string `json:"region"`
	Date string `json:"date"`
}

type Data struct { //contain data to be displayed
	Date string `json:"date"`
	Positive string `json:"positive"`
	Tests string `json:"tests"`
	Expired string `json:"expired"`
	Admitted string `json:"admitted"`
	Discharged string `json:"discharged"`
	Region string `json:"region"`
}

type Result struct { //contain the whole dataset respone
	DataSet []Data `json:"response"`
}

type Node struct { //used to store data in link list
	date string
	positive string
	tests string
	expired string
	admitted string
	discharged string
	region string
	next *Node
	prev *Node
}
  
type List struct { //contains list of nodes
    head *Node
    tail *Node
}

func server() { //create server
	// listen on a port
	ln, err := net.Listen("tcp", ":4040")
	if err != nil {
	  fmt.Println(err)
	  return
	}
	for {
	  // accept a connection
	  c, err := ln.Accept()
	  if err != nil {
		fmt.Println(err)
		continue
	  }
	  // handle the connection
	  go handleServerConnection(c)
	}
  }
  
  func handleServerConnection(c net.Conn) { //maintain server connection and fetch data from csv file acccording to the query provided

    // we create a decoder that reads directly from the socket
    d := json.NewDecoder(c)
	var msg Query
    err := d.Decode(&msg)
	fmt.Println(msg, err)
	// fmt.Println(msg.Query.Date)
	data := ReadCsvFile("./covid_final_data.csv", msg.Query.Date, msg.Query.Region)
	jsonData := data.ConvertToJSON()
	fmt.Println(jsonData)
	

	c.Close()
	return

}

func (L *List) Insert(date string, positive string, tests string, expired string , admitted string, discharged string , region string) {
    list := &Node{
        next: L.head,
		date:  date,
		positive : positive,
		tests : tests,
		expired : expired,
		admitted : admitted,
		discharged : discharged,
		region : region,
	}
	
    if L.head != nil {
        L.head.prev = list
    }
    L.head = list
 
    l := L.head
    for l.next != nil {
        l = l.next
    }
    L.tail = l
}

func (l *List) Display() {
    list := l.head
    for list != nil {
        fmt.Printf("%+v ->", list.region)
        list = list.next
    }
    fmt.Println()
}


func (l *List) ConvertToJSON()(jsonData string) {
	var result Result
	var dataSet []Data
	list := l.head
    for list != nil {
		dataSet = append(dataSet, Data{
			Date : list.date,
			Positive : list.positive,
			Tests : list.tests,
			Expired : list.expired,
			Admitted : list.admitted,
			Discharged : list.discharged,
			Region : list.region,
		})
        list = list.next
	}
	result.DataSet = dataSet

	fmt.Println()

	pagesJson, err := json.Marshal(result)
    if err != nil {
        fmt.Printf("Cannot encode to JSON ", err)
	}
	jsonData = string(pagesJson)
	return 
}


  
func ReadCsvFile(filePath string, date string, region string)(link List)  {
    // Load a csv file.
    f, _ := os.Open(filePath)
    // Create a new reader.
	r := csv.NewReader(f)
	
	link = List{}

    for {
        record, err := r.Read()
        // Stop at EOF.
        if err == io.EOF {
            break
        }

        if err != nil {
            panic(err)
        }
		if date == record[4] || region == record[9] {
			link.Insert(record[4], record[2], record[3], record[6] , record[10], record[5] , record[9])
		}

		
	}
	return

}


func main() {
	go server()
	var input string
	fmt.Scanln(&input)
  }

