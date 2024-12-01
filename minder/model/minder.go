package model

type RegisterReq struct {
	Username       string `json:"username"  schema:"username" validate:"required"`
	Email          string `json:"email"  schema:"email" validate:"required"`
	PhoneNumber    string `json:"phoneNumber" schema:"phoneNumber" validate:"required"`
	Password       string `json:"password"  schema:"password" validate:"required"`
	FullName       string `json:"fullName"  schema:"fullName" validate:"required"`
	Gender         string `json:"gender"  schema:"gender" validate:"required,oneof=male female"`
	DateOfBirth    string `json:"dateOfBirth"  schema:"dateOfBirth" validate:"required,omitempty,dateformat"`
	ProfilePicture string `json:"profilePicture"  schema:"profilePicture" validate:"required"`
}

type LoginReq struct {
	Username string `json:"username" schema:"username" validate:"required"`
	Password string `json:"password" schema:"password" validate:"required"`
}

type SwipeReq struct {
	Id       int64  `json:"id" schema:"id" validate:"required"`
	TargetId int64  `json:"targetId" schema:"targetId" validate:"required"`
	Action   string `json:"action" schema:"action" validate:"required,oneof=like pass"`
}

type UserData struct {
	UserId         int64  `json:"userId"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	FullName       string `json:"fullName"`
	Gender         string `json:"gender"`
	DateOfBirth    string `json:"dateOfBirth"`
	IsUpgraded     bool   `json:"isUpgraded"`
	ProfilePicture string `json:"profilePicture"`
}

type TargetUserData struct {
	UserId         int64  `json:"userId"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phoneNumber"`
	FullName       string `json:"fullName"`
	Gender         string `json:"gender"`
	DateOfBirth    string `json:"dateOfBirth"`
	ProfilePicture string `json:"profilePicture"`
}
