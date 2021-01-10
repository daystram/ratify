package datatransfers

type ApplicationInfo struct {
	ClientID    string `json:"client_id,omitempty" binding:"-"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty" binding:"-"`
	LoginURL    string `json:"login_url,omitempty" binding:"required"`
	CallbackURL string `json:"callback_url,omitempty" binding:"required"`
	LogoutURL   string `json:"logout_url,omitempty" binding:"required"`
	Metadata    string `json:"metadata" binding:"-"`
	Locked      *bool  `json:"locked,omitempty" binding:"-"`
	CreatedAt   int64  `json:"created_at,omitempty" binding:"-"`
	UpdatedAt   int64  `json:"updated_at,omitempty" binding:"-"`
}
