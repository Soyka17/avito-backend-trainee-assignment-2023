package router

type Router interface {
	Get(string, func(string) (code int, resp map[string]any))
	Post(string, map[string]any, func(map[string]any) (int, map[string]any))
	Run()
}
