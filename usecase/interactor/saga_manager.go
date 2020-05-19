package interactor

import (
	"context"
	"encoding/json"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/usecase/repo"
	"github.com/google/uuid"
	"github.com/looplab/fsm"
)

// type sagaManager struct {
// 	OutboxRepo repo.OutboxRepo
// 	createPostSagaState
// }

type createPostSagaState struct {
	FSM        *fsm.FSM
	SagaID     string
	Post       *models.Post
	OutboxRepo repo.OutboxRepo
	// CmdPublisher CmdPublisher
}

// type cmdPublisher struct {
// 	OutboxRepo repo.OutboxRepo
// }

// type CmdPublisher interface {
// 	publishCreateRoomCmd(ctx context.Context, o *models.Outbox, i *models.SagaInstance) error
// }

// func (*cmdPublisher) publishCreateRoomCmd(ctx context.Context, o *models.Outbox, i *models.SagaInstance) error {}

func newCreatePostSagaState(p *models.Post, outboxRepo repo.OutboxRepo, sagaID string) *createPostSagaState {
	s := &createPostSagaState{
		SagaID:     sagaID,
		Post:       p,
		OutboxRepo: outboxRepo,
	}

	s.FSM = fsm.NewFSM(
		"init",
		fsm.Events{
			// {Name: "UploadImage", Src: []string{"Init"}, Dst: "UploadingImage"},
			// {Name: "UploadImage", Src: []string{"Init"}, Dst: "UploadingImage"},
			// {Name: "", Src: []string{"UploadingImage"}, Dst: "CreatingCreateRoom"},

			{Name: "CreateRoom", Src: []string{"init"}, Dst: "CreatingRoom"},
			{Name: "RejectPost", Src: []string{"CreatingRoom"}, Dst: "PostRejected"},
		},
		fsm.Callbacks{
			// "UploadImage": func(e *fsm.Event) { s.uploadImage(e) },
			"CreateRoom": func(e *fsm.Event) { s.createRoom(e) },
			// "enter_state": func(e *fsm.Event) { s.enterState(e) },
		},
	)

	return s
}

// func (s *createPostSagaState) uploadImage(e *fsm.Event) {
// 	fmt.Printf("uploadImage: %#v\n", e)
// }

func (s *createPostSagaState) createRoom(e *fsm.Event) {
	jsonPost, err := json.Marshal(s.Post)
	if err != nil {
		e.Cancel(err)
	}

	id, err := uuid.NewUUID()
	if err != nil {
		e.Cancel(err)
	}
	outbox := &models.Outbox{
		ID:        id.String(),
		EventType: "create.room",
		EventData: jsonPost,
	}

	sagaI := &models.SagaInstance{
		ID:           s.SagaID,
		SagaType:     "CreatePostSaga",
		SagaData:     jsonPost,
		CurrentState: "CreatingRoom",
	}
	if err := s.OutboxRepo.CreateOutbox(context.Background(), outbox, sagaI); err != nil {
		e.Cancel(err)
	}
}

// func (s *createPostSagaState) enterState(e *fsm.Event) {
// 	fmt.Printf("enterState: %#v\n", e)
// }
