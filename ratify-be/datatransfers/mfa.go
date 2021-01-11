package datatransfers

type TOTPRequest struct {
	OTP string `json:"otp" binding:"required"`
}
