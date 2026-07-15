package model

type RoleMenusList struct {
	MenuID      string `gorm:"column:menu_id" json:"menu_id"`
	MenuTitle   string `gorm:"column:menu_title" json:"menu_title"`
	Level       int    `gorm:"column:level" json:"level"`
	FormPos     int    `gorm:"column:form_pos" json:"form_pos"`
	FormClass   string `gorm:"column:form_class" json:"form_class"`
	IconIndex   int    `gorm:"column:icon_index" json:"icon_index"`
	IsHeader    bool   `gorm:"column:is_header" json:"is_header"`
	Breadcrumbs string `gorm:"column:breadcrumbs" json:"breadcrumbs"`
	Shortcut    string `gorm:"column:shortcut" json:"shortcut"`
	TrCode      string `gorm:"column:tr_code" json:"tr_code"`
	Params      string `gorm:"column:params" json:"params"`
	SortIndex   string `gorm:"column:sort_index" json:"sort_index"`
	ParentID    string `gorm:"column:parent_id" json:"parent_id"`
}

func (RoleMenusList) TableName() string {
	return "sys.role_menus"
}

type HakAksesList struct {
	Id   string `gorm:"column:package_id" json:"package_id"`
	Name string `gorm:"column:package_name" json:"package_name"`
	File string `gorm:"column:package_file" json:"package_file"`
}

func (HakAksesList) TableName() string {
	return "sys.role_menus"
}
