### Install Deps

Install Mongo (2.6.5)

Install Go (1.4.2) and prepare default GOPATH/GOROOT env variable

Install Node (0.12.7)

##### Next build and run (all paths from root of project sources)

    cp -r server/vendor ~/go/src
    cd client/builder && npm install && bower install && sudo npm install -g gulp bower

##### Build and run

    mkdir ~/mongo && mongod --dbpath ~/mongo
    go build server/main.go && MARTINI_ENV=production server/main
    cd client/builder && gulp build

##### Run Dev Server

    mkdir ~/mongo && mongod --dbpath ~/mongo
    go run server/main.go
    cd client/builder && gulp run
