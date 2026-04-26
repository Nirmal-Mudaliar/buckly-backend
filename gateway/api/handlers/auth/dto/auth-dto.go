package dto

type User struct {
	Id              int64  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	PhoneNo         string `json:"phone_no"`
	DateOfBirth     string `json:"date_of_birth"`
	Gender          string `json:"gender"`
	Bio             string `json:"bio"`
	ProfilePhotoUrl string `json:"profile_photo_url"`
	HomeCountryId   int64  `json:"home_country_id"`
	HomeStateId     int64  `json:"home_state_id"`
	HomeCityId      int64  `json:"home_city_id"`
	IsPhoneVerified bool   `json:"is_phone_verified"`
	TrustScore      int64  `json:"trust_score"`
	Status          string `json:"status"`
	InsertTs        string `json:"insert_ts"`
	ModifiedTs      string `json:"modified_ts"`
}

type SignUpRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNo     string `json:"phone_no"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}

type SignUpResponse struct {
	Message string `json:"message"`
}
