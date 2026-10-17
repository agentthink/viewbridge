package sip

type RouteGroup struct {
	method      string
	middlewares []HandlerFunc
	s           *Server
}

type MessageReceive struct {
	CmdType string `xml:"CmdType"`
	SN      int    `xml:"SN"`
}

func newRouteGroup(method string, s *Server, ms ...HandlerFunc) *RouteGroup {
	return &RouteGroup{
		method:      method,
		middlewares: ms,
		s:           s,
	}
}

// Use 注册组级中间件，仅对该分组下的路由生效
func (g *RouteGroup) Use(middleware ...HandlerFunc) *RouteGroup {
	g.middlewares = append(g.middlewares, middleware...)
	return g
}

func (g *RouteGroup) addGroup(pattern string, handler ...HandlerFunc) {
	key := g.method + "-" + pattern
	g.s.addRoute(key, append(g.middlewares, handler...)...)
}

func (g *RouteGroup) Handle(pattern string, handler ...HandlerFunc) {
	g.addGroup(pattern, handler...)
}
