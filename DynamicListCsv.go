package main

// If running on linux and wishing to build a windows executable run this:
// GOOS=windows GOARCH=386 go build -o DynamicListCsv.exe DynamicListCsv.go
//
// If running on windows and wishing to build a windows executable run this:
// set GOOS=linux
// set GOARCH=amd64 
// go build -o DynamicListCsv DynamicListCsv.go

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type OrgDetail struct {
	UserNamePrefix string
	OrgDomainInfo  OrgDomainInfo
	UserCount      int
	StartingIndex  int
	ZeroPadded     bool
	Password       string
}

type OrgDomainInfo struct {
	EmailDomainPrefix string
	EmailDomainSuffix string
	Count             int
	StartingIndex     int
	ZeroPadded        bool
}

type OrgInfo struct {
	UserNamePrefix string
	OrgEmailDomain string
	UserCount      int
	StartingIndex  int
	ZeroPadded     bool
	Password       string
}

const INPUTFILE = "config.json"
const OUTPUTFILE = "output.csv"

func main() {
    go beginProcess()
	fmt.Println("Press enter to exit")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func beginProcess(){    
    orgDetails := getOrgDetails()
    orgInfo := getOrgInfoFromOrgDetails(orgDetails)
    outputCsv := createCsv(orgInfo)
    go writeCsv(outputCsv)
}

func getOrgDetails() []OrgDetail {
	data, _ := ioutil.ReadFile(INPUTFILE)
	var orgDetails []OrgDetail
	json.Unmarshal(data, &orgDetails)
	return orgDetails
}

func getOrgInfoFromOrgDetails(orgDetails []OrgDetail) []OrgInfo {
	var orgInfos []OrgInfo
	for _, orgDetail := range orgDetails {
        if orgDetail.OrgDomainInfo.Count > 0 {
            for i := 0; i < orgDetail.OrgDomainInfo.Count; i++ {
                newOrgInfo := getNewOrgInfo(orgDetail, i)
                orgInfos = append(orgInfos, newOrgInfo)
            }
        } else {
            newOrgInfo := getNewOrgInfo(orgDetail, 0)
            orgInfos = append(orgInfos, newOrgInfo)
        }
	}
	return orgInfos
}

func getNewOrgInfo(orgDetail OrgDetail, index int) OrgInfo {
	orgInfo := OrgInfo{
		orgDetail.UserNamePrefix,
		getOrgEmailDomain(orgDetail.OrgDomainInfo, index),
		orgDetail.UserCount,
		orgDetail.StartingIndex,
		orgDetail.ZeroPadded,
		orgDetail.Password}
	return orgInfo
}

func getOrgEmailDomain(orgDomainInfo OrgDomainInfo, index int) string {
	var orgNumber string
	if orgDomainInfo.Count == 0 {
		orgNumber = ""
	} else {
		orgNumber = getOrgNumber(orgDomainInfo, index)
	}
	orgEmailDomain := orgDomainInfo.EmailDomainPrefix + orgNumber + orgDomainInfo.EmailDomainSuffix
	return orgEmailDomain
}

func getOrgNumber(orgDomainInfo OrgDomainInfo, index int) string {
	var orgNumber string
	if orgDomainInfo.ZeroPadded {
        maxIndexLength := len([]rune(strconv.Itoa(orgDomainInfo.StartingIndex + orgDomainInfo.Count)))
        currentIndex := strconv.Itoa(index + orgDomainInfo.StartingIndex)
        currentIndexLength := len([]rune(currentIndex))
        for currentIndexLength < maxIndexLength {
            currentIndex = "0" + currentIndex
            currentIndexLength = len([]rune(currentIndex))
        }
		orgNumber = currentIndex
	} else {
		orgNumber = strconv.Itoa(index + orgDomainInfo.StartingIndex)
	}
	return orgNumber
}

func createCsv(orgs []OrgInfo) [][]string {
	var outputCsv [][]string
	currentUserIndices := getUserStartingIndices(orgs)
	currentOrgIndex := 0
	for i := 0; i < getTotalUserCount(orgs); i++ {
		newRow := getNewRow(orgs, currentOrgIndex, currentUserIndices[currentOrgIndex])
		outputCsv = append(outputCsv, newRow)
		currentUserIndices[currentOrgIndex]++
		currentOrgIndex = getNextOrgIndex(orgs, currentOrgIndex, currentUserIndices)
	}
	return outputCsv
}

func getUserStartingIndices(orgs []OrgInfo) []int {
	currentUserIndices := make([]int, len(orgs))
	for index, org := range orgs {
		currentUserIndices[index] = org.StartingIndex
	}
	return currentUserIndices
}

func getTotalUserCount(orgs []OrgInfo) int {
	totalUserCount := 0
	for _, org := range orgs {
		totalUserCount += org.UserCount
	}
	return totalUserCount
}

func getNewRow(orgs []OrgInfo, orgIndex int, userIndex int) []string {
	org := orgs[orgIndex]
	newRow := make([]string, 2)
	newRow[0] = org.UserNamePrefix + getUserNumber(org, userIndex) + "@" + org.OrgEmailDomain
	newRow[1] = org.Password
	return newRow
}

func getUserNumber(org OrgInfo, userIndex int) string {
	var userNumber string
	if org.ZeroPadded {
        maxIndexLength := len([]rune(strconv.Itoa(org.StartingIndex + org.UserCount)))
        currentIndex := strconv.Itoa(userIndex + org.StartingIndex)
        currentIndexLength := len([]rune(currentIndex))
        for currentIndexLength < maxIndexLength {
            currentIndex = "0" + currentIndex
            currentIndexLength = len([]rune(currentIndex))
        }
        userNumber = currentIndex
	} else {
		userNumber = strconv.Itoa(userIndex)
	}
	return userNumber
}

func getNextOrgIndex(orgs []OrgInfo, currentOrgIndex int, currentUserIndices []int) int {
	newIndex := currentOrgIndex + 1
	if newIndex >= len(orgs) {
		newIndex = 0
	}
	if currentUserIndices[newIndex] > orgs[newIndex].UserCount+orgs[newIndex].StartingIndex {
		newIndex = getNextOrgIndex(orgs, newIndex, currentUserIndices)
	}
	return newIndex
}

func writeCsv(outputCsv [][]string) {
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
