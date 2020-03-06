package registry

import (
	"github.com/ezio1119/fishapp-post/post/controllers"
	"github.com/ezio1119/fishapp-post/post/controllers/post_grpc"
	"github.com/ezio1119/fishapp-post/post/interactor"
	"github.com/ezio1119/fishapp-post/post/presenter"
	"github.com/ezio1119/fishapp-post/post/repository"
)

func (r *registry) NewPostController() post_grpc.PostServiceServer {
	return controllers.NewPostController(
		interactor.NewPostInteractor(
			repository.NewPostRepository(r.db),
			presenter.NewPostPresenter(),
			r.ctxTimeout,
		))
}
