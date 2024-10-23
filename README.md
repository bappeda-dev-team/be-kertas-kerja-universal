# API E-KAK 


### Kebutuhan

- Go versi 1.22
- MySQL versi 8.0.33
- [golang-migrate](https://github.com/golang-migrate/migrate)


### Cara install

- buat database bernama db_ekak

- install golang-migrate di cmd

```sh
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

- add create migration
``` sh
migrate create -ext sql -dir db create_new_table
```

- migrasi database

```sh
migrate -path db -database "mysql://root@tcp(localhost:3306)/db_ekak" up
```


### Run server

ketikkan perintah:

```sh
go run main.go
```

untuk menghentikan server, tekan Ctrl + c
