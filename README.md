# GoProject

A Final year project on blood donors registration website named ***DonorSpace***.

### Setup guide
* Read MONGO.md file to setup mongoDB tables first.
* Setup Gmail API with proper authorization as given in link in Dependencies section.
* Install following dependencies:
    * go get -u google.golang.org/api/gmail/v1
    * go get -u golang.org/x/oauth2/...
    * go get -u github.com/gorilla/mux
    * go get -u github.com/kardianos/govendor
    * go get gopkg.in/yaml.v2


### Dependencies
* **mux**: implements a request router and dispatcher for matching incoming requests to their respective handler. Link: [gorilla/mux](https://github.com/gorilla/mux)
* **govendor**: downloads go dependencies in local. Link: [govendor](https://github.com/kardianos/govendor)
* __yaml__: configuration key values are stored as yaml files. Link: [yaml.v2](https://gopkg.in/yaml.v2)
* __Gmail API__: sends emails to users registering to be donors. Link: [Google API doc](https://developers.google.com/gmail/api/quickstart/go)
