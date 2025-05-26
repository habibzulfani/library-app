package user_models

type UserRegisterInput struct {
    Email     string `json:"email" form:"email" validate:"required,email"`
    Name      string `json:"name" form:"name" validate:"required,min=2,max=100"`
    Password  string `json:"password" form:"password" validate:"required,min=8"`
    UserType  string `json:"user_type" form:"user_type" validate:"required,oneof=student teacher"`
    IDNumber  string `json:"id_number" form:"id_number" validate:"required"`
    JurusanID uint   `json:"jurusan_id" form:"jurusan_id" validate:"required,numeric"`
    Jurusan   string `json:"jurusan" form:"jurusan" validate:"required"`
    Address   string `json:"address" form:"address" validate:"required"`
    Role      string `json:"role" form:"role" validate:"required,oneof=user admin"`
}