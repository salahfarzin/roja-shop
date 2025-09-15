## Generate database migration
```
migrate create -ext sql -dir migrations -seq create_products_table
```
## Migrate database
```
 ~/go/bin/migrate -database "sqlite://$PWD/storage/database/roja.db" -path "$PWD/migrations" up 
 ```