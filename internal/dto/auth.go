package dto

type RegisterUser struct {
	// The email of the user.
	// Required: true
	Email string `json:"email" binding:"required,email"`

	// The password of the user.
	// Required: true
	Password string `json:"password" binding:"required,min=8,max=64"`

	// The first name of the user.
	// Required: true
	FirstName string `json:"first_name" binding:"required"`

	// The last name of the user.
	// Required: true
	LastName string `json:"last_name" binding:"required"`

	// The phone number of the user.
	Phone string `json:"phone" binding:"omitempty,e164"`

	// The birth date of the user.
	BirthDate string `json:"birth_date" binding:"omitempty,datetime=2006-01-02"`
}
