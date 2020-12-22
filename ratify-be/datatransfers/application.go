package datatransfers

type ApplicationInfo struct {
	ClientID     string `json:"client_id" binding:"-"`
	ClientSecret string `json:"client_secret" binding:"-"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"-"`
	LoginURL     string `json:"login_url" binding:"required"`
	CallbackURL  string `json:"callback_url" binding:"required"`
	LogoutURL    string `json:"logout_url" binding:"required"`
	Metadata     string `json:"metadata" binding:"-"`
	CreatedAt    int64  `json:"created_at" binding:"-"`
	UpdatedAt    int64  `json:"updated_at" binding:"-"`
}

type ApplicationClientIDURI struct {
	ClientID string `uri:"client_id"`
}
