GOPATH = C:\Users\Pawel\repos;C:\Users\Pawel\repos\ppppp

all:
	go build -o weather.exe weather
	go build -o webserver.exe webserver
	go build -o worker.exe worker
	docker-compose up --build

clean:
	rm -f weather.exe webserver.exe worker.exe
