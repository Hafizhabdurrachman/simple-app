package user

// UserProfile describe detail his/her profile
type UserProfile struct {
	ID      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
	Gender  string `db:"gender" json:"gender"`
}

// UserFamily describe detail users family
type UserFamily struct {
	UserID   int64  `db:"user_id" json:"user_id"`
	Name     string `db:"name" json:"name"`
	Relation string `db:"relation" json:"relation"`
}

// UserTransportation describe detail users vehicle
type UserTransportation struct {
	UserID      int64  `db:"user_id" json:"user_id"`
	Name        string `db:"name" json:"name"`
	TypeVehicle string `db:"type" json:"type"`
	Colour      string `db:"colour" json:"colour"`
}

// UserDetail describe detail users vehicle
type UserDetail struct {
	Profile        UserProfile
	Family         []UserFamily
	Transportation []UserTransportation
}
