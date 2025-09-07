package helper

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// CustomValidator カスタムバリデーター
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator バリデーターを初期化
func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate バリデーションを実行
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return FormatValidationErrors(err)
	}
	return nil
}

// ValidationError カスタムバリデーションエラー
type ValidationError struct {
	Field   string `json:"field"`   // フィールド名
	Tag     string `json:"tag"`     // バリデーションタグ
	Value   string `json:"value"`   // 実際の値
	Message string `json:"message"` // エラーメッセージ
}

// ValidationErrors バリデーションエラーのスライス
type ValidationErrors []ValidationError

// Error エラーインターフェースの実装
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// FormatValidationErrors バリデーションエラーをフォーマット
func FormatValidationErrors(err error) ValidationErrors {
	var validationErrors ValidationErrors

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErr {
			ve := ValidationError{
				Field: strings.ToLower(fieldErr.Field()),
				Tag:   fieldErr.Tag(),
				Value: fmt.Sprintf("%v", fieldErr.Value()),
				Message: generateErrorMessage(fieldErr),
			}
			validationErrors = append(validationErrors, ve)
		}
	}

	return validationErrors
}

// generateErrorMessage フィールドエラーからメッセージを生成
func generateErrorMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%sは必須項目です", field)
	case "email":
		return fmt.Sprintf("%sは有効なメールアドレスを入力してください", field)
	case "min":
		return fmt.Sprintf("%sは%s文字以上で入力してください", field, fe.Param())
	case "max":
		return fmt.Sprintf("%sは%s文字以下で入力してください", field, fe.Param())
	case "len":
		return fmt.Sprintf("%sは%s文字で入力してください", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%sは次の値のいずれかを指定してください: %s", field, fe.Param())
	case "url":
		return fmt.Sprintf("%sは有効なURLを入力してください", field)
	case "numeric":
		return fmt.Sprintf("%sは数字を入力してください", field)
	case "alpha":
		return fmt.Sprintf("%sはアルファベットのみ入力可能です", field)
	case "alphanum":
		return fmt.Sprintf("%sは英数字のみ入力可能です", field)
	default:
		return fmt.Sprintf("%sの入力値が正しくありません", field)
	}
}

// GetValidationErrorDetails バリデーションエラーの詳細を文字列で取得
func GetValidationErrorDetails(err error) string {
	if validationErrors, ok := err.(ValidationErrors); ok {
		var details []string
		for _, ve := range validationErrors {
			details = append(details, fmt.Sprintf("フィールド'%s': %s", ve.Field, ve.Message))
		}
		return strings.Join(details, ", ")
	}
	return err.Error()
}

// ValidateStruct 構造体をバリデーション
func ValidateStruct(s interface{}) error {
	validator := NewValidator()
	return validator.Validate(s)
}