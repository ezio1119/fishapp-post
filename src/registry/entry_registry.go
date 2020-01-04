package registry

import (
	"github.com/ezio1119/fishapp-post/entry/controllers"
	"github.com/ezio1119/fishapp-post/entry/controllers/entry_post_grpc"
	"github.com/ezio1119/fishapp-post/entry/interactor"
	"github.com/ezio1119/fishapp-post/entry/presenter"
	"github.com/ezio1119/fishapp-post/entry/repository"
	_postRepo "github.com/ezio1119/fishapp-post/post/repository"
)

func (r *registry) NewEntryController() entry_post_grpc.EntryServiceServer {
	return controllers.NewEntryController(
		interactor.NewEntryInteractor(
			repository.NewEntryRepository(r.db),
			presenter.NewEntryPresenter(),
			_postRepo.NewPostRepository(r.db),
			r.ctxTimeout,
		))
}
