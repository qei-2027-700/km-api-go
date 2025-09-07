package company

import (
	"fmt"

	"km-api-go/internal/domain"
	"km-api-go/internal/helper"
	"km-api-go/internal/company/repository"
)

// CompanyUsecase 会社のビジネスロジックを処理
type CompanyUsecase struct {
	companyRepo     repository.CompanyRepository
	companyUserRepo repository.CompanyUserRepository
}

// NewCompanyUsecase 会社ユースケースのコンストラクタ
func NewCompanyUsecase(companyRepo repository.CompanyRepository, companyUserRepo repository.CompanyUserRepository) *CompanyUsecase {
	return &CompanyUsecase{
		companyRepo:     companyRepo,
		companyUserRepo: companyUserRepo,
	}
}

// GetAllCompanies 全会社を取得
func (uc *CompanyUsecase) GetAllCompanies() ([]domain.Company, error) {
	companies, err := uc.companyRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all companies: %w", err)
	}

	// レスポンス用データに変換
	var responseCompanies []domain.Company
	for _, c := range companies {
		responseCompanies = append(responseCompanies, c.ToResponseCompany())
	}

	return responseCompanies, nil
}

// GetCompanyByID IDで会社を取得
func (uc *CompanyUsecase) GetCompanyByID(id uint) (*domain.Company, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid company id: %d", id)
	}

	company, err := uc.companyRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get company by id %d: %w", id, err)
	}

	// レスポンス用データに変換
	responseCompany := company.ToResponseCompany()
	return &responseCompany, nil
}

// CreateCompany 会社を作成
func (uc *CompanyUsecase) CreateCompany(name, email, phone, address, website, description string) (*domain.Company, error) {
	// 入力バリデーション
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// メール重複チェック
	exists, err := uc.companyRepo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("company with email %s already exists", email)
	}

	// 会社オブジェクト作成
	company := &domain.Company{
		Name:        name,
		Email:       email,
		Phone:       phone,
		Address:     address,
		Website:     website,
		Description: description,
	}

	// バリデーション
	if !company.IsValidCompany() {
		return nil, fmt.Errorf("invalid company data")
	}

	// リポジトリで保存
	if err := uc.companyRepo.Create(company); err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	// レスポンス用データに変換
	responseCompany := company.ToResponseCompany()
	return &responseCompany, nil
}

// UpdateCompany 会社を更新
func (uc *CompanyUsecase) UpdateCompany(id uint, name, email, phone, address, website, description string) (*domain.Company, error) {
	// 入力バリデーション
	if id == 0 {
		return nil, fmt.Errorf("invalid company id: %d", id)
	}
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// 既存会社取得
	existingCompany, err := uc.companyRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get company for update: %w", err)
	}

	// メール変更時の重複チェック
	if existingCompany.Email != email {
		exists, err := uc.companyRepo.ExistsByEmail(email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("company with email %s already exists", email)
		}
	}

	// 会社情報更新
	existingCompany.Name = name
	existingCompany.Email = email
	existingCompany.Phone = phone
	existingCompany.Address = address
	existingCompany.Website = website
	existingCompany.Description = description

	// バリデーション
	if !existingCompany.IsValidCompany() {
		return nil, fmt.Errorf("invalid company data")
	}

	// リポジトリで更新
	if err := uc.companyRepo.Update(existingCompany); err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	// レスポンス用データに変換
	responseCompany := existingCompany.ToResponseCompany()
	return &responseCompany, nil
}

// DeleteCompany 会社を削除
func (uc *CompanyUsecase) DeleteCompany(id uint) error {
	// 入力バリデーション
	if id == 0 {
		return fmt.Errorf("invalid company id: %d", id)
	}

	// 会社存在確認
	exists, err := uc.companyRepo.Exists(id)
	if err != nil {
		return fmt.Errorf("failed to check company existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("company with id %d not found", id)
	}

	// リポジトリで削除（CASCADE設定により関連データも自動削除）
	if err := uc.companyRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	return nil
}

// GetCompaniesPaginated ページネーション付きで会社を取得
func (uc *CompanyUsecase) GetCompaniesPaginated(page, limit int) ([]domain.Company, *helper.PaginationResponse, error) {
	// ページネーションパラメータ正規化
	paginationReq := &helper.PaginationRequest{
		Page:  page,
		Limit: limit,
	}
	offset := paginationReq.GetOffset()
	normalizedLimit := paginationReq.GetLimit()

	// 総件数取得
	total, err := uc.companyRepo.Count()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count companies: %w", err)
	}

	// ページネーション付きで会社取得
	companies, err := uc.companyRepo.GetPaginated(offset, normalizedLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get paginated companies: %w", err)
	}

	// レスポンス用データに変換
	var responseCompanies []domain.Company
	for _, c := range companies {
		responseCompanies = append(responseCompanies, c.ToResponseCompany())
	}

	// ページネーション情報作成
	pagination := helper.NewPaginationResponse(paginationReq.Page, normalizedLimit, total)

	return responseCompanies, pagination, nil
}

// SearchCompanies 名前で会社を検索
func (uc *CompanyUsecase) SearchCompanies(name string) ([]domain.Company, error) {
	if name == "" {
		return nil, fmt.Errorf("search name is required")
	}

	companies, err := uc.companyRepo.SearchByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed to search companies: %w", err)
	}

	// レスポンス用データに変換
	var responseCompanies []domain.Company
	for _, c := range companies {
		responseCompanies = append(responseCompanies, c.ToResponseCompany())
	}

	return responseCompanies, nil
}

// AddUserToCompany ユーザーを会社に追加
func (uc *CompanyUsecase) AddUserToCompany(userID, companyID uint, role string) (*domain.CompanyUser, error) {
	// 入力バリデーション
	if userID == 0 {
		return nil, fmt.Errorf("invalid user id: %d", userID)
	}
	if companyID == 0 {
		return nil, fmt.Errorf("invalid company id: %d", companyID)
	}
	if role == "" {
		role = "member" // デフォルト役割
	}

	// 会社存在確認
	exists, err := uc.companyRepo.Exists(companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to check company existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("company with id %d not found", companyID)
	}

	// 既存関係チェック
	relationExists, err := uc.companyUserRepo.Exists(userID, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to check relation existence: %w", err)
	}
	if relationExists {
		return nil, fmt.Errorf("user %d is already associated with company %d", userID, companyID)
	}

	// 関係作成
	companyUser := &domain.CompanyUser{
		UserID:    userID,
		CompanyID: companyID,
		Role:      role,
	}

	if err := uc.companyUserRepo.Create(companyUser); err != nil {
		return nil, fmt.Errorf("failed to add user to company: %w", err)
	}

	return companyUser, nil
}

// UpdateUserRole ユーザーの役割を更新
func (uc *CompanyUsecase) UpdateUserRole(userID, companyID uint, role string) (*domain.CompanyUser, error) {
	// 入力バリデーション
	if userID == 0 {
		return nil, fmt.Errorf("invalid user id: %d", userID)
	}
	if companyID == 0 {
		return nil, fmt.Errorf("invalid company id: %d", companyID)
	}
	if role == "" {
		return nil, fmt.Errorf("role is required")
	}

	// 既存関係取得
	companyUser, err := uc.companyUserRepo.GetRelation(userID, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user-company relation: %w", err)
	}

	// 役割更新
	companyUser.Role = role
	if err := uc.companyUserRepo.Update(companyUser); err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	return companyUser, nil
}

// RemoveUserFromCompany ユーザーを会社から削除
func (uc *CompanyUsecase) RemoveUserFromCompany(userID, companyID uint) error {
	// 入力バリデーション
	if userID == 0 {
		return fmt.Errorf("invalid user id: %d", userID)
	}
	if companyID == 0 {
		return fmt.Errorf("invalid company id: %d", companyID)
	}

	// 関係存在確認
	exists, err := uc.companyUserRepo.Exists(userID, companyID)
	if err != nil {
		return fmt.Errorf("failed to check relation existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user %d is not associated with company %d", userID, companyID)
	}

	// 関係削除
	if err := uc.companyUserRepo.Delete(userID, companyID); err != nil {
		return fmt.Errorf("failed to remove user from company: %w", err)
	}

	return nil
}

// GetUsersByCompany 会社のユーザー一覧を取得
func (uc *CompanyUsecase) GetUsersByCompany(companyID uint) ([]domain.CompanyUser, error) {
	if companyID == 0 {
		return nil, fmt.Errorf("invalid company id: %d", companyID)
	}

	// 会社存在確認
	exists, err := uc.companyRepo.Exists(companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to check company existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("company with id %d not found", companyID)
	}

	companyUsers, err := uc.companyUserRepo.GetUsersByCompanyID(companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by company: %w", err)
	}

	return companyUsers, nil
}

// GetCompaniesByUser ユーザーの会社一覧を取得
func (uc *CompanyUsecase) GetCompaniesByUser(userID uint) ([]domain.CompanyUser, error) {
	if userID == 0 {
		return nil, fmt.Errorf("invalid user id: %d", userID)
	}

	companyUsers, err := uc.companyUserRepo.GetCompaniesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get companies by user: %w", err)
	}

	return companyUsers, nil
}
