package middleware

type GoMiddleware struct{}

func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
