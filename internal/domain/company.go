package domain

import (
	"time"
)

// Company 会社エンティティ
// @Description 会社情報
type Company struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement" example:"1"`                     // 会社ID
	Name        string    `json:"name" gorm:"size:100;not null" validate:"required,min=2,max=100" example:"株式会社サンプル"`  // 会社名
	Email       string    `json:"email" gorm:"size:255;uniqueIndex;not null" validate:"required,email" example:"info@sample.co.jp"` // 会社メールアドレス
	Phone       string    `json:"phone" gorm:"size:20" validate:"omitempty,min=10,max=20" example:"03-1234-5678"`            // 電話番号
	Address     string    `json:"address" gorm:"size:500" validate:"omitempty,max=500" example:"東京都渋谷区..."`               // 住所
	Website     string    `json:"website" gorm:"size:255" validate:"omitempty,url" example:"https://sample.co.jp"`           // ウェブサイト
	Description string    `json:"description" gorm:"type:text" validate:"omitempty,max=1000" example:"IT関連のサービスを提供しています"`  // 会社説明
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`                                    // 作成日時
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`                                    // 更新日時
}

// TableName テーブル名を指定
func (Company) TableName() string {
	return "companies"
}

// IsValidCompany 会社データが有効か検証
func (c *Company) IsValidCompany() bool {
	return c.Name != "" && c.Email != ""
}

// GetDisplayName 表示用の名前を取得
func (c *Company) GetDisplayName() string {
	if c.Name != "" {
		return c.Name
	}
	return c.Email
}

// HasContact 連絡先情報があるか確認
func (c *Company) HasContact() bool {
	return c.Phone != "" || c.Email != "" || c.Address != ""
}

// HasWebsite ウェブサイト情報があるか確認
func (c *Company) HasWebsite() bool {
	return c.Website != ""
}

// ToResponseCompany レスポンス用のCompanyを返す
func (c *Company) ToResponseCompany() Company {
	return Company{
		ID:          c.ID,
		Name:        c.Name,
		Email:       c.Email,
		Phone:       c.Phone,
		Address:     c.Address,
		Website:     c.Website,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// CompanyUser User-Company関係エンティティ（多対多関係用）
// @Description ユーザーと会社の関係
type CompanyUser struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null;index" example:"1"`        // ユーザーID
	CompanyID uint      `json:"company_id" gorm:"not null;index" example:"1"`     // 会社ID
	Role      string    `json:"role" gorm:"size:50;default:'member'" example:"admin"` // 役割（admin, member, etc.）
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`                 // 関係作成日時
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`                 // 関係更新日時
}

// TableName テーブル名を指定
func (CompanyUser) TableName() string {
	return "company_users"
}

// IsAdmin 管理者権限があるか確認
func (cu *CompanyUser) IsAdmin() bool {
	return cu.Role == "admin"
}

// IsMember メンバー権限があるか確認
func (cu *CompanyUser) IsMember() bool {
	return cu.Role == "member" || cu.Role == "admin"
}
