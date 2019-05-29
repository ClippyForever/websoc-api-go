package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//XMLresult is Go
type XMLresult struct {
	XMLName xml.Name `xml:"websoc_results"`
	GetP    string   `xml:"get_parm,attr"`
	Schools schools  `xml:"course_list"`
}

type schools struct {
	Schools []school `xml:"school"`
}

type school struct {
	Name        string       `xml:"school_name,attr"`
	Comment     string       `xml:"school_comment"`
	School_code string       `xml:"school_code,attr"`
	Departments []Department `xml:"department"`
}

type Department struct {
	XMLName                     xml.Name `xml:"department"`
	Name                        string   `xml:"dept_name,attr"`
	Dept_Code                   string   `xml:"dept_code,attr"`
	Department_Comment          string   `xml:"department_comment"`
	Course_Number_Range_Comment string   `xml:"course_number_range_comment"`
	Course_Code_Range_Comment   string   `xml:"course_code_range_comment"`
	Courses                     []Course `xml:"course"`
}

type Course struct {
	XMLName          xml.Name  `xml:"course"`
	Course_Number    string    `xml:"course_number,attr"`
	Course_Title     string    `xml:"course_title,attr"`
	PrerequisiteLink string    `xml:"course_prereq_link"`
	Comment          string    `xml:"course_comment"`
	Sections         []Section `xml:"section"`
}

type Section struct {
	XMLName        xml.Name       `xml:"section"`
	ClassCode      string         `xml:"course_code"`
	ClassType      string         `xml:"sec_type"`
	SectionCode    string         `xml:"sec_num"`
	Units          string         `xml:"sec_units"`
	FinalExam      FinalExam      `xml:"sec_final"`
	Restrictions   string         `xml:"sec_restrictions"`
	Status         string         `xml:"sec_status"`
	Comment        string         `xml:"sec_comment"`
	Sec_Instructor Instructors    `xml:"sec_instructor"`
	Sec_Enrollment Sec_Enrollment `xml:"sec_enrollment"`
	Sec_Meeting    Sec_Meeting    `xml:"sec_meeting"`
}

type FinalExam struct {
	Sec_Final_date string `xml:"sec_final_date"`
	Sec_Final_day  string `xml:"sec_final_day"`
	Sec_final_time string `xml:"sec_final_time"`
}

type Instructors struct {
	XMLName     xml.Name `xml:"sec_instructor"`
	Instructors []string `xml:"instructor"`
}

type Sec_Enrollment struct {
	XMLName              xml.Name `xml:"sec_enrollment"`
	MaxCapacity          string   `xml:"sec_max_enroll"`
	NumCurrentlyEnrolled string   `xml:"sec_enrolled"`
	NumOnWaitlist        string   `xml:"sec_waitlist"`
	NumRequested         string   `xml:"sec_enroll_requests"`
	NumNewOnlyReserved   string   `xml:"sec_new_only_reserved"`
}

type Sec_Meeting struct {
	XMLName  xml.Name  `xml:"sec_meeting"`
	Meetings []Meeting `xml:"sec_meet"`
}

type Meeting struct {
	XMLname       xml.Name `xml:"sec_meet"`
	Sec_Days      string   `xml:"sec_days"`
	Sec_Time      string   `xml:"sec_time"`
	Sec_Bldg      string   `xml:"sec_bldg"`
	Sec_Room      string   `xml:"sec_room"`
	Sec_Room_Link string   `xml:"sec_room_link"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	formData := url.Values{
		"Submit":           {"Display XML Results"},
		"YearTerm":         {getCodedTerm(strings.ToLower(request.QueryStringParameters["term"]))},
		"ShowComments":     {"on"},
		"ShowFinals":       {"on"},
		"Breadth":          {checkExistParams(request.QueryStringParameters["GE"], "ANY")},
		"Dept":             {checkExistParams(request.QueryStringParameters["department"], "ALL")},
		"CourseNum":        {request.QueryStringParameters["courseNum"]},
		"Division":         {checkExistParams(getCodedDiv(request.QueryStringParameters["division"]), "ANY")},
		"CourseCodes":      {request.QueryStringParameters["courseCodes"]},
		"InstrName":        {request.QueryStringParameters["instructorName"]},
		"CourseTitle":      {request.QueryStringParameters["courseTitle"]},
		"ClassType":        {checkExistParams(request.QueryStringParameters["courseCodes"], "ALL")},
		"Units":            {request.QueryStringParameters["units"]},
		"Days":             {request.QueryStringParameters["days"]},
		"StartTime":        {request.QueryStringParameters["startTime"]},
		"EndTime":          {request.QueryStringParameters["endTimes"]},
		"MaxCap":           {request.QueryStringParameters["maxCap"]},
		"FullCourses":      {checkExistParams(request.QueryStringParameters["fullCourses"], "ANY")},
		"CancelledCourses": {checkExistParams(request.QueryStringParameters["cancelledCourses"], "EXCLUDE")},
		"Bldg":             {request.QueryStringParameters["building"]},
		"Room":             {request.QueryStringParameters["room"]},
	}
	resp, err := http.PostForm("https://www.reg.uci.edu/perl/WebSoc", formData)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var results XMLresult

	xml.Unmarshal(body, &results)

	ok, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}

	return events.APIGatewayProxyResponse{Body: string(ok), Headers: request.QueryStringParameters, StatusCode: 200}, nil
}

func checkExistParams(exist string, notExist string) string {
	if exist == "" {
		return notExist
	}

	return exist
}

//uncomment this function to run local
// func getResult(w http.ResponseWriter, r *http.Request) {
// 	result := MakeRequest()
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(result.Schools.Schools)
// }

func main() {
	//comment 3 linese below to test local
	lambda.Start(handleRequest)

	formData := url.Values{"dsdsd": {"dsdsds"}}
	fmt.Println(formData.Get("dsdssd"))

	//uncomment  below to test local
	// r := mux.NewRouter()
	// r.HandleFunc("/api/", getResult).Methods("GET")

	// log.Fatal(http.ListenAndServe(":8000", r))

	// fmt.Println("Hello")

}

func getCodedTerm(term string) string {
	actualTerm := ""
	splittedTerm := strings.Split(term, " ")

	if splittedTerm[1] == "fall" {
		actualTerm = splittedTerm[0] + "-92"
	} else if splittedTerm[1] == "winter" {
		actualTerm = splittedTerm[0] + "-03"
	} else if splittedTerm[1] == "spring" {
		actualTerm = splittedTerm[0] + "-14"
	} else if splittedTerm[1] == "summer1" {
		actualTerm = splittedTerm[0] + "-25"
	} else if splittedTerm[1] == "summer2" {
		actualTerm = splittedTerm[0] + "-76"
	} else if splittedTerm[1] == "summer10wk" {
		actualTerm = splittedTerm[0] + "-39"
	}

	return actualTerm
}
func getCodedDiv(div string) string {
	codedDiv := strings.ToLower(div)

	if codedDiv == "all" {
		codedDiv = "all"
	} else if codedDiv == "lowerdiv" {
		codedDiv = "0xx"
	} else if codedDiv == "upperdiv" {
		codedDiv = "1xx"
	} else if codedDiv == "graduate" {
		codedDiv = "2xx"
	}

	return codedDiv
}

// func MakeRequest() XMLresult {

// 	formData := url.Values{
// 		"Submit":           {"Display XML Results"},
// 		"YearTerm":         {getCodedTerm("2019 fall")},
// 		"ShowComments":     {"on"},
// 		"ShowFinals":       {"on"},
// 		"Breadth":          {"ANY"},
// 		"Dept":             {"BIO SCI"},
// 		"CourseNum":        {""},
// 		"Division":         {"ANY"},
// 		"CourseCodes":      {""},
// 		"InstrName":        {""},
// 		"CourseTitle":      {""},
// 		"ClassType":        {"ALL"},
// 		"Units":            {""},
// 		"Days":             {""},
// 		"StartTime":        {""},
// 		"EndTime":          {""},
// 		"MaxCap":           {""},
// 		"FullCourses":      {"ANY"},
// 		"CancelledCourses": {"EXCLUDE"},
// 		"Bldg":             {""},
// 		"Room":             {""},
// 	}

// 	resp, err := http.PostForm("https://www.reg.uci.edu/perl/WebSoc", formData)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)

// 	var results XMLresult

// 	xml.Unmarshal(body, &results)

// 	return results
// }
