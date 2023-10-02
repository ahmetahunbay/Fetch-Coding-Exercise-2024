This is my Fetch 2024 internship take home project. It is written entirely in Go, but uses SQLite for a persistent data store

To run this code, you will need to download SQLite [here](https://www.sqlite.org/download.html), and [Go](https://go.dev/doc/install) here.

The database uses the name "FetchExerciseData.db", so if you hava another db of the same name, I would ask you to please change it. 

There are 4 main API routes you can run. The 3 from the spec, and an additional one: /clear. This additional endpoint is a DELETE method that clears the database and cache, if you want to start fresh.

To build the executable, run the following command:
go build main.go

You can then run the code by running the main executable with the following command:
./main

<h2>Calling the Endpoints</h2>

You can call the endpoints from [here](https://www.postman.com/ahmetahunbay/workspace/fetch-coding-exercise-2024/request/30177900-f8c5bdb0-8f0f-47e6-94af-cf8dc765d6ab) -- my public postman workspace

Otherwise, if you have [curl](https://www.bing.com/search?q=curl+download&cvid=2e9113c250074d28b74890c8399c27e6&gs_lcrp=EgZjaHJvbWUyBggAEEUYOTIECAEQADIECAIQADIECAMQADIECAQQADIECAUQADIECAYQADIECAcQADIECAgQANIBCDM5MDFqMGo5qAIAsAIA&FORM=ANAB01&PC=LCTS) installed, you can call the api endpoints using the following commands:

for /add route, run:
curl -X POST -H "Content-Type: application/json" -d '{ "payer": "DANNON", "points": 400, "timestamp": "2022-11-02T14:00:00Z" }' http://localhost:8000/add

for /spend route, run:
curl -X POST -H "Content-Type: application/json" -d '{"points": 50}' http://localhost:8000/spend

for /balance route, run:
curl http://localhost:8000/balance

for /clear route, run;
curl http://localhost:8000/clear
