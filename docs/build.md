## Application building

Install dependencies
```bash
go get -d ./src/
```

Build app
```bash
go build -o ac-export ./src/
```

Build app for Linux
```bash
env GOOS=linux GOARCH=amd64 GOARM=7 go build -o ac-export-linux ./src/
```