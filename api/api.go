package api

type Api interface {
	StartServer()
	GracefulStopServer()
}
