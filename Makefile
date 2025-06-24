up:
	migrate -path migrations -database "postgres://aleksey:qwerty@localhost:5444/myDB?sslmode=disable" up

down:
	migrate -path migrations -database "postgres://aleksey:qwerty@localhost:5444/myDB?sslmode=disable" down