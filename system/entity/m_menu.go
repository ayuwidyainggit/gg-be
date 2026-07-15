package entity

type WebMenuResp struct {
	Id         string        `json:"id"`
	ParentId   string        `json:"parent_id"`
	Text       string        `json:"text"`
	Icon       string        `json:"icon"`
	Url        string        `json:"url,omitempty"`
	TargetType string        `json:"targetType"`
	TrCode     string        `json:"tr_code"`
	TrCode2    string        `json:"tr_code2"`
	IsHeader   bool          `json:"isHeader,omitempty"`
	Children   []WebMenuResp `json:"children,omitempty"`
}

type WebMenuMap struct {
	Db   map[string][]WebMenuResp
	Resp []WebMenuResp
}

func NewWebMenuMap() *WebMenuMap {
	return &WebMenuMap{
		Db: make(map[string][]WebMenuResp),
	}
}
func (w *WebMenuMap) SetChildrenRecursively(res *WebMenuResp) {
	// append to make a copy of the slice (otherwise we will be changing items in the 'database')
	res.Children = append([]WebMenuResp{}, w.Db[res.Id]...) // Get the children from simulated database

	for i := range res.Children {
		w.SetChildrenRecursively(&res.Children[i])
	}
}
func (w *WebMenuMap) Append(resp WebMenuResp) {
	w.Resp = append(w.Resp, resp)
}

type DesktopMenuResp struct {
	Menu    []DesktopMenu    `json:"menu"`
	Package []DesktopPackage `json:"package"`
}

type DesktopMenu struct {
	MenuID     string        `json:"menu_id"`
	MenuTitle  string        `json:"menu_title"`
	Level      int           `json:"level"`
	FormPos    int           `json:"form_pos"`
	FormClass  string        `json:"form_class"`
	IconIndex  int           `json:"icon_index"`
	MenuAction int           `json:"menu_action"`
	IsHeader   bool          `json:"is_header"`
	Shortcut   string        `json:"shortcut"`
	TrCode     string        `json:"tr_code"`
	Params     string        `json:"params"`
	Children   []DesktopMenu `json:"children,omitempty"`
}
type DesktopPackage struct {
	PackageID   string `json:"package_id"`
	PackageName string `json:"package_name"`
	PackageFile string `json:"package_file"`
}

type DesktopMenuMap struct {
	Db   map[string][]DesktopMenu
	Resp []DesktopMenu
}

func NewDesktopMenuMap() *DesktopMenuMap {
	return &DesktopMenuMap{
		Db: make(map[string][]DesktopMenu),
	}
}
func (w *DesktopMenuMap) SetChildrenRecursively(res *DesktopMenu) {
	// append to make a copy of the slice (otherwise we will be changing items in the 'database')
	res.Children = append([]DesktopMenu{}, w.Db[res.MenuID]...) // Get the children from simulated database

	for i := range res.Children {
		w.SetChildrenRecursively(&res.Children[i])
	}
}
func (w *DesktopMenuMap) Append(resp DesktopMenu) {
	w.Resp = append(w.Resp, resp)
}
