package main

import (
        "fmt"
        "io/ioutil"
        "net/http"
        "regexp"
        "encoding/json"
        "flag"
        "os"
)

func main() {

        stationID := flag.String("s", "", "Station ID - see https://www.dwd.de/DE/leistungen/opendata/help/warnungen/cap_warncellids_csv.csv")
        flag.Parse()

        if *stationID == "" {
                ExitUnknown("No station id given")
        }
        if len(*stationID) != 9 {
                ExitUnknown("Given station ID doesn't exist!")
        }

        // Perform http request and convert byte array to string
        http_response := string(http_request())

        // Convert JSONP to JSON
        re := regexp.MustCompile(`warnWetter.loadWarnings\(`)
        s := re.ReplaceAllString(http_response, "")
        s = s[:len(s)-2]

        // Declared an empty map interface
        var main_hash map[string]interface{}

        // Unmarshal or Decode the JSON to the interface.
        json.Unmarshal([]byte(s), &main_hash)

        // check if a warning for stationID is available
        if _, ok := main_hash["warnings"].(map[string]interface{})[*stationID]; ok {

                level := fmt.Sprintf("%v", main_hash["warnings"].(map[string]interface{})[*stationID].([]interface{})[0].(map[string]interface{})["level"])

                if (level == "1" || level == "2" || level == "3") {
                        ExitWarning(getData(main_hash, *stationID))
                } else if (level == "4") {
                        ExitCritical(getData(main_hash, *stationID))
                } else {
                        ExitWarning(getData(main_hash, *stationID))
                }
        } else {
                ExitOK()
        }
}

func http_request() (arr_resp []byte){
        client := &http.Client{
                CheckRedirect: func(req *http.Request, via []*http.Request) error {
                return http.ErrUseLastResponse
                },
        }
        resp, err := client.Get("https://www.dwd.de/DWD/warnungen/warnapp/json/warnings.json")
        if (err != nil) {
                ExitUnknown(err.Error())
        }
        if (resp.StatusCode != 200) {
                comb := fmt.Sprint(resp.StatusCode)
                ExitUnknown(comb)
        }
        arr_resp, err = ioutil.ReadAll(resp.Body)
        if (err != nil) {
                ExitUnknown(err.Error())
        }
        return
}

func ExitUnknown (reason string) {
        fmt.Printf("%s %s", "UNKNOWN -", reason)
        os.Exit(3)
}
func ExitOK () {
        fmt.Println("OK - No warning found")
        os.Exit(0)
}
func ExitCritical (headline string, description string) {
        fmt.Printf("%s %s %s %s", "CRITICAL -", headline, "\n", description)
        os.Exit(2)
}
func ExitWarning (headline string, description string) {
        fmt.Printf("%s %s %s %s", "WARNING -", headline, "\n", description)
        os.Exit(1)
}
func getData (main_hash map[string]interface{}, stationID string) (headline string, description string) {
        headline = fmt.Sprintf("%v", main_hash["warnings"].(map[string]interface{})[stationID].([]interface{})[0].(map[string]interface{})["headline"])
        description = fmt.Sprintf("%v", main_hash["warnings"].(map[string]interface{})[stationID].([]interface{})[0].(map[string]interface{})["description"])
        return
}
