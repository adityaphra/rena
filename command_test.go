package main

import (
	"fmt"
	"testing"
)

func TestCommand(t *testing.T) {
	f := CreateFile("/home/user/Better Call Saul - S0E01.mkv")
	filePath := f.getFullPath()
	testCases := map[string]string{
		"s/saul/Saul Goodman":                   "/home/user/Better Call Saul Goodman - S0E01.mkv",
		"s/saul/Saul Goodman/m":                 "/home/user/Better Call Saul - S0E01.mkv",
		"s/sa(ul)/%0/r":                         "/home/user/Better Call Saul - S0E01.mkv",
		"s/sa(ul)/%1/r":                         "/home/user/Better Call ul - S0E01.mkv",
		"s/sa(ul)/%1/rm":                        "/home/user/Better Call Saul - S0E01.mkv",
		"s/(sa)(ul)/%2%1/r":                     "/home/user/Better Call ulSa - S0E01.mkv",
		"s/sa(?P<first>ul)/%first/r":            "/home/user/Better Call ul - S0E01.mkv",
		"s/(saul|call)//rw":                     "/home/user/Better   - S0E01.mkv",
		"s/(saul|call)//r":                      "/home/user/Better - S0E01.mkv",
		"d/(saul|call)/rw":                      "/home/user/Better   - S0E01.mkv",
		"d/(saul|call)/r":                       "/home/user/Better - S0E01.mkv",
		"d/saul ":                               "/home/user/Better Call - S0E01.mkv",
		"d/(call|saul)/r":                       "/home/user/Better - S0E01.mkv",
		"d/saul /m":                             "/home/user/Better Call Saul - S0E01.mkv",
		"d/sa.l /r":                             "/home/user/Better Call - S0E01.mkv",
		"d/sa.l /rm":                            "/home/user/Better Call Saul - S0E01.mkv",
		"d/Sa.l /rm":                            "/home/user/Better Call - S0E01.mkv",
		"t/%n.mp4":                              "/home/user/Better Call Saul - S0E01.mp4",
		"t/%%n.mp4":                             "/home/user/%n.mp4",
		"t/%n.%x":                               "/home/user/Better Call Saul - S0E01.mkv",
		"t/[prefix] %f":                         "/home/user/[prefix] Better Call Saul - S0E01.mkv",
		"t/%n [suffix].%x":                      "/home/user/Better Call Saul - S0E01 [suffix].mkv",
		"m,saul,/home/user/Better Call Saul":    "/home/user/Better Call Saul/Better Call Saul - S0E01.mkv",
		"m,sa.l,/home/user/Better Call Saul,r":  "/home/user/Better Call Saul/Better Call Saul - S0E01.mkv",
		"m,Sa.l,/home/user/Better Call Saul,rm": "/home/user/Better Call Saul/Better Call Saul - S0E01.mkv",
		"m,saul,/home/user/Better Call Saul,m":  "/home/user/Better Call Saul - S0E01.mkv",
		"m,Saul,/home/user/Better Call Saul,m":  "/home/user/Better Call Saul/Better Call Saul - S0E01.mkv",
		"m,(.*) - .*,/home/user/%1,r":           "/home/user/Better Call Saul/Better Call Saul - S0E01.mkv",
		"m,(?P<first>Call) (?P<second>Saul),/home/user/%second-%first,r": "/home/user/Saul-Call/Better Call Saul - S0E01.mkv",
	}

	passCount := 0
	failCount := 0
	for command, expected := range testCases {
		c, err := parseCommand(command)
		if err != nil {
			t.Error(err)
		}

		copiedFile := f
		c.Execute(&copiedFile)
		copiedFilePath := copiedFile.getFullPath()
		if copiedFilePath != expected {
			t.Errorf(`---------------------------------
Input:      (%v, %v)
Expecting:  %v
Actual:     %v
Fail
`, command, filePath, expected, copiedFilePath)
			failCount++
		} else {
			fmt.Printf(`---------------------------------
Input:      (%v, %v)
Expecting:  %v
Actual:     %v
Pass
`, command, filePath, expected, copiedFilePath)
			passCount++
		}

	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}
