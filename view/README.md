# Set Up Tutorial
1. Install mongodb for save data, redis for cache and fill info to file config.yaml
2. Start : go run main.go
# Flow
1. Question will be load from file to mongodb and client will get from mongo db, log will be write in both console and file server.log
2. When user login, token will be create and save on redis db, also a correspoding cookie will be create at client, if you want to log out you have to delete cookie
3. Each time user save an answer, answer of that question will be send to server and save, this principle and cookie will allow for saving state of survey
4. When user finish survey, list true answer of question wil response to client and client will use that to evaluate survey ,user will be label have done and can not access to survey any more
5. Front-end is not done but that will be the flow if front-end will be done