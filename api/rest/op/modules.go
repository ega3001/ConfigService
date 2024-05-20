package op

import (
	"context"
	"main/api/rest/inputStructs"
	"main/core/node"
	"main/core/utils"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

func RegisterListRoutes(api huma.API, root *node.Node) {
	huma.Register(api, huma.Operation{
		OperationID:   "get-list-module-names",
		Method:        http.MethodGet,
		Path:          "/modules/names",
		Summary:       "Get all modules names",
		Description:   "Get all modules names",
		Tags:          []string{"Modules"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *struct{}) (*inputStructs.ListResp, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}

		resp := &inputStructs.ListResp{}
		resp.Body.List = utils.ArrToArrAny(modulesN.ListChildKeys())

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "get-list-modules",
		Method:        http.MethodGet,
		Path:          "/modules",
		Summary:       "Get all modules data",
		Description:   "Get all modules with contents",
		Tags:          []string{"Modules"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *struct{}) (*inputStructs.ModulesResp, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}

		resp := &inputStructs.ModulesResp{}
		resp.Body.List = make([]inputStructs.BaseModule, 0, modulesN.ChildsAmount())
		for _, childN := range modulesN.GetChilds() {
			resultM := inputStructs.BaseModule{}
			if err = mapstructure.Decode(childN.Get(), &resultM); err != nil {
				return nil, err
			}
			resp.Body.List = append(
				resp.Body.List,
				resultM,
			)
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "get-module",
		Method:        http.MethodGet,
		Path:          "/modules/{moduleName}",
		Summary:       "Get module data",
		Description:   "Get module contents",
		Tags:          []string{"Modules"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.ModulePathName) (*inputStructs.ModuleResp, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}

		resp := &inputStructs.ModuleResp{}
		if err = mapstructure.Decode(moduleN.Get(), &resp.Body); err != nil {
			return nil, err
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "get-module-list-ids",
		Method:        http.MethodGet,
		Path:          "/modules/{moduleName}/lists/ids",
		Summary:       "Get all module list ids",
		Description:   "Get all module list ids",
		Tags:          []string{"Modules/Lists"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.ModulePathName) (*inputStructs.ListResp, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}
		listsN, err := moduleN.GetChild("lists")
		if err != nil {
			return nil, err
		}

		resp := &inputStructs.ListResp{}
		resp.Body.List = utils.ArrToArrAny(listsN.ListChildKeys())

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "get-module-lists",
		Method:        http.MethodGet,
		Path:          "/modules/{moduleName}/lists",
		Summary:       "Get all module lists data",
		Description:   "Get all module lists with content",
		Tags:          []string{"Modules/Lists"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.ModulePathName) (*inputStructs.MListsResp, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}
		listsN, err := moduleN.GetChild("lists")
		if err != nil {
			return nil, err
		}

		resp := &inputStructs.MListsResp{}
		resp.Body.List = make([]inputStructs.MList, 0, listsN.ChildsAmount())
		for _, listNode := range listsN.GetChilds() {
			list := inputStructs.MList{}
			if err = mapstructure.Decode(listNode.Get(), &list); err != nil {
				return nil, err
			}

			elems := make([]map[string]any, 0, listNode.ChildsAmount())
			for _, elemNode := range listNode.GetChilds() {
				elems = append(elems, elemNode.Get())
			}

			list.Elems = elems

			resp.Body.List = append(resp.Body.List, list)
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "get-module-list",
		Method:        http.MethodGet,
		Path:          "/modules/{moduleName}/{listId}",
		Summary:       "Get module list content",
		Description:   "Get module list content",
		Tags:          []string{"Modules/Lists"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.ModuleMListPathId) (*inputStructs.MListResp, error) {
		modulesNode, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleNode, err := modulesNode.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}
		listsNode, err := moduleNode.GetChild("lists")
		if err != nil {
			return nil, err
		}
		listNode, err := listsNode.GetChild(input.ListId)
		if err != nil {
			return nil, err
		}

		resp := &inputStructs.MListResp{}
		if err = mapstructure.Decode(listNode.Get(), &resp.Body); err != nil {
			return nil, err
		}

		elems := make([]map[string]any, 0, listNode.ChildsAmount())
		for _, elemNode := range listNode.GetChilds() {
			elems = append(elems, elemNode.Get())
		}

		resp.Body.Elems = elems

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "post-module",
		Method:        http.MethodPost,
		Path:          "/modules",
		Summary:       "Create new Module",
		Description:   "Create new Module",
		Tags:          []string{"Modules"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *inputStructs.ModulePost) (*inputStructs.ModuleResp, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}

		moduleN := modulesN.CreateChild(input.Body.Name)
		if err = moduleN.Init(); err != nil {
			return nil, err
		}
		var newModuleValue map[string]any
		if err = mapstructure.Decode(input.Body, &newModuleValue); err != nil {
			return nil, err
		}
		if err = moduleN.Put(newModuleValue); err != nil {
			return nil, err
		}
		listsN := moduleN.CreateChild("lists")
		if err = listsN.Init(); err != nil {
			return nil, err
		}

		resp := &inputStructs.ModuleResp{}
		if err = mapstructure.Decode(moduleN.Get(), &resp.Body); err != nil {
			return nil, err
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "post-module-list",
		Method:        http.MethodPost,
		Path:          "/modules/{moduleName}",
		Summary:       "Create new Module list",
		Description:   "Create new Module list",
		Tags:          []string{"Modules/Lists"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *inputStructs.MListPost) (*inputStructs.MListResp, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}
		listsN, err := moduleN.GetChild("lists")
		if err != nil {
			return nil, err
		}

		listId := uuid.New().String()
		listN := listsN.CreateChild(listId)
		if err = listN.Init(); err != nil {
			return nil, err
		}

		var listData map[string]any
		if err = mapstructure.Decode(input.Body, &listData); err != nil {
			return nil, err
		}

		listData["id"] = listId
		if err = listN.Put(listData); err != nil {
			return nil, err
		}

		resp := &inputStructs.MListResp{}
		if err = mapstructure.Decode(listData, &resp.Body); err != nil {
			return nil, err
		}

		resp.Body.Elems = make([]inputStructs.BaseMListElem, 0)
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "patch-module",
		Method:        http.MethodPatch,
		Path:          "/modules/{moduleName}",
		Summary:       "Patch Module",
		Description:   "Change Module contents",
		Tags:          []string{"Modules"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.ModulePatch) (*inputStructs.ModuleResp, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}

		var newModData map[string]any
		if err = mapstructure.Decode(input.Body, &newModData); err != nil {
			return nil, err
		}
		if err = moduleN.Put(utils.MergeMap(moduleN.Get(), newModData)); err != nil {
			return nil, err
		}

		resp := &inputStructs.ModuleResp{}
		if err = mapstructure.Decode(moduleN.Get(), &resp.Body); err != nil {
			return nil, err
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "patch-module-list",
		Method:        http.MethodPatch,
		Path:          "/modules/{moduleName}/{listId}",
		Summary:       "Patch Module list",
		Description:   "Change Module list contents",
		Tags:          []string{"Modules/Lists"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *inputStructs.MListPatch) (*struct{}, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}
		listsN, err := moduleN.GetChild("lists")
		if err != nil {
			return nil, err
		}
		listN, err := listsN.GetChild(input.ListId)
		if err != nil {
			return nil, err
		}

		var newListData map[string]any
		if err = mapstructure.Decode(input.Body, &newListData); err != nil {
			return nil, err
		}
		if err = listN.Put(utils.MergeMap(listN.Get(), newListData)); err != nil {
			return nil, err
		}

		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "patch-module-add-to-list",
		Method:        http.MethodPost,
		Path:          "/modules/{moduleName}/{listId}",
		Summary:       "Add content to Module list",
		Description:   "Adds new records to Module list",
		Tags:          []string{"Modules/Lists/Elements"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *inputStructs.MListElemsPatch) (*inputStructs.ListElemsResp, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}
		listsN, err := moduleN.GetChild("lists")
		if err != nil {
			return nil, err
		}
		listN, err := listsN.GetChild(input.ListId)
		if err != nil {
			return nil, err
		}

		resp := &inputStructs.ListElemsResp{}
		for _, elem := range input.Body {
			var elemMap map[string]any
			if err = mapstructure.Decode(elem, &elemMap); err != nil {
				return nil, err
			}

			elemMap["id"] = uuid.New().String()
			elemNode := listN.CreateChild(elemMap["id"].(string))

			if err = elemNode.Init(); err != nil {
				return nil, err
			}
			if err = elemNode.Put(elemMap); err != nil {
				return nil, err
			}

			resp.Body.List = append(resp.Body.List, elemMap)
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "patch-module-patch-list-element",
		Method:        http.MethodPatch,
		Path:          "/modules/{moduleName}/{listId}/{elementId}",
		Summary:       "Patch list element",
		Description:   "Patches specified record in Module list",
		Tags:          []string{"Modules/Lists/Elements"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *inputStructs.MListElemPatch) (*struct{}, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}
		listsN, err := moduleN.GetChild("lists")
		if err != nil {
			return nil, err
		}
		listN, err := listsN.GetChild(input.ListId)
		if err != nil {
			return nil, err
		}
		elem, err := listN.GetChild(input.ElemId)
		if err != nil {
			return nil, err
		}

		var newElemData map[string]any
		if err = mapstructure.Decode(input.Body, &newElemData); err != nil {
			return nil, err
		}
		if err = elem.Put(utils.MergeMap(elem.Get(), newElemData)); err != nil {
			return nil, err
		}

		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "patch-module-remove-from-list",
		Method:        http.MethodDelete,
		Path:          "/modules/{moduleName}/{listId}/{elementId}",
		Summary:       "Remove content from Module list",
		Description:   "Removes specified records from Module list",
		Tags:          []string{"Modules/Lists/Elements"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *inputStructs.ModuleMListElemPathId) (*struct{}, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}
		listsN, err := moduleN.GetChild("lists")
		if err != nil {
			return nil, err
		}
		listN, err := listsN.GetChild(input.ListId)
		if err != nil {
			return nil, err
		}

		err = listN.RemoveChild(input.ElemId)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "patch-module-replace-list",
		Method:        http.MethodPatch,
		Path:          "/modules/{moduleName}/{listId}/replace",
		Summary:       "Replace content for Module list",
		Description:   "Replaces all elements of specified Module list",
		Tags:          []string{"Modules/Lists/Elements"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *inputStructs.MListElemsPatch) (*inputStructs.ListElemsResp, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}
		listsN, err := moduleN.GetChild("lists")
		if err != nil {
			return nil, err
		}
		listN, err := listsN.GetChild(input.ListId)
		if err != nil {
			return nil, err
		}

		for _, elemN := range listN.GetChilds() {
			err := listN.RemoveChild(elemN.Get()["id"].(string))
			if err != nil {
				return nil, err
			}
		}

		resp := &inputStructs.ListElemsResp{}
		for _, elem := range input.Body {
			elemMap := map[string]any{}
			if err = mapstructure.Decode(elem, &elemMap); err != nil {
				return nil, err
			}

			elemMap["id"] = uuid.New().String()
			objN := listN.CreateChild(elemMap["id"].(string))

			if err = objN.Init(); err != nil {
				return nil, err
			}
			if err = objN.Put(elemMap); err != nil {
				return nil, err
			}

			resp.Body.List = append(resp.Body.List, elemMap)
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-module",
		Method:        http.MethodDelete,
		Path:          "/modules/{moduleName}",
		Summary:       "Delete Module",
		Description:   "Delete Module contents by id",
		Tags:          []string{"Modules"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *inputStructs.ModulePathName) (*struct{}, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}

		if err = modulesN.RemoveChild(input.ModuleName); err != nil {
			return nil, err
		}

		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-module-list",
		Method:        http.MethodDelete,
		Path:          "/modules/{moduleName}/{listId}",
		Summary:       "Delete Module list",
		Description:   "Delete Module list contents by id",
		Tags:          []string{"Modules/Lists"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *inputStructs.ModuleMListPathId) (*struct{}, error) {
		modulesN, err := root.GetChild("modules")
		if err != nil {
			return nil, err
		}
		moduleN, err := modulesN.GetChild(input.ModuleName)
		if err != nil {
			return nil, err
		}
		listsN, err := moduleN.GetChild("lists")
		if err != nil {
			return nil, err
		}

		if err = listsN.RemoveChild(input.ListId); err != nil {
			return nil, err
		}

		return nil, nil
	})
}
