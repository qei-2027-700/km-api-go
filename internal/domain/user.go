package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User ユーザーエンティティ
// @Description ユーザー情報
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`                    // ユーザーID
	Name      string    `json:"name" gorm:"size:50;not null" validate:"required,min=2,max=50" example:"山田太郎"`          // ユーザー名
	Email     string    `json:"email" gorm:"size:255;uniqueIndex;not null" validate:"required,email" example:"yamada@example.com"` // メールアドレス
	Password  string    `json:"-" gorm:"size:255;not null" validate:"required,min=8"`              // パスワード（レスポンスに含めない）
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`                                   // 作成日時
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`                                   // 更新日時
}

// TableName テーブル名を指定
func (User) TableName() string {
	return "users"
}

// HashPassword パスワードをハッシュ化
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword パスワードを検証
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsValidUser ユーザーデータが有効か検証
func (u *User) IsValidUser() bool {
	return u.Name != "" && u.Email != "" && u.Password != ""
}

// BeforeCreate GORM フック - 作成前にパスワードをハッシュ化
func (u *User) BeforeCreate() error {
	if u.Password != "" {
		return u.HashPassword()
	}
	return nil
}

// GetDisplayName 表示用の名前を取得
func (u *User) GetDisplayName() string {
	if u.Name != "" {
		return u.Name
	}
	return u.Email
}

// ToResponseUser レスポンス用にパスワードを除いたUserを返す
func (u *User) ToResponseUser() User {
	return User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
