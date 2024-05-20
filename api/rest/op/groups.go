package op

import (
	"context"
	"errors"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"main/api/rest/inputStructs"
	"main/core/node"
	"main/core/utils"
	"net/http"
	"regexp"
)

func RegisterGroupRoutes(api huma.API, root *node.Node) {
	huma.Register(api, huma.Operation{
		OperationID:   "get-all-groups",
		Method:        http.MethodGet,
		Path:          "/groups",
		Summary:       "Get all groups",
		Description:   "Get all groups",
		Tags:          []string{"Groups"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *struct{}) (*inputStructs.GroupsResp, error) {
		groups, err := root.GetChild("groups")
		if err != nil {
			return nil, err
		}

		resp := &inputStructs.GroupsResp{}
		resp.Body.Groups = make([]inputStructs.GroupWithCams, groups.ChildsAmount())
		for i, group := range groups.GetChilds() {
			err = mapstructure.Decode(group.Get(), &resp.Body.Groups[i])
			if err != nil {
				return nil, err
			}

			cameras := make([]inputStructs.BaseCamera, group.ChildsAmount())
			for j, camera := range group.GetChilds() {
				cam := inputStructs.BaseCamera{}
				err = mapstructure.Decode(camera.Get(), &cam)
				if err != nil {
					return nil, err
				}
				cameras[j] = cam
			}
			resp.Body.Groups[i].Cameras = cameras
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "get-group",
		Method:        http.MethodGet,
		Path:          "/groups/{groupId}",
		Summary:       "Get group",
		Description:   "Get group",
		Tags:          []string{"Groups"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.GroupPathId) (*inputStructs.GroupResp, error) {
		groups, err := root.GetChild("groups")
		if err != nil {
			return nil, err
		}
		group, err := groups.GetChild(input.GroupId)
		if err != nil {
			return nil, err
		}

		resp := &inputStructs.GroupResp{}
		if err = mapstructure.Decode(group.Get(), &resp.Body); err != nil {
			return nil, err
		}

		cameras := make([]inputStructs.BaseCamera, group.ChildsAmount())
		for j, camera := range group.GetChilds() {
			cam := inputStructs.BaseCamera{}
			err = mapstructure.Decode(camera.Get(), &cam)
			if err != nil {
				return nil, err
			}
			cameras[j] = cam
		}
		resp.Body.Cameras = cameras

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "post-group",
		Method:        http.MethodPost,
		Path:          "/groups",
		Summary:       "Create group",
		Description:   "Create group",
		Tags:          []string{"Groups"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *inputStructs.GroupPost) (*inputStructs.GroupResp, error) {
		groups, err := root.GetChild("groups")
		if err != nil {
			return nil, err
		}

		for _, sch := range input.Body.Schedule {
			r := regexp.MustCompile("^(([01][0-9])|(2[0-3])):[0-5][0-9]$")
			if !r.MatchString(sch.StartTime) || !r.MatchString(sch.EndTime) {
				return nil, errors.New("wrong time format, use hh:mm")
			}
		}

		groupId := uuid.New().String()
		group := groups.CreateChild(groupId)
		if err = group.Init(); err != nil {
			return nil, err
		}

		groupData := map[string]any{}
		groupData["id"] = groupId
		groupData["schedule"] = []inputStructs.Schedule{}
		if err = mapstructure.Decode(input.Body, &groupData); err != nil {
			return nil, err
		}
		if err = group.Put(groupData); err != nil {
			return nil, err
		}

		resp := &inputStructs.GroupResp{}
		if err = mapstructure.Decode(group.Get(), &resp.Body); err != nil {
			return nil, err
		}
		resp.Body.Cameras = []inputStructs.BaseCamera{}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "patch-group",
		Method:        http.MethodPatch,
		Path:          "/groups/{groupId}",
		Summary:       "Patch group",
		Description:   "Patch group",
		Tags:          []string{"Groups"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.GroupPatch) (*inputStructs.GroupResp, error) {
		groups, err := root.GetChild("groups")
		if err != nil {
			return nil, err
		}

		group, err := groups.GetChild(input.GroupId)
		if err != nil {
			return nil, err
		}

		for _, sch := range input.Body.Schedule {
			r := regexp.MustCompile("^(([01][0-9])|(2[0-3])):[0-5][0-9]$")
			if !r.MatchString(sch.StartTime) || !r.MatchString(sch.EndTime) {
				return nil, errors.New("wrong time format, use hh:mm")
			}
		}

		var newGroupData map[string]any
		if err = mapstructure.Decode(input.Body, &newGroupData); err != nil {
			return nil, err
		}

		if input.Body.Name == "" {
			delete(newGroupData, "name")
		}
		if input.Body.Schedule == nil {
			delete(newGroupData, "schedule")
		}

		if err = group.Put(utils.MergeMap(group.Get(), newGroupData)); err != nil {
			return nil, err
		}

		resp := &inputStructs.GroupResp{}
		if err = mapstructure.Decode(group.Get(), &resp.Body); err != nil {
			return nil, err
		}

		resp.Body.Cameras = make([]inputStructs.BaseCamera, group.ChildsAmount())
		for j, camera := range group.GetChilds() {
			cam := inputStructs.BaseCamera{}
			err = mapstructure.Decode(camera.Get(), &cam)
			if err != nil {
				return nil, err
			}
			resp.Body.Cameras[j] = cam
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-group",
		Method:        http.MethodDelete,
		Path:          "/groups/{groupId}",
		Summary:       "Delete group",
		Description:   "Delete group",
		Tags:          []string{"Groups"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *inputStructs.GroupPathId) (*struct{}, error) {
		groups, err := root.GetChild("groups")
		if err != nil {
			return nil, err
		}

		if err := groups.RemoveChild(input.GroupId); err != nil {
			return nil, err
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "set-secure",
		Method:        http.MethodPatch,
		Path:          "/groups/{groupId}/secure",
		Summary:       "Set group on secure",
		Description:   "Set group on secure",
		Tags:          []string{"Groups"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.SecurePatch) (*inputStructs.GroupResp, error) {
		groups, err := root.GetChild("groups")
		if err != nil {
			return nil, err
		}
		group, err := groups.GetChild(input.GroupId)
		if err != nil {
			return nil, err
		}

		var newGroupData = map[string]any{}
		if err = mapstructure.Decode(input.Body, &newGroupData); err != nil {
			return nil, err
		}
		if err = group.Put(utils.MergeMap(group.Get(), newGroupData)); err != nil {
			return nil, err
		}
		resp := &inputStructs.GroupResp{}
		if err = mapstructure.Decode(group.Get(), &resp.Body); err != nil {
			return nil, err
		}

		resp.Body.Cameras = make([]inputStructs.BaseCamera, group.ChildsAmount())
		for j, camera := range group.GetChilds() {
			cam := inputStructs.BaseCamera{}
			err = mapstructure.Decode(camera.Get(), &cam)
			if err != nil {
				return nil, err
			}
			resp.Body.Cameras[j] = cam
		}

		return resp, nil
	})
}
