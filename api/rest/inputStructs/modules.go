package inputStructs

// Base
type BaseList struct {
	List []any `json:"list" mapstructure:"list" contentType:"application/json" required:"true"`
}

type BaseModuleData struct {
}

type BaseModule struct {
	Name           string `json:"name" mapstructure:"name" example:"LPR" doc:"Module name for user"`
	BaseModuleData `mapstructure:",squash"`
}

type BaseMListData struct {
	Name  string `json:"name" mapstructure:"name" contentType:"application/json" required:"true"`
	Alarm bool   `json:"alarm" mapstructure:"alarm" contentType:"application/json" required:"true"`
	Color string `json:"color" mapstructure:"color" contentType:"application/json" required:"true"`
}

type BaseMList struct {
	Id            string `json:"id" mapstructure:"id" example:"xxxx-xxxx-xxxx-xxxx" doc:"List Id"`
	BaseMListData `mapstructure:",squash"`
}

type BaseMListElemData any

type BaseMListElem struct {
	Id     string              `json:"id" mapstructure:"id" example:"xxxx-xxxx-xxxx-xxxx" doc:"List elem Id"`
	Fields []BaseMListElemData `json:"fields"`
}

type ListElemsResp struct {
	Body struct {
		List []map[string]any `json:"list"`
	}
}

// List
type ListResp struct {
	Body BaseList
}

// Path
type ModulePathName struct {
	ModuleName string `path:"moduleName" example:"LPR" doc:"Module name"`
}

type MListPathId struct {
	ListId string `path:"listId" example:"xxxx-xxxx-xxxx-xxxx" doc:"List name as UUID"`
}

type MListElemPathId struct {
	ElemId string `path:"elementId" example:"xxxx-xxxx-xxxx-xxxx" doc:"Element name as UUID"`
}

type ModuleMListPathId struct {
	ModulePathName
	MListPathId
}

type ModuleMListElemPathId struct {
	ModulePathName
	MListPathId
	MListElemPathId
}

// MList
type MList struct {
	BaseMList `mapstructure:",squash"`
	Elems     any `json:"list"`
}

type MListResp struct {
	Body MList
}

type MListsResp struct {
	Body struct {
		List []MList `json:"list" doc:"List of lists"`
	}
}

type MListPost struct {
	ModulePathName
	Body BaseMListData
}

type MListPatch struct {
	ModuleMListPathId
	Body BaseMListData
}

type MListElemsPatch struct {
	ModuleMListPathId
	Body []BaseMListElemData
}

type MListElemPatch struct {
	ModuleMListElemPathId
	Body BaseMListElemData
}

// Module
type ModuleResp struct {
	Body BaseModule
}

type ModulesResp struct {
	Body struct {
		List []BaseModule `json:"list" doc:"List of modules"`
	}
}

type ModulePost struct {
	Body BaseModule
}

type ModulePatch struct {
	Body BaseModuleData
	ModulePathName
}
