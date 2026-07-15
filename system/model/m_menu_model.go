package model

import "gorm.io/gorm"

type MMenu struct {
	MenuID     string `gorm:"column:menu_id" json:"id"`
	MenuName   string `gorm:"column:menu_name" json:"menu_name"`
	MenuTitle  string `gorm:"column:menu_title" json:"text"`
	Level      int    `gorm:"column:level" json:"level"`
	SortIndex  int    `gorm:"column:sort_index" json:"sort_index"`
	PackageID  string `gorm:"column:package_id" json:"package_id"`
	FormPos    int    `gorm:"column:form_pos" json:"form_pos"`
	FormClass  string `gorm:"column:form_class" json:"form_class"`
	IconIndex  int    `gorm:"column:icon_index" json:"icon_index"`
	ParentID   string `gorm:"column:parent_id" json:"parent_id"`
	MenuAction int    `gorm:"column:menu_action" json:"menu_action"`
	MenuTab    int    `gorm:"column:menu_tab" json:"menu_tab"`
	IconWeb    string `gorm:"column:icon_web" json:"icon"`
	UrlWeb     string `gorm:"column:url_web" json:"url"`
	IsHeader   bool   `gorm:"column:is_header" json:"isHeader"`
	TrCode     string `gorm:"column:tr_code" json:"tr_code"`
	TrCode2    string `gorm:"column:tr_code2" json:"tr_code2"`
	TargetType string `gorm:"column:targetType" json:"targetType"`
}

func (MMenu) TableName() string {
	return "sys.role_menus"
}

func (m *MMenu) BeforeCreate(trx *gorm.DB) (err error) {
	// m.CreatedDate = time.Now()
	return nil
}

type MMenuNoToken struct {
	MenuID     string `gorm:"column:menu_id" json:"id"`
	MenuName   string `gorm:"column:menu_name" json:"menu_name"`
	MenuTitle  string `gorm:"column:menu_title" json:"text"`
	Level      int    `gorm:"column:level" json:"level"`
	SortIndex  int    `gorm:"column:sort_index" json:"sort_index"`
	PackageID  string `gorm:"column:package_id" json:"package_id"`
	FormPos    int    `gorm:"column:form_pos" json:"form_pos"`
	FormClass  string `gorm:"column:form_class" json:"form_class"`
	IconIndex  int    `gorm:"column:icon_index" json:"icon_index"`
	ParentID   string `gorm:"column:parent_id" json:"parent_id"`
	MenuAction int    `gorm:"column:menu_action" json:"menu_action"`
	MenuTab    int    `gorm:"column:menu_tab" json:"menu_tab"`
	IconWeb    string `gorm:"column:icon_web" json:"icon"`
	UrlWeb     string `gorm:"column:url_web" json:"url"`
	IsHeader   bool   `gorm:"column:is_header" json:"isHeader"`
	TrCode     string `gorm:"column:tr_code" json:"tr_code"`
	TrCode2    string `gorm:"column:tr_code2" json:"tr_code2"`
	TargetType string `gorm:"column:targetType" json:"targetType"`
}

func (MMenuNoToken) TableName() string {
	return "sys.role_menus"
}

func (m *MMenuNoToken) BeforeCreate(trx *gorm.DB) (err error) {
	// m.CreatedDate = time.Now()
	return nil
}
