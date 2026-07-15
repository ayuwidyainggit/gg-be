package constant

const (
	PasswordResetStatusPendingOTP    = "pending_otp"
	PasswordResetStatusOTPValidated  = "otp_validated"
	PasswordResetStatusCompleted     = "completed"
	PasswordResetStatusLocked         = "locked"
	PasswordResetStatusExpired        = "expired"
)

const (
	ErrCodeEmailNotRegistered           = "EMAIL_NOT_REGISTERED"
	ErrCodeOTPInvalidOrExpired          = "OTP_INVALID_OR_EXPIRED"
	ErrCodeOTPMaxAttempts               = "OTP_MAX_ATTEMPTS"
	ErrCodeOTPCooldown                  = "OTP_COOLDOWN"
	ErrCodeOTPMaxResend                 = "OTP_MAX_RESEND"
	ErrCodeResetTokenInvalidOrExpired   = "RESET_TOKEN_INVALID_OR_EXPIRED"
	ErrCodePasswordPolicy               = "PASSWORD_POLICY_VIOLATION"
	ErrCodePasswordConfirmMismatch      = "PASSWORD_CONFIRM_MISMATCH"
	ErrCodeDeprecatedPasswordPatch      = "DEPRECATED_USE_RESET_PASSWORD"
)
