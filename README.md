# TODO

* After every change to th prisma data model, we need to run `prisma deploy`.
* Need to add `managementApiSecret` to the `docker-compose.yml`.

# Usage

## Development Mode
To run the project in *development* mode, we need to run the `docker-compose.dev.yml` file, like so:
```
sudo docker-compose -f docker-compose.dev.yml up -d
```
And then, run the server's executable:
```
go run server/cmd/main.go
```

If `redis` is installed on the local machine, you may need to disable it:
```
sudo systemctl stop redis
```
You can start it back up again, when you're done, with the command:
```
sudo systemctl start redis
```

## Release Mode
To run the project in *Release Mode*, we only need to run the `docker-compose.yml` file, like so:
```
sudo docker-compose up -d --build
```

(You may also have to disable `redis`, if you have it installed.)

# Tests
To run all the unit tests included, use the command:
```
go test ./...
```
