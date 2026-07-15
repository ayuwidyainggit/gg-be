package repository

import (
	"context"
	"errors"
	"time"

	"mobile/model"
	"mobile/pkg/constant"

	"gorm.io/gorm"
)

type PasswordResetRequestRepository interface {
	ExpirePendingByUserID(ctx context.Context, userID int64) error
	Create(ctx context.Context, row *model.PasswordResetRequest) error
	FindByRequestID(ctx context.Context, requestID string) (model.PasswordResetRequest, error)
	FindByResetToken(ctx context.Context, resetTokenHMAC string) (model.PasswordResetRequest, error)
	Update(ctx context.Context, id int64, updates map[string]interface{}) error
}

type passwordResetRequestRepoImpl struct {
	*gorm.DB
}

func NewPasswordResetRequestRepository(db *gorm.DB) *passwordResetRequestRepoImpl {
	return &passwordResetRequestRepoImpl{db}
}

func (r *passwordResetRequestRepoImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.WithContext(ctx)
}

func (r *passwordResetRequestRepoImpl) ExpirePendingByUserID(ctx context.Context, userID int64) error {
	res := r.model(ctx).Model(&model.PasswordResetRequest{}).
		Where("user_id = ? AND status IN ?", userID, []string{
			constant.PasswordResetStatusPendingOTP,
			constant.PasswordResetStatusOTPValidated,
		}).
		Updates(map[string]interface{}{
			"status":     constant.PasswordResetStatusExpired,
			"updated_at": time.Now(),
		})
	return res.Error
}

func (r *passwordResetRequestRepoImpl) Create(ctx context.Context, row *model.PasswordResetRequest) error {
	return r.model(ctx).Create(row).Error
}

func (r *passwordResetRequestRepoImpl) FindByRequestID(ctx context.Context, requestID string) (model.PasswordResetRequest, error) {
	var row model.PasswordResetRequest
	err := r.model(ctx).Where("request_id = ?", requestID).Take(&row).Error
	if err != nil {
		return row, err
	}
	return row, nil
}

func (r *passwordResetRequestRepoImpl) FindByResetToken(ctx context.Context, resetTokenHMAC string) (model.PasswordResetRequest, error) {
	var row model.PasswordResetRequest
	err := r.model(ctx).Where("reset_token = ?", resetTokenHMAC).Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return row, err
		}
		return row, err
	}
	return row, nil
}

func (r *passwordResetRequestRepoImpl) Update(ctx context.Context, id int64, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	res := r.model(ctx).Model(&model.PasswordResetRequest{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
