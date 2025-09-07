package helper

// PaginationRequest ページネーション用リクエスト
type PaginationRequest struct {
	Page  int `query:"page" validate:"omitempty,min=1" example:"1"`        // ページ番号（1から開始）
	Limit int `query:"limit" validate:"omitempty,min=1,max=100" example:"10"` // 1ページあたりの件数
}

// GetOffset オフセット値を計算
func (p *PaginationRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetLimit()
}

// GetLimit リミット値を取得（デフォルト値設定）
func (p *PaginationRequest) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	return p.Limit
}

// PaginationResponse ページネーション用レスポンス
type PaginationResponse struct {
	Page       int   `json:"page" example:"1"`        // 現在のページ番号
	Limit      int   `json:"limit" example:"10"`      // 1ページあたりの件数
	Total      int64 `json:"total" example:"100"`     // 総件数
	TotalPages int   `json:"total_pages" example:"10"` // 総ページ数
}

// NewPaginationResponse ページネーションレスポンスを作成
func NewPaginationResponse(page, limit int, total int64) *PaginationResponse {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 0 {
		totalPages = 0
	}
	
	return &PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// SearchRequest 検索用リクエスト
type SearchRequest struct {
	Query string `query:"q" validate:"omitempty,max=100" example:"search term"` // 検索クエリ
	*PaginationRequest
}

// SortRequest ソート用リクエスト
type SortRequest struct {
	SortBy  string `query:"sort_by" validate:"omitempty,oneof=id name email created_at updated_at" example:"created_at"` // ソート項目
	SortDir string `query:"sort_dir" validate:"omitempty,oneof=asc desc" example:"desc"`                               // ソート方向
}

// GetSortBy ソート項目を取得（デフォルト値設定）
func (s *SortRequest) GetSortBy() string {
	if s.SortBy == "" {
		return "created_at"
	}
	return s.SortBy
}

// GetSortDir ソート方向を取得（デフォルト値設定）
func (s *SortRequest) GetSortDir() string {
	if s.SortDir == "" {
		return "desc"
	}
	return s.SortDir
}

// IDRequest ID用リクエスト
type IDRequest struct {
	ID uint `param:"id" validate:"required,min=1" example:"1"` // エンティティID
}

// ErrorCode エラーコード
type ErrorCode string

const (
	ErrorCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrorCodeAlreadyExists   ErrorCode = "ALREADY_EXISTS"
	ErrorCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrorCodeInternalError   ErrorCode = "INTERNAL_ERROR"
	ErrorCodeDatabaseError   ErrorCode = "DATABASE_ERROR"
	ErrorCodeExternalAPI     ErrorCode = "EXTERNAL_API_ERROR"
)

// String ErrorCodeの文字列表現
func (e ErrorCode) String() string {
	return string(e)
}