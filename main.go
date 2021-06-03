package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/simple-app/database"
	entityPostgresql "github.com/simple-app/entity/postgresql"
	entityRedis "github.com/simple-app/entity/redis"
	entityUser "github.com/simple-app/entity/user"
)

type handlerUser struct {
	db    *sql.DB
	redis *redis.Client
}

// TrackerDuration execution
func TrackerDuration(start time.Time, funcName string) {
	duration := time.Since(start)
	log.Printf("%s took %s", funcName, duration)
}

// GetUserDetail from given param user id
func (h *handlerUser) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	//for tracking time execution
	startTime := time.Now()

	// get param request
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	// This is the option when you want to compare the performance
	// set true to using cache (using redis first)
	// set false to not using cache (purely psotgresql)
	useCache := true

	// We set context with timeout, so if the duration for calling this function
	// is more then how it set, it would be timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var (
		errs               = []error{}
		wg                 = sync.WaitGroup{}
		mux                = sync.Mutex{}
		userProfile        entityUser.UserProfile
		userFamily         []entityUser.UserFamily
		userTransportation []entityUser.UserTransportation
	)

	// parallel process for user profile
	wg.Add(1)
	go func() {
		defer wg.Done()
		userProfile, err = h.GetUserProfile(ctx, id, useCache)
		if err != nil {
			// locking process while append error to array of errors
			mux.Lock()
			errs = append(errs, err)
			mux.Unlock()
			return
		}
		TrackerDuration(startTime, "GetUserProfile")

	}()

	// parallel process for user family
	wg.Add(1)
	go func() {
		defer wg.Done()
		userFamily, err = h.GetUserFamily(ctx, id, useCache)
		if err != nil {
			// locking process while append error to array of errors
			mux.Lock()
			errs = append(errs, err)
			mux.Unlock()
			return
		}
		TrackerDuration(startTime, "GetUserFamily")

	}()

	// parallel process for user transportation
	wg.Add(1)
	go func() {
		defer wg.Done()
		userTransportation, err = h.GetUserTransportation(ctx, id, useCache)
		if err != nil {
			// locking process while append error to array of errors
			mux.Lock()
			errs = append(errs, err)
			mux.Unlock()
			return
		}
		TrackerDuration(startTime, "GetUserTransportation")

	}()

	// waiting all data complete
	// and check if there is any error
	wg.Wait()
	if len(errs) > 0 {
		log.Println(errs)
		json.NewEncoder(w).Encode(errs)
		return
	}

	//aggregate data user detail
	userDetail := entityUser.UserDetail{}
	userDetail.Profile = userProfile
	userDetail.Family = userFamily
	userDetail.Transportation = userTransportation
	TrackerDuration(startTime, "All function")
	log.Println("==========================================")
	json.NewEncoder(w).Encode(userDetail)
	return
}

// GetUserProfile detail data user profile
func (h *handlerUser) GetUserProfile(ctx context.Context, userID int64, useCache bool) (entityUser.UserProfile, error) {

	var (
		userProfile entityUser.UserProfile
		err         error
	)

	if useCache {
		// get data from cache redis
		strData, err := database.GetCacheUserData(ctx, h.redis, entityRedis.UserProfile, userID)
		if err != nil && err != redis.Nil {
			log.Println("error getting data user profile from redis :  ", err)
		}

		if strData != "" {
			err = json.Unmarshal([]byte(strData), &userProfile)
			if err == nil {
				return userProfile, err
			}
			log.Println("error getting data user profile when unmarshal struct :  ", err)
		}
	}

	// get data from DB
	userProfile, err = database.GetUserProfile(ctx, userID, h.db)
	if err != nil {
		return userProfile, err
	}

	if useCache {
		// set data on redis
		err = database.SetCacheUserData(ctx, h.redis, entityRedis.UserProfile, userID, userProfile, entityRedis.SetTimeExp)
		if err != nil {
			log.Println("error when set data user profile on redis :  ", err)
		}

	}

	return userProfile, nil
}

// GetUserFamily detail data user family
func (h *handlerUser) GetUserFamily(ctx context.Context, userID int64, useCache bool) ([]entityUser.UserFamily, error) {

	var (
		userFamilies []entityUser.UserFamily
		err          error
	)

	if useCache {
		// get data from cache redis
		strData, err := database.GetCacheUserData(ctx, h.redis, entityRedis.UserFamily, userID)
		if err != nil && err != redis.Nil {
			log.Println("error getting data user family from redis :  ", err)
		}

		if strData != "" {
			err = json.Unmarshal([]byte(strData), &userFamilies)
			if err == nil {
				return userFamilies, err
			}
			log.Println("error getting data user family when unmarshal struct :  ", err)
		}
	}

	// get data from DB
	userFamilies, err = database.GetUserFamily(ctx, userID, h.db)
	if err != nil {
		return userFamilies, err
	}

	if useCache {
		// set data on redis
		err = database.SetCacheUserData(ctx, h.redis, entityRedis.UserFamily, userID, userFamilies, entityRedis.SetTimeExp)
		if err != nil {
			log.Println("error when set data user family on redis :  ", err)
		}
	}

	return userFamilies, nil
}

// GetUserTransportation detail data user transportation
func (h *handlerUser) GetUserTransportation(ctx context.Context, userID int64, useCache bool) ([]entityUser.UserTransportation, error) {

	var (
		userTransportations []entityUser.UserTransportation
		err                 error
	)

	if useCache {
		// get data from cache redis
		strData, err := database.GetCacheUserData(ctx, h.redis, entityRedis.UserTransportation, userID)
		if err != nil && err != redis.Nil {
			log.Println("error getting data user transportation from redis :  ", err)
		}

		if strData != "" {
			err = json.Unmarshal([]byte(strData), &userTransportations)
			if err == nil {
				return userTransportations, err
			}
			log.Println("error getting data user transportation when unmarshal struct :  ", err)
		}
	}

	// get data from DB
	userTransportations, err = database.GetUserTransportation(ctx, userID, h.db)
	if err != nil {
		return userTransportations, err
	}

	if useCache {
		// set data on redis
		err = database.SetCacheUserData(ctx, h.redis, entityRedis.UserTransportation, userID, userTransportations, entityRedis.SetTimeExp)
		if err != nil {
			log.Println("error when set data user transportation on redis :  ", err)
		}
	}

	return userTransportations, nil
}

func handleRequests(h *handlerUser) {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/user/{id}", h.GetUserDetail)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {

	// create string conection for psql
	psql := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		entityPostgresql.Host, entityPostgresql.Port, entityPostgresql.User, entityPostgresql.Password, entityPostgresql.DBname)

	// connect to DB
	dbconn, err := sql.Open("postgres", psql)
	if err != nil {
		log.Println(err)
		return
	}
	defer dbconn.Close()

	// string address for redis
	redisAddr := fmt.Sprintf("%s:%d",
		entityRedis.Host, entityRedis.Port)

	// connect to Redis
	rediClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	handler := &handlerUser{
		db:    dbconn,
		redis: rediClient,
	}

	fmt.Println("===== || Starting Apps ||=====")
	handleRequests(handler)

}
