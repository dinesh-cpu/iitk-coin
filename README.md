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
- go version: go1.16.4 linux/amd64    # https://golang.org/dl/
- OS: ubuntu-20.04 LTS   # https://docs.microsoft.com/en-us/windows/wsl/install-win10
- text editor: VSCode    	# https://code.visualstudio.com/download
- terminal: ubuntu terminal   		# https://ohmyz.sh/
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
2021/06/20 23:59:40 User Database opened and table created (if not existed) successfully!
2021/06/20 23:59:40 Transaction Database opened and table created (if not existed) successfully!
2021/06/20 23:59:40 Wallet Database opened and table created (if not existed) successfully!
2021/06/20 23:59:40 Serving at 8080
```

## Endpoints
POST requests take place via `JSON` requests. A typical usage would look like

```bash

```

- `/login` : `POST`
```json
{"name":"<name>", "rollno":"<rollno>", "password":"<password>"}
```

- `/signup` : `POST`
```json
{"rollno":"<rollno>", "password":"<password>"}
```

- `/reward` : `POST`
```json
{"rollno":"<rollno>", "coins":"<coins>"}
```

- `/transfer` : `POST`
```json
{"sender":"<senderRollNo>", "receiver":"<receiverRollNo>", "coins":"<coins>"}
```

GET requests:

- `/secretpage` : `GET`
```bash
curl http://localhost:8080/secretpage
```

- `/balance` : `GET`
```bash
curl http://localhost:8080/balance?rollno=<rollno>
```

## Models

-  User
```go
	Name     string `json:"name"`
	Rollno   string  `json:"rollno"`
	Password string `json:"password"`
```

- RewardPayload
```go
	Rollno string `json:"rollno"`
	Coins  int64 `json:"coins,string"`
```

- TransferPayload
```go
	SenderRollno   string `json:"sender"`
	ReceiverRollno string `json:"receiver"`
	Coins          int64 `json:"coins,string"`
```
