package op

import (
	"context"
	"errors"
	"main/api/rest/inputStructs"
	"main/core/node"
	"main/core/utils"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

func RegisterCamRoutes(api huma.API, root *node.Node) {
	huma.Register(api, huma.Operation{
		OperationID:   "get-camera",
		Method:        http.MethodGet,
		Path:          "/cameras/{cameraId}",
		Summary:       "Get a camera",
		Description:   "Get a camera",
		Tags:          []string{"Groups/Cameras"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.CamPathId) (*inputStructs.CamResp, error) {
		groups, err := root.GetChild("groups")
		if err != nil {
			return nil, err
		}

		var cam, camGroup *node.Node

		for _, group := range groups.GetChilds() {
			if cam, err = group.GetChild(input.CamId); err == nil {
				camGroup = group
				break
			}
		}

		if cam == nil {
			return nil, node.ErrNodeNotExists
		}

		resp := &inputStructs.CamResp{}
		group := inputStructs.BaseGroup{}

		if err = mapstructure.Decode(cam.Get(), &resp.Body); err != nil {
			return nil, err
		}
		if err = mapstructure.Decode(camGroup.Get(), &group); err != nil {
			return nil, err
		}

		resp.Body.Group = group

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "post-cam",
		Method:        http.MethodPost,
		Path:          "/groups/{groupId}",
		Summary:       "Create a camera",
		Tags:          []string{"Groups/Cameras"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *inputStructs.CamPost) (*inputStructs.CamResp, error) {
		groups, err := root.GetChild("groups")
		if err != nil {
			return nil, err
		}

		group, err := groups.GetChild(input.GroupId)
		if err != nil {
			return nil, err
		}

		resp := &inputStructs.CamResp{}
		if err = mapstructure.Decode(group.Get(), &resp.Body.Group); err != nil {
			return nil, err
		}

		modules, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}

		for _, module := range input.Body.Modules {
			if _, err := modules.GetChild(module); err != nil {
				return nil, err
			}
		}

		if input.Body.Mask == nil {
			input.Body.Mask = []inputStructs.Point{}
		} else if len(input.Body.Mask) > 0 && len(input.Body.Mask) < 3 {
			return nil, errors.New("mask should contain at least 3 points")
		}

		camData := map[string]any{}
		camData["id"] = uuid.New().String()

		cam := group.CreateChild(camData["id"].(string))
		if err = cam.Init(); err != nil {
			return nil, err
		}

		if err = mapstructure.Decode(input.Body, &camData); err != nil {
			return nil, err
		}
		if err = cam.Put(camData); err != nil {
			return nil, err
		}
		if err = mapstructure.Decode(cam.Get(), &resp.Body); err != nil {
			return nil, err
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "patch-cam",
		Method:        http.MethodPatch,
		Path:          "/cameras/{cameraId}",
		Summary:       "Patch a camera",
		Tags:          []string{"Groups/Cameras"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.CamPatch) (*inputStructs.CamResp, error) {
		groups, err := root.GetChild("groups")
		if err != nil {
			return nil, err
		}

		var cam, camGroup *node.Node

		for _, group := range groups.GetChilds() {
			if cam, err = group.GetChild(input.CamId); err == nil {
				camGroup = group
				break
			}
		}

		if cam == nil {
			return nil, node.ErrNodeNotExists
		}

		modules, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}

		for _, module := range input.Body.Modules {
			if _, err := modules.GetChild(module); err != nil {
				return nil, err
			}
		}

		if len(input.Body.Mask) > 0 && len(input.Body.Mask) < 3 {
			return nil, errors.New("mask should contain at least 3 points")
		}

		var newCamData map[string]any
		if err = mapstructure.Decode(input.Body, &newCamData); err != nil {
			return nil, err
		}

		if input.Body.Name == "" {
			delete(newCamData, "name")
		}
		if input.Body.Modules == nil {
			delete(newCamData, "modules")
		}
		if input.Body.Mask == nil {
			delete(newCamData, "mask")
		}

		if err = cam.Put(utils.MergeMap(cam.Get(), newCamData)); err != nil {
			return nil, err
		}

		resp := &inputStructs.CamResp{}
		if err = mapstructure.Decode(cam.Get(), &resp.Body); err != nil {
			return nil, err
		}
		if err = mapstructure.Decode(camGroup.Get(), &resp.Body.Group); err != nil {
			return nil, err
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-cam",
		Method:        http.MethodDelete,
		Path:          "/cameras/{cameraId}",
		Summary:       "Delete a camera",
		Tags:          []string{"Groups/Cameras"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *inputStructs.CamPathId) (*struct{}, error) {
		groups, err := root.GetChild("groups")
		if err != nil {
			return nil, err
		}

		for _, group := range groups.GetChilds() {
			if err = group.RemoveChild(input.CamId); err == nil {
				return nil, nil
			}
		}

		return nil, node.ErrNodeNotExists
	})

}
