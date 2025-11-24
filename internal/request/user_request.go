package request

// CreateUserRequest adalah struct untuk binding request JSON pada create user
// Digunakan untuk validasi input sebelum map ke model.User
type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`      // Wajib diisi
	Email    string `json:"email" binding:"required,email"`    // Wajib diisi dan harus format email
	Password string `json:"password" binding:"required,min=6"` // Wajib diisi, minimal 6 karakter
}
