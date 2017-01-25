package main

// GOOS=windows GOARCH=386 go build -o CsvCombiner.exe CsvCombiner.go

import (
	"flag"
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

type OrgDetails struct {
    UserNamePrefix string
    OrgEmailDomain string
	UserCount int
	StartingIndex int
	ZeroPadded bool
	Password string
}

const OUTPUTFILE = "test.csv"

var configFile string
var orgs []OrgDetails

func main(){
	getInput()
	getOrgDetails()
	outputCsv = createCsv()
	writeCsv(outputCsv)
    
    fmt.Println("Press enter to exit")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}


func getInput() {
	configFilePtr := flag.String("config", "config.json", "")
	flag.Parse()
	configFile = *configFilePtr
}

func getOrgDetails(){
    //TODO read from config.json
    orgs = make([]OrgDetails, 12)
    for i := 0; i < len(orgs); i++ {
        orgNum := i+1;
        org[i] = OrgDetails{
            "user",
            "testacdload"+orgNum+".test",
            2000,
            0,
            false,
            "test1234"
        }
    }
}

func createCsv() [][]string {
    var outputCsv [][]string
    currentUserIndices := getUserStartingIndices()
    currentOrgIndex := 0
    for i := 0; i < getTotalUserCount(); i++{
        newRow = getNewRow(currentOrgIndex, currentUserIndices[currentOrgIndex])
        outputCsv = append(outputCsv, newRow)
        currentUserIndices[currentOrgIndex]++
        currentOrgIndex = getNextOrgIndex(currentOrgIndex, currentUserIndices)
    }
    return outputCsv
}

func getNextOrgIndex(currentOrgIndex int, currentUserIndices []int) int {
    newIndex := currentOrgIndex + 1
    if newIndex >= len(orgs) {
        newIndex = 0
    }
    if currentUserIndices[newIndex] > orgs[newIndex].UserCount + orgs[newIndex].StartingIndex {
        newIndex = getNextOrgIndex(newIndex, currentUserIndices)
    }
    return newIndex
}

func getNewRow(orgIndex int, userIndex int) []string {
    org := orgs[orgIndex]
    newRow := make([]string, 2)
    newRow[0] = org.UserNamePrefix + getUserNumber(org, userIndex) + "@" + org.OrgEmailDomain
    newRow[1] = org.Password
    return newRow
}

func getUserNumber(org OrgDetails, userIndex int) string {
    var userNumber string
    if (org.ZeroPadded){
        //TODO pad with zeros
        userNumber = userIndex
    } else {
        userNumber = userIndex
    }
    return userNumber
}

func getUserStartingIndices() []int {
    currentUserIndices = make([]int, len(orgs))
    for index, org := range orgs {
        currentUserIndices[index] = org.StartingUserIndex
    }
    return currentUserIndices
}

func getTotalUserCount() int {
    totalUserCount := 0
    for _, org := range orgs {
        totalUserCount += org.UserCount
    }
    return totalUserCount
}

func writeCsv(outputCsv [][]string){
    fmt.Println("Writing output file: " + OUTPUTFILE)
	outputFile, err := os.Create(OUTPUTFILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()
	writer := csv.NewWriter(outputFile)
	writer.WriteAll(outputCsv)
}