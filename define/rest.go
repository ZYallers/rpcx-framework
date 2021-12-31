package define

type Restful map[string][]RestHandler
type RestHandler struct {
	Sort    int
	Signed  bool
	Logged  bool
	Path    string
	Version string
	Method  string
	Service IService
}
