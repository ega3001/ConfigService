package inputStructs

// Cameras
type CameraData struct {
	Name     string   `json:"name,omitempty" mapstructure:"name" contentType:"application/json" doc:"Human-friendly Camera name"`
	Modules  []string `json:"modules,omitempty" mapstructure:"modules" example:"LPR" doc:"Collection of registered camera modules"`
	Mask     []Point  `json:"mask,omitempty" mapstructure:"mask" doc:"Polygon that describes detection zone"`
	Source   string   `json:"source,omitempty" mapstructure:"source" doc:"RTSP URI"`
	Reserved bool     `json:"reserved,omitempty" mapstructure:"reserved" doc:"Read-only utility field"`
}

type BaseCamera struct {
	Id         string `json:"id" example:"xxxx-xxxx-xxxx-xxxx" doc:"Camera Id"`
	CameraData `mapstructure:",squash"`
}

type CameraWithGroup struct {
	BaseCamera `mapstructure:",squash"`
	Group      BaseGroup `json:"group"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Path
type CamPathId struct {
	CamId string `path:"cameraId" example:"xxxx-xxxx-xxxx-xxxx" doc:"Camera Id in UUID format"`
}

// Post/Patch
type CamPost struct {
	GroupPathId
	Body CameraData
}

type CamPatch struct {
	CamPathId
	Body CameraData
}

// Responses
type CamResp struct {
	Body CameraWithGroup
}

type CamsResp struct {
	Body struct {
		Cameras []CameraWithGroup `json:"cameras"`
	}
}
