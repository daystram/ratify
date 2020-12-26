package datatransfers

type UserLogin struct {
	Username string `json:"username" binding:"-"`
	Password string `json:"password" binding:"-"`
}

type UserSignup struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdate struct {
	Email string `json:"email" binding:"-"`
}

type UserInfo struct {
	Subject   string `json:"sub"`
	Username  string `uri:"username" json:"username"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}
