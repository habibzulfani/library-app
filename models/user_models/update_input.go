package user_models

type UserUpdateInput struct {
	Email     string `json:"email" form:"email" validate:"omitempty,email"`
	Name      string `json:"name" form:"name" validate:"omitempty,min=2,max=100"`
	Password  string `json:"password" form:"password" validate:"omitempty,min=8"`
	UserType  string `json:"user_type" form:"user_type" validate:"omitempty,oneof=student teacher"`
	IDNumber  string `json:"id_number" form:"id_number" validate:"omitempty"`
	JurusanID uint   `json:"jurusan_id" form:"jurusan_id" validate:"omitempty,numeric"`
	Address   string `json:"address" form:"address" validate:"omitempty"`
	Role      string `json:"role" form:"role" validate:"omitempty,oneof=user admin"`
}
