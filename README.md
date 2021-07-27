# IITK Coin
## SnT Project 2021, Programming Club 

This repository contains the code for the IITK Coin project done so far.

### Relevant Links

- [Midterm Evaluation presentation](https://docs.google.com/presentation/d/1kriN-7A3v1RlXUDL5NETX3roJKRMJInptkWofIxY8dg/edit?usp=sharing)
- [Midterm Documentation](https://docs.google.com/document/d/1bvOWH4k0U-l2pQ1jLWIDzOkJ2wbHNW4jJw7tMWkUV6o/edit?usp=sharing)

## Table Of Content
- [Development Environment](#development-environment)
- [Directory Structure](#directory-structure)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Models](#models)

## Development Environment

```bash
- go version: go1.16.4 linux/amd64   
- OS: ubuntu-20.04 LTS   
- text editor: VSCode    	
- terminal: ubuntu terminal 
```

## Directory Structure
```
.
├── README.md
├── functions
│   └── functiuons.go
├── handlers
│   └── handlers.go
├── go.mod
├── go.sum
├── coindatabase.db
├── main.go

2 directories, 6 files
```

## Usage
```bash
cd $GOPATH/src/github.com/<username>
git clone https://github.com/dinesh-cpu/iitk-coin.git
cd repo
go run main.go     
#, or build the program and run the executable
go build
./iitk-coin
```

Output should look like

```
2021/06/20 23:59:40 FINALDATA table created (if not existed) successfully!
2021/06/20 23:59:40 EVENTS table created (if not existed) successfully!
2021/06/20 23:59:40 REDEEM table created (if not existed) successfully!
2021/06/20 23:59:40 Serving at 8080
```

## Endpoints
POST requests take place via `JSON` requests. A typical usage would look like

```bash

```

- `/login` : `POST`
```json
{"rollno":"<rollno>", "password":"<password>"}
```

- `/signup` : `POST`
```json
{"name":"<username>","rollno":"<user rollno>", "password":"<password>","batch":"<user batch>"}
```

- `/logout` : `POST`
```json

```
- `/redeemcoins` : `POST`
```json
  {"coin":"<How much coin want to redeem>", "item":"item name"}
```

- `/acion` : `POST`
```json
{"id":"<Id of redeem request>", "action":"<0 or 1>"}
```

GET requests:

- `/pendingrequests` : `GET`
```bash
curl http://localhost:8080/pendingrequests
```

- `/getcoin` : `GET`
```bash
curl http://localhost:8080/getcoin
```

## Models

-  Credentials
```go
	Name     string `json:"username"`
	Password string `json:"password"`
	Rollno   int    `json:"rollno"`
	Batch    string `json:"batch"`
```

- Redeem
```go
	Coin int    `json:"coin"`
	Item string `json:"item"`
```

- Transfercoin
```go
	Rollno2 int `json:"rollno1"`
	Coin    int `json:"coin"`
```
