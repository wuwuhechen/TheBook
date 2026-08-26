package structs

type User struct {
	UserID   uint   `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`
	Username string `gorm:"column:username;type:text;not null" json:"username"`
	Password string `gorm:"column:password;type:text;not null" json:"password"`
	Nickname string `gorm:"column:nickname;type:text;not null" json:"nickname"`
	Role     string `gorm:"column:role;type:text;not null" json:"role"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
