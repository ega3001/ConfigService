package inputStructs

// Group
type GroupData struct {
	Name     string     `json:"name,omitempty" mapstructure:"name" contentType:"application/json"`
	Schedule []Schedule `json:"schedule,omitempty" mapstructure:"schedule"`
}

type IsSecure struct {
	IsSecure bool `json:"isSecure" mapstructure:"isSecure" doc:"True if group is on"`
}

type BaseGroup struct {
	Id        string `json:"id" mapstructure:"id" example:"xxxx-xxxx-xxxx-xxxx" doc:"Group Id"`
	GroupData `mapstructure:",squash"`
	IsSecure  `mapstructure:",squash"`
}

type GroupWithCams struct {
	BaseGroup `mapstructure:",squash"`
	Cameras   []BaseCamera `json:"cameras" doc:"Array of cameras"`
}

// Schedule
type Schedule struct {
	StartTime string `json:"startTime" mapstructure:"startTime" example:"22:30"`
	EndTime   string `json:"endTime" mapstructure:"endTime" example:"08:30"`
}

// Path
type GroupPathId struct {
	GroupId string `path:"groupId" example:"xxxx-xxxx-xxxx-xxxx" doc:"Group Id in UUID format"`
}

// Post/Patch
type GroupPost struct {
	Body GroupData
}

type GroupPatch struct {
	GroupPathId
	Body GroupData
}

type SecurePatch struct {
	GroupPathId
	Body IsSecure
}

// Responses
type GroupResp struct {
	Body GroupWithCams
}

type GroupsResp struct {
	Body struct {
		Groups []GroupWithCams `json:"groups"`
	}
}
