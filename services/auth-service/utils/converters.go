package utils

import (
	auth_gen "buckly-ms/proto/auth-gen"
	database_gen "buckly-ms/proto/database-gen"
)

// DatabaseUserToAuthUser converts a database_gen.User to auth_gen.User
func DatabaseUserToAuthUser(dbUser *database_gen.User) *auth_gen.User {
	if dbUser == nil {
		return nil
	}

	return &auth_gen.User{
		Id:              dbUser.Id,
		FirstName:       dbUser.FirstName,
		LastName:        dbUser.LastName,
		Email:           dbUser.Email,
		PhoneNo:         dbUser.PhoneNo,
		DateOfBirth:     dbUser.DateOfBirth,
		Gender:          dbUser.Gender,
		Bio:             dbUser.Bio,
		ProfilePhotoUrl: dbUser.ProfilePhotoUrl,
		HomeCountryId:   dbUser.HomeCountryId,
		HomeStateId:     dbUser.HomeStateId,
		HomeCityId:      dbUser.HomeCityId,
		IsPhoneVerified: dbUser.IsPhoneVerified,
		TrustScore:      dbUser.TrustScore,
		Status:          dbUser.Status,
		InsertTs:        dbUser.InsertTs,
		ModifiedTs:      dbUser.ModifiedTs,
	}
}
