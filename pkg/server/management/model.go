package management

type CopyBlueprintFromFileReq struct {
	BlueprintName string `json:"blueprintName"`
	FileName      string `json:"fileName"`
}

type CopyBlueprintReq struct {
	Token   string `json:"token"`
	NewName string `json:"newName"`
}

type BlueprintMethodReq struct {
	BlueprintName string `json:"blueprintName"`
	FileName      string `json:"filename"`
	MethodName    string `json:"methodName"`
}

type ListAppendBlueprintReq struct {
	BlueprintName string   `json:"blueprintName"`
	FileName      string   `json:"filename"`
	MethodName    []string `json:"methodName"`
}
