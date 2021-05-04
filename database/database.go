package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"

	entityUser "github.com/simple-app/entity/user"
)

// GetUserProfile is function for getting detail profile from given user ID
func GetUserProfile(ctx context.Context, userID int64, dbconn *sql.DB) (entityUser.UserProfile, error) {

	userProfile := entityUser.UserProfile{}

	//verify connection
	err := dbconn.Ping()
	if err != nil {
		return userProfile, err
	}

	query := `
	SELECT
			id,
			name,
			address,
			gender
		FROM
			user_profile
		WHERE
			id = $1
	`

	rows, err := dbconn.Query(query, userID)
	if err != nil {
		log.Println("failed when get data from table user_profile, cause :  ", err)
		return userProfile, err
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(
			&userProfile.ID,
			&userProfile.Name,
			&userProfile.Address,
			&userProfile.Gender,
		)
		if err != nil {
			return userProfile, err
		}
	}

	return userProfile, nil
}

// GetUserFamily is function for getting detail family from given user ID
func GetUserFamily(ctx context.Context, userID int64, dbconn *sql.DB) ([]entityUser.UserFamily, error) {

	//verify connection
	err := dbconn.Ping()
	if err != nil {
		return nil, err
	}

	query := `
	SELECT
			user_id,
			name,
			relation
		FROM
			user_family
		WHERE
			user_id = $1
	`

	rows, err := dbconn.QueryContext(ctx, query, userID)
	if err != nil {
		log.Println("failed when get data from table user_family, cause :  ", err)
		return nil, err
	}

	defer rows.Close()

	var userFamilies []entityUser.UserFamily

	for rows.Next() {
		userFamily := entityUser.UserFamily{}
		err = rows.Scan(
			&userFamily.UserID,
			&userFamily.Name,
			&userFamily.Relation,
		)
		if err != nil {
			return nil, err
		}
		userFamilies = append(userFamilies, userFamily)
	}

	return userFamilies, nil
}

// GetUserTransportation is function for getting detail transportation from given user ID
func GetUserTransportation(ctx context.Context, userID int64, dbconn *sql.DB) ([]entityUser.UserTransportation, error) {

	//verify connection
	err := dbconn.Ping()
	if err != nil {
		return nil, err
	}

	query := `
	SELECT
			user_id,
			name,
			type,
			colour
		FROM
			user_transportation
		WHERE
			user_id = $1
	`

	rows, err := dbconn.QueryContext(ctx, query, userID)
	if err != nil {
		log.Println("failed when get data from table user_transportation, cause :  ", err)
		return nil, err
	}

	defer rows.Close()

	var userTransportations []entityUser.UserTransportation
	for rows.Next() {
		userTransport := entityUser.UserTransportation{}
		err = rows.Scan(
			&userTransport.UserID,
			&userTransport.Name,
			&userTransport.TypeVehicle,
			&userTransport.Colour,
		)
		if err != nil {
			return nil, err
		}
		userTransportations = append(userTransportations, userTransport)
	}

	return userTransportations, nil
}

// GetCacheUserData is function for fet data from redis by given key
func GetCacheUserData(ctx context.Context, redis *redis.Client, key string, userID int64) (stringData string, err error) {

	// create key redis
	keyRedis := fmt.Sprintf(key, userID)

	return redis.Get(ctx, keyRedis).Result()
}

// SetCacheUserData is function for set data to redis
func SetCacheUserData(ctx context.Context, redis *redis.Client, key string, userID int64, data interface{}, expiredTime time.Duration) error {

	// create key redis
	keyRedis := fmt.Sprintf(key, userID)

	// marshaling data
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return redis.SetEX(ctx, keyRedis, json, expiredTime).Err()
}
