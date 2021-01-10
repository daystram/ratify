package datatransfers

type UserLogin struct {
	Username string `json:"preferred_username" binding:"-"`
	Password string `json:"password" binding:"-"`
}

type UserSignup struct {
	GivenName  string `json:"given_name" binding:"required"`
	FamilyName string `json:"family_name" binding:"required"`
	Username   string `json:"preferred_username" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type UserUpdate struct {
	GivenName  string `json:"given_name" binding:"required"`
	FamilyName string `json:"family_name" binding:"required"`
	Email      string `json:"email" binding:"required"`
}

type UserVerify struct {
	Token string `json:"token" binding:"required"`
}

type UserInfo struct {
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Subject       string `json:"sub"`
	Username      string `uri:"preferred_username" json:"preferred_username"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	CreatedAt     int64  `json:"created_at"`
}
