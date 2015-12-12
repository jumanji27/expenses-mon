### Install Deps

Install Mongo (2.6.5)

Install Go (1.4.2) and prepare default GOPATH/GOROOT env variable

Install Node (0.12.7)

##### Next build and run (all paths from root of project sources)

    sudo cp -r server/vend ~/go/src
    cd client/builder && npm install && sudo npm install -g gulp@3.9.0 bower@1.6.5 && bower install

##### Build and run

    mkdir ~/mongo && mongod --dbpath ~/mongo
    go build server/main.go && MARTINI_ENV=production server/main
    cd client/builder && gulp build

##### Run Dev Server

    mkdir ~/mongo && mongod --dbpath ~/mongo
    go run server/main.go
    cd client/builder && gulp run