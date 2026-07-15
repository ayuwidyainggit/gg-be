package repository

import (
	"system/model"

	"gorm.io/gorm"
)

type (
	RepositoryMMenuImpl struct {
		*gorm.DB
	}
)
type MMenuRepository interface {
	FindAllMenu(userId int64, custId string) ([]model.MMenu, error)
	FindAllMenuDesktop(userID int64) ([]model.RoleMenusList, error)
	FindAllMenuWithoutCustId(menuInt int) ([]model.MMenuNoToken, error)
	FindAllHakAksesPkg(userID int64) ([]model.HakAksesList, error)
}

func NewMMenuRepository(db *gorm.DB) *RepositoryMMenuImpl {
	return &RepositoryMMenuImpl{db}
}

// func (repo *RepositoryMMenuImpl) model(ctx context.Context) *gorm.DB {
// 	tx := extractTx(ctx)
// 	if tx != nil {
// 		return tx.WithContext(ctx)
// 	}
// 	return repo.WithContext(ctx)
// }

func (repository *RepositoryMMenuImpl) FindAllMenu(userID int64, custId string) ([]model.MMenu, error) {
	var data []model.MMenu

	// query := repository.Select("*").Where("menu_type in (0,2)").Order("level ASC, sort_index ASC ")
	query := repository.Select(`DISTINCT 
			m.menu_id, m.menu_name, 
			m.level, m.sort_index, 
			m.parent_id, m.menu_action, m.menu_type,
			m.icon_web, m.is_header, m.url_web, m."targetType", 
			ml.menu_title,
			mtn.prefix AS tr_code,
			mtn2.prefix AS tr_code2`).
		Joins("JOIN sys.m_menu m ON (sys.role_menus.menu_id=m.menu_id)").
		Joins("JOIN sys.m_user us ON (us.user_id=?)", userID).
		Joins("JOIN sys.menu_lang ml ON ((m.menu_id=ml.menu_id) AND (ml.lang_id=us.lang_id))").
		Joins("LEFT JOIN sys.m_trans_no mtn ON mtn.tr_code = m.tr_code AND mtn.cust_id = ?", custId).
		Joins("LEFT JOIN sys.m_trans_no mtn2 ON mtn2.tr_code = m.tr_code2 AND mtn.cust_id = ?", custId).
		Where("m.menu_type IN (0, 2) AND ( sys.role_menus.role_id IN (SELECT role_id FROM sys.user_roles WHERE user_id=? ) AND sys.role_menus.cust_id IN (SELECT cust_id FROM sys.user_roles WHERE user_id=? ) ) ", userID, userID).
		Order("m.level, m.sort_index")
	err := query.Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (repository *RepositoryMMenuImpl) FindAllMenuDesktop(userID int64) ([]model.RoleMenusList, error) {
	var data []model.RoleMenusList

	query := repository.Select("DISTINCT m.menu_id, ml.menu_title, m.level, m.form_pos, m.form_class,m.icon_index, m.is_header, m.breadcrumbs, m.tr_code, m.params, m.shortcut, m.sort_index,m.parent_id").
		Joins("JOIN sys.m_menu m ON (sys.role_menus.menu_id=m.menu_id)").
		Joins("JOIN sys.menu_lang ml ON ((m.menu_id=ml.menu_id) AND (ml.lang_id='id'))").
		Where("m.menu_type IN (0, 1) AND ( sys.role_menus.role_id IN (SELECT role_id FROM sys.user_roles WHERE user_id=? ) AND sys.role_menus.cust_id IN (SELECT cust_id FROM sys.user_roles WHERE user_id=? ) ) ", userID, userID).
		Order("m.level, m.sort_index")
	err := query.Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (repository *RepositoryMMenuImpl) FindAllMenuWithoutCustId(menuInt int) ([]model.MMenuNoToken, error) {
	var data []model.MMenuNoToken

	// query := repository.Select("*").Where("menu_type in (0,2)").Order("level ASC, sort_index ASC ")
	query := repository.Select(`DISTINCT 
			m.menu_id, m.menu_name, 
			m.level, m.sort_index, 
			m.parent_id, m.menu_action, m.menu_type,
			m.icon_web, m.is_header, m.url_web, m."targetType", 
			ml.menu_title`).
		Table("sys.m_menu AS m").
		Joins("JOIN sys.menu_lang ml ON m.menu_id=ml.menu_id")
	if menuInt != 0 {
		query = query.Where("m.menu_type = ?", menuInt)
	}

	query = query.Order("m.level, m.sort_index")

	err := query.Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (repository *RepositoryMMenuImpl) FindAllHakAksesPkg(userID int64) ([]model.HakAksesList, error) {
	var data []model.HakAksesList

	query := repository.Select("DISTINCT pc.package_id, pc.package_name, pc.package_file ").
		Joins("JOIN sys.m_menu m ON (sys.role_menus.menu_id=m.menu_id)").
		Joins("JOIN sys.m_package pc ON (m.package_id=pc.package_id)").
		Where("m.menu_type IN (0, 1) AND sys.role_menus.role_id IN (SELECT role_id FROM sys.user_roles WHERE user_id=?)  ", userID)
	err := query.Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}
