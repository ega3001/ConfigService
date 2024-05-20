package op

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/mitchellh/mapstructure"
	"main/api/rest/inputStructs"
	"main/core/node"
	"main/core/utils"
	"net/http"
)

func RegisterMediaRoutes(api huma.API, root *node.Node) {
	huma.Register(api, huma.Operation{
		OperationID:   "get-media",
		Method:        http.MethodGet,
		Path:          "/system/media",
		Summary:       "Get media config",
		Description:   "Get media config",
		Tags:          []string{"System/Media"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *struct{}) (*inputStructs.MediaResp, error) {
		system, err := root.GetChild("system")
		if err != nil {
			return nil, err
		}
		mediaNode, err := system.GetChild("media")
		if err != nil {
			return nil, err
		}

		media := &inputStructs.Media{}
		if err = mapstructure.Decode(mediaNode.Get(), &media); err != nil {
			return nil, err
		}

		resp := &inputStructs.MediaResp{Body: *media}
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "patch-media",
		Method:        http.MethodPatch,
		Path:          "/system/media",
		Summary:       "Patch media config",
		Description:   "Patch media config",
		Tags:          []string{"System/Media"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.MediaInput) (*inputStructs.MediaResp, error) {
		system, err := root.GetChild("system")
		if err != nil {
			return nil, err
		}
		mediaNode, err := system.GetChild("media")
		if err != nil {
			return nil, err
		}

		newMediaData := map[string]any{}
		if err := mapstructure.Decode(input.Body, &newMediaData); err != nil {
			return nil, err
		}
		if err = mediaNode.Put(utils.MergeMap(mediaNode.Get(), newMediaData)); err != nil {
			return nil, err
		}

		media := &inputStructs.Media{}
		if err = mapstructure.Decode(mediaNode.Get(), &media); err != nil {
			return nil, err
		}
		resp := &inputStructs.MediaResp{Body: *media}

		return resp, nil
	})
}
