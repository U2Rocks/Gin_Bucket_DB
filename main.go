package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

// BUCKETDB - A simple go database that stores items in a "bucket"
// figure out how to remove items later

// simple databases -> store data(G) -> retrieve data(G) -> update data -> delete data
// make secure version -> pass trusted struct that has username and pass into functions that get and read data

// implement reading file in reverse for latest writes first
// add another read function that filters results based on text input(test it now)
// add ability to remove item from top or bottom***

// _ = normal / * = skip / @ = removed or delete on next cycle / ! = only read if security variable = true
// add special mode to line writer functions -> add conditions to reading functions that do not put line into slice
// add special line to top of created bucket that has name and date of creation

// add implementation to delete current db at path and test if delete lines works

// STRUCTS

// admin user info for db object
type adminUser struct {
	user string
	pass string
}

// db object that holds the mode, path, and admin information
type db struct {
	mode       string
	admin      adminUser
	databaseID string
	path       string
}

// client user object for communicating with a db
type dbClient struct {
	user string
	pass string
}

// IMPORTANT STARTING FUNCTIONS

// function to create new database with specific admin settings
func newBucket(setMode string, setAdmin adminUser, ID string, newPath string) *db {
	newDatabase := db{
		mode:       setMode,
		admin:      setAdmin,
		databaseID: ID,
		path:       newPath,
	}
	return &newDatabase
}

// function to create new database with default admin settings
func newDefaultBucket(setMode string, ID string) *db {
	defaultAdmin := adminUser{
		user: "admin",
		pass: "admin",
	}
	newDatabase := db{
		mode:       setMode,
		admin:      defaultAdmin,
		databaseID: ID,
		path:       "./storage.txt",
	}
	return &newDatabase
}

// METHODS TO GET DB OBJECT INFO

// function that returns string for current db mode
func (d *db) getBucketMode() string {
	return d.mode
}

// function that returns db admin information in a slice
func (d *db) getBucketAdmin() []string {
	user := d.admin.user
	pass := d.admin.pass
	loginSlice := []string{user, pass}
	return loginSlice
}

// function that returns a databases ID value
func (d *db) getBucketId() string {
	return d.databaseID
}

// function that returns a databases path
func (d *db) getBucketPath() string {
	return d.path
}

// METHODS TO SET DB INFO

// function that sets string for current db mode(valid options are "normal" and "secure")
func (d *db) setBucketMode(newMode string) {
	d.mode = newMode
}

// function that sets db admin information in a slice [user, pass]
func (d *db) setBucketAdmin(userInfo []string) {
	d.admin.user = userInfo[0]
	d.admin.pass = userInfo[0]
}

// function that sets a databases ID value
func (d *db) setBucketId(newID string) {
	d.databaseID = newID
}

// function that returns a databases path
func (d *db) setBucketPath(newPath string) {
	d.path = newPath
}

// MAIN FUNCTION

func main() {
	// printlns to make conosle more readable
	fmt.Println()
	fmt.Println("<----------------->")
	fmt.Println("<-------NEW RUN...------->")
	fmt.Println("<----------------->")
	fmt.Println()

	// create new base client
	newClient := &dbClient{
		user: "admin",
		pass: "admin",
	}

	dbObject := newDefaultBucket("normal", "first")
	fmt.Println("@@@ DB OBJECT @@@")
	fmt.Println(dbObject.getBucketId())
	fmt.Println(dbObject.getBucketAdmin())
	fmt.Println(dbObject.getBucketMode())
	fmt.Println(dbObject.getBucketPath())
	fmt.Println("@@@ DB OBJECT END @@@")

	// declare content to write and file path and security bool
	writeContent := "this might be deleted"
	filterString := "might"

	// run meta function for writing raw data into storage file
	dbObject.writeToStorage(writeContent, newClient)

	// read file and print contents of file
	dbObject.writeToStorage(filterString, newClient)
	dbObject.writeToStorage(filterString, newClient)
	dbObject.writeToStorage(filterString, newClient)
	dbObject.writeToStorage(filterString, newClient)
	dbObject.writeToStorage("this is fine", newClient)

	// delete lines that mention deleted
	dbObject.deleteLines(newClient, filterString)

	// read new db
	returnList2 := dbObject.readBucketFiltered(newClient, filterString)
	fmt.Println(returnList2)

}

// function to reduce code and handle error printing logic
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// meta function that checks path and runs correct write functionx
func (d *db) writeToStorage(writeContent string, client *dbClient) {
	// check if file path is valid
	isPathGood := checkFilePath(d.path)

	// check for secure mode and act normally otherwise
	if d.mode == "secure" {
		if d.admin.user == client.user && d.admin.pass == client.pass {
			// run the appropriate function depending on if the file exists
			if isPathGood {
				writeToExisitingFile(d.path, writeContent, d.mode)
			} else {
				writeToNewFile(d.path, writeContent, d.mode, d.databaseID)
			}
		}
	} else {
		if isPathGood {
			writeToExisitingFile(d.path, writeContent, d.mode)
		} else {
			writeToNewFile(d.path, writeContent, d.mode, d.databaseID)
		}
	}
}

// function to read the storage.txt file line by line and return a slice with all values
func (d *db) readBucket(c *dbClient) []string {
	fmt.Println()
	file, err := os.Open(d.path)
	handleError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanList := []string{}
	for scanner.Scan() {
		if string(scanner.Text()[0]) == "!" {
			if d.admin.user == c.user && d.admin.pass == c.pass {
				scanList = append(scanList, scanner.Text()[1:])
			}
		} else if string(scanner.Text()[0]) == "*" {
			continue
		} else {
			scanList = append(scanList, scanner.Text()[1:])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return scanList

}

// reads storage line by line and returns a string slice with all values but applies a filter based on input
func (d *db) readBucketFiltered(c *dbClient, filter string) []string {
	fmt.Println()
	file, err := os.Open(d.path)
	handleError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanList := []string{}
	for scanner.Scan() {
		if string(scanner.Text()[0]) == "!" {
			if d.admin.user == c.user && d.admin.pass == c.pass {
				if strings.Contains(string(scanner.Text()[1:]), filter) {
					scanList = append(scanList, scanner.Text()[1:])
				}
			}
		} else if string(scanner.Text()[0]) == "*" {
			continue
		} else {
			if strings.Contains(string(scanner.Text()[1:]), filter) {
				scanList = append(scanList, scanner.Text()[1:])
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return scanList

}

// function to check if the file path is valid with console logging sugar
func checkFilePath(filePath string) bool {
	fmt.Println("<----------------->")
	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("The file path: ", filePath, " is good")
		fmt.Println("<----------------->")
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("The file path: ", filePath, " is bad...")
		fmt.Println("<----------------->")
		return false

	} else {
		fmt.Println("The file path: ", filePath, " may or may not exist...")
		fmt.Println("<----------------->")
		return false
	}
}

// WRITE FUNCTIONS

// function to write data to an existing storage.txt file
func writeToExisitingFile(inputPath string, writeContent string, mode string) {
	switch mode {
	case "normal":
		writeContent = "_" + writeContent + "\n"
	case "secure":
		writeContent = "!" + writeContent + "\n"
	default:
		writeContent = "_" + writeContent + "\n"
	}
	file, err := os.OpenFile(inputPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	handleError(err)

	length, err := io.WriteString(file, writeContent)

	// this logic forces length to stop erroring and will not print
	noError := false
	if noError {
		fmt.Print(length)
	}

	handleError(err)
	defer file.Close()
}

// function to write data to a nonexistent storage.txt file
func writeToNewFile(inputPath string, writeContent string, mode string, ID string) {
	switch mode {
	case "normal":
		writeContent = "_" + writeContent + "\n"
	case "secure":
		writeContent = "!" + writeContent + "\n"
	default:
		writeContent = "_" + writeContent + "\n"
	}
	file, err := os.Create(inputPath)
	handleError(err)

	currentTime := time.Now()

	// header message for new bucket
	headerMessage := "*** Bucket ID: \"" + ID + "\"; Creation Date: " + currentTime.String() + " ***\n"

	length, err := io.WriteString(file, headerMessage)
	handleError(err)

	// this logic forces length to stop erroring and will not print
	noError := false
	if noError {
		fmt.Print(length)
	}

	length2, err := io.WriteString(file, writeContent)
	handleError(err)

	// this logic forces length to stop erroring and will not print
	noError2 := false
	if noError2 {
		fmt.Print(length2)
	}

	defer file.Close()
}

// creates a line in the database that will be ignored by the reader
func (d *db) writeSpecialLine(writeContent string) {
	writeContent = "*" + writeContent + "\n"

	file, err := os.OpenFile(d.path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	handleError(err)

	length, err := io.WriteString(file, writeContent)
	handleError(err)

	// this logic forces length to stop erroring and will not print
	noError := false
	if noError {
		fmt.Print(length)
	}

	defer file.Close()

}

// function filters out line based on a pattern and copies all non matching lines to a new database
func (d *db) deleteLines(c *dbClient, pattern string) {
	// check if pattern is empty
	if pattern == "" {
		return
	}
	// only delete line if path is valid
	if checkFilePath(d.path) {
		// create slice to hold values for new storage bucket...
		newBucketValues := []string{}
		// create a regex pattern
		r, _ := regexp.Compile(pattern)
		// open the file with read only permissions
		file, err := os.Open(d.path)
		handleError(err)
		defer file.Close()
		// create a new scanner object with the opened file as an argument
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			if string(scanner.Text()[0]) == "*" {
				continue
			} else if string(scanner.Text()[0]) == "!" {
				if d.admin.user == c.user && d.admin.pass == c.pass {
					lineValue := string(scanner.Text())
					if r.MatchString(lineValue) == false {
						newBucketValues = append(newBucketValues, lineValue)
					}
				}
			} else {
				if strings.Contains(string(scanner.Text()), pattern) {
					lineValue := string(scanner.Text())
					if r.MatchString(lineValue) == false {
						newBucketValues = append(newBucketValues, lineValue)
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		// delete current storage db(not implemented)
		deleteBucket(d.path)
		// create new database
		createEmptyBucket(d.path, d.databaseID)
		// input new Bucket list to new database(not implemented)
		copyBucket(newBucketValues, d.path)
	}
}

// DELETE FUNCTIONS - CREATE EMPTY, DELETE BUCKET, COPY LIST

// create an empty bucket/db that only contains a header row
func createEmptyBucket(inputPath string, ID string) {
	// create new file at specific path
	file, err := os.Create(inputPath)
	handleError(err)

	// get current time object
	currentTime := time.Now()

	// header message for new bucket(no write mode required)
	headerMessage := "*** Bucket ID: \"" + ID + "\"; Creation Date: " + currentTime.String() + " ***\n"

	// write header message to top of bucket
	length, err := io.WriteString(file, headerMessage)
	handleError(err)

	// this logic forces length to stop erroring and will not print
	noError := false
	if noError {
		fmt.Print(length)
	}

	defer file.Close()
}

// delete the current bucket with the path found in the db object
func deleteBucket(inputPath string) {
	err := os.Remove("storage.txt")
	handleError(err)
}

// copy all values from a slice of strings into the file at the input path
func copyBucket(rowList []string, inputPath string) {

	file, err := os.OpenFile(inputPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	handleError(err)

	for _, value := range rowList {

		length, err := io.WriteString(file, value)
		handleError(err)

		// this logic forces length to stop erroring and will not print
		noError := false
		if noError {
			fmt.Print(length)
		}
	}

	defer file.Close()
}
