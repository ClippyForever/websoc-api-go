
# Welcome to My Go!



# Example

Go Here [Go](https://971k1pm7de.execute-api.us-west-1.amazonaws.com/dev/websoc?department=COMPSCI&term=2019%20Fall&GE=ANY&courseNum=&courseCodes=&instructorName=&units=&endTime=&startTime=&fullCourses=ANY&building=), .

# List of parameters:

*department
*term
*GE
*courseNum
*courseCodes
*instructorName
*units
*endTime
*startTime
*fullCourses
*building

Full Link: https://971k1pm7de.execute-api.us-west-1.amazonaws.com/dev/websoc?department=COMPSCI&term=2019%20Fall&GE=ANY&courseNum=&courseCodes=&instructorName=&units=&endTime=&startTime=&fullCourses=ANY&building=

# Usage

```Go
 func getResult(w http.ResponseWriter, r *http.Request) {

 result := MakeRequest()

 w.Header().Set("Content-Type", "application/json")

 json.NewEncoder(w).Encode(result.Schools.Schools)

}

func main() {

 r := mux.NewRouter()

 r.HandleFunc("/api/", getResult).Methods("GET")
//run on local
 log.Fatal(http.ListenAndServe(":8000", r))

 fmt.Println("Hello")

}
```
# To Run On Local 

Get Mux router library first 
Type in and enter "go get -u github.com/gorilla/mux" in the terminal

Remember to import Mux library
```Go
import (
"github.com/gorilla/mux"
)
```


Two options to run:
1. Type in and enter  "go run main.go" to run without compiling the program first.
2. Type in and enter "go build" then "./nameOftheProgram" to run with compiling.
