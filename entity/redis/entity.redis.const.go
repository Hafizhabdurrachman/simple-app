package redis

import "time"

// connection redis const
const (
	Host = "127.0.0.1"
	Port = 6379
)

// config redis
const (
	// set time key expiration on second
	SetTimeExp = 1 * time.Minute
)

// possible key caching data
const (
	// key for detail data user profile
	UserProfile = "user_profile:id:%d"

	// key for detail data user family
	UserFamily = "user_family:id:%d"

	// key for detail data user transportation
	UserTransportation = "user_transportation:id:%d"
)
