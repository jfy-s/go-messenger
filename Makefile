migrations up:
	goose -dir migrations postgres "postgres://postgres:test@localhost:5432/messenger" up
migrations down:
	goose -dir migrations postgres "postgres://postgres:test@localhost:5432/messenger" down