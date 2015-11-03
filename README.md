### Install Deps

Install Mongo (2.6.5)

Install Go (1.4.2) and prepare default GOPATH/GOROOT env variable

Install Node (0.12.7)

##### Next build and run (all paths from root of project sources)

  cd server && go get github.com/go-martini/martini github.com/martini-contrib/render gopkg.in/mgo.v2 gopkg.in/mgo.v2/bson
  cd client/build && npm install && bower install && sudo npm install -g gulp bower

##### Build and run

  mkdir ~/mongo && mongod --dbpath ~/mongo
  go build server/main.go && MARTINI_ENV=production server/main
  cd client/build && gulp build

##### Run Dev Server

  mkdir ~/mongo && mongod --dbpath ~/mongo
  go run server/main.go
  cd client/build && gulp run
