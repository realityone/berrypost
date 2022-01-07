package management

type Common struct {
	Annotation map[string]string `json:"annotation"`
}

const (
	AppBerrypostManagementInvokeDefaultTarget = "app.berrypost.management.invoke.default.target"
	AppBerrypostManagementInvokePreferTarget  = "app.berrypost.management.invoke.prefer.target"
)
