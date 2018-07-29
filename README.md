# subdomain
easyway to create a proxy redirect to from subdomain with apache


## Installing Go
```
apt update
apt install golang-go -y
```

## Making working folder
```
mkdir ~/go
export GOPATH=~/go
```

## Install Cobra
``` 
go get -u github.com/spf13/cobra/cobra
```

## Getting and building this project
```
go get -u -v -d github.com/martijn1279/subdomain
cd $GOPATH/src/github.com/martijn1279/subdomain
go build
```

## Running program
```
./subdomain add test.com subdomain http://localhost.com/
```
