### Install Deps

Install Mongo (2.6.5)

Install Go (1.4.2) and prepare default GOPATH/GOROOT env variable

Install Node (0.12.7)

  cd server && go get github.com/go-martini/martini github.com/martini-contrib/render gopkg.in/mgo.v2 gopkg.in/mgo.v2/bson
  cd client/build && npm install && bower install && npm install -g gulp bower

### Run

  gulp build
  mkdir ~/mongo && mongod --dbpath ~/mongo
  go run server/main.go

### Run Dev Server

  gulp run
  mkdir ~/mongo && mongod --dbpath ~/mongo
  go run server/main.go
