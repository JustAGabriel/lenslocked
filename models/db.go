package models

type DbConnectionConfig struct {
}

func GetConnectionString() string {
	return "host=localhost user=baloo password=junglebook dbname=lenslocked port=5432 sslmode=disable"
}
