# BUCKET DB

<div style="text-align: center"><img src="https://www.seekpng.com/png/full/28-281861_metal-bucket-png-clipart-bucket.png" alt="A metal bucket" height=200 width=400></div>

## Summary

A minimal file based database system that stores all data as lines in a .txt file. The system works off of storage files known as buckets.

## Build

Bucket Database is completely built with standard gin libraries.

## How to Start

- Create a new bucket with <span style="background-color: yellow; color: black;">[ dbObject := newDefaultBucket("normal", "new bucket") ]</span> which returns a pointer
- Create a new dbClient object with <span style="background-color: yellow; color: black;">[ newClient := &dbClient{ user: "admin", pass: "admin",} ]</span> which returns a pointer
- Add your first item to the database with <span style="background-color: yellow; color: black;">[ dbObject.writeToStorage("Hello World", newClient) ]</span>
- Check out your new DB with <span style="background-color: yellow; color: black;">[ returnList2 := dbObject.readBucket(newClient) ]</span> which returns a slice of strings
- Print the contents of the new DB with <span style="background-color: yellow; color: black;">[ fmt.Println(returnlist2) ]</span>

### Notes and Final Comments

- Deleting lines is not functional and just empties the entire database.

- New items are added to the bottom of the bucket/file.
