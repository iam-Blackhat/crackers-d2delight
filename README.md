# crackers-d2delight

# Migration table

# Create table

## migrate create -ext sql -dir db/migrations create_users_table

## migrate -path db/migrations -database "postgres://datasirpi:YourPass@localhost:5432/crackers?sslmode=disable" up
