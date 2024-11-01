package db

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgxpool"
)

func Connect(route string)(*pgxpool.Pool, error){
 dbpool,err := pgxpool.New(context.Background(), route)
 if err != nil {
	return nil, fmt.Errorf("don't connect to data base: %w",err)
 }
//  defer dbpool.Close()



 fmt.Println("Conexion establecida")
 return dbpool,nil;
}