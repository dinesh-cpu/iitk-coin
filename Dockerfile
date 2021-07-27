FROM golang:1.16

LABEL maintainer = "dineshcpu <dineshmatrix2@gmail.com>"
#Set the current working directory inside the container

WORKDIR $GOPATH/src/github.com/dinesh-cpu/iitk-coin

#commands
copy go.mod .
copy go.sum .
RUN go mod download

#copy everything 
COPY . .

#BUILD Executable 
RUN go build

EXPOSE 8080

#run executable
CMD ["./dinesh-cpu"]
