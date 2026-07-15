package repository

import (
	"context"
	"errors"
	"log"
	"math"
	"mobile/entity"
	"mobile/model"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryUserImpl struct {
		*gorm.DB
	}
)
type UserRepository interface {
	FindOneByEmailAndCustId(email, custId string) (model.User, error)
	FindOneByEmailCustIdEmpId(email, custId string, empId int64) (model.UserRead, error)
	Store(c context.Context, data *model.User) error
	FindOneByUserID(userID int64) (model.User, error)
	Update(c context.Context, userID int64, data model.User) error
	FindDetail(userID int64, custId string) (Details model.User, err error)
	Delete(c context.Context, custId string, UserId int, deletedBy int64) error
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.User, int64, int, error)
	UpdateFcmToken(c context.Context, userID int64, fcmToken string) error
	UpdatePassword(c context.Context, userID int64, password string) error
}

func NewUserRepository(db *gorm.DB) *RepositoryUserImpl {
	return &RepositoryUserImpl{db}
}

func (repo *RepositoryUserImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryUserImpl) FindOneByEmailAndCustId(email, custId string) (model.User, error) {
	user := model.User{}
	err := repository.
		Where("is_del = ? AND email = ? AND cust_id = ?", false, email, custId).
		Take(&user).Error

	if err != nil {
		log.Println("err.Error():", err.Error())
		return user, err
	}

	return user, nil
}

func (repository *RepositoryUserImpl) FindOneByEmailCustIdEmpId(email, custId string, empId int64) (model.UserRead, error) {
	user := model.UserRead{}

	err := repository.
		Select("sys.m_user.*, r.role_name").
		Joins("join sys.user_roles mu on mu.user_id = sys.m_user.user_id AND mu.cust_id = sys.m_user.cust_id").
		Joins("join sys.m_role r on mu.role_id = r.role_id").
		Where("sys.m_user.is_del = ? AND sys.m_user.email = ? AND sys.m_user.cust_id = ? AND sys.m_user.emp_id = ? AND r.cust_id=?", false, email, custId, empId, custId).
		Take(&user).Error

	if err != nil {
		log.Println("err.Error():", err.Error())
		return user, err
	}

	return user, nil
}

func (repository *RepositoryUserImpl) FindOneByUserID(userID int64) (model.User, error) {
	user := model.User{}

	err := repository.
		Where("user_id=?", userID).
		Take(&user).Error
	if err != nil {
		log.Println("err.Error():", err.Error())
		return user, err
	}

	return user, nil
}

func (repository *RepositoryUserImpl) Store(c context.Context, data *model.User) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryUserImpl) Update(c context.Context, userID int64, data model.User) error {
	result := repository.model(c).Model(&data).Where("user_id=?", userID).Updates(&data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryUserImpl) FindDetail(userID int64, custId string) (Details model.User, err error) {
	err = repository.
		Where("user_id = ? AND cust_id=?", userID, custId).
		Take(&Details).Error
	return Details, err
}

func (repository *RepositoryUserImpl) Delete(c context.Context, custId string, UserId int, deletedBy int64) error {
	var data model.User
	result := repository.model(c).Model(&data).Where("user_id=? AND cust_id = ? AND is_del= ? ", UserId, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryUserImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.User, int64, int, error) {
	var data []model.User
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("user_id")
	query := repository.Select("*")

	if dataFilter.Query != "" {

	}
	if dataFilter.Sort != "" {

	} else {
		query.Order("user_id DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&data).Error
	if err != nil {
		return data, total, 0, err
	}
	err = queryCount.Model(&data).Count(&total).Error
	if err != nil {
		return data, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return data, total, lastPage, nil
}
func (repository *RepositoryUserImpl) UpdateFcmToken(c context.Context, userID int64, fcmToken string) error {
	var userModel model.User
	result := repository.model(c).Model(&userModel).Where("user_id=?", userID).Update("fcm_token", fcmToken)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryUserImpl) UpdatePassword(c context.Context, userID int64, password string) error {
	var userModel model.User
	passHash, err := model.HashPasswordString(password)
	if err != nil {
		return err
	}
	result := repository.model(c).Model(&userModel).Where("user_id=?", userID).Update("user_pass", passHash)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
