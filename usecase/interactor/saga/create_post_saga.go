package saga

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ezio1119/fishapp-post/models"
	"github.com/ezio1119/fishapp-post/pb"
	"github.com/ezio1119/fishapp-post/usecase/repo"
	"github.com/google/uuid"
	"github.com/looplab/fsm"
	"google.golang.org/protobuf/encoding/protojson"
)

type createPostSagaState struct {
	sagaID       string
	sagaType     string
	currentState string
	post         *models.Post
	createdAt    time.Time
	updatedAt    time.Time
}

func NewCreatePostSagaState(p *models.Post, state, sagaID string) *createPostSagaState {
	now := time.Now()
	return &createPostSagaState{
		sagaType:     "CreatePostSaga",
		post:         p,
		sagaID:       sagaID,
		currentState: state,
		createdAt:    now,
		updatedAt:    now,
	}
}

func (s *createPostSagaState) convSagaInstance(state string) (*models.SagaInstance, error) {
	jsonPost, err := json.Marshal(s.post)
	if err != nil {
		return nil, err
	}
	return &models.SagaInstance{
		ID:           s.sagaID,
		SagaType:     s.sagaType,
		SagaData:     jsonPost,
		CurrentState: state,
		CreatedAt:    s.createdAt,
		UpdatedAt:    s.updatedAt,
	}, nil
}

type CreatePostSagaManager struct {
	FSM                 *fsm.FSM
	createPostSagaState *createPostSagaState
	outboxRepo          repo.OutboxRepo
	postRepo            repo.PostRepo
	sagaInstanceRepo    repo.SagaInstanceRepo
	transactionRepo     repo.TransactionRepo
	imageUploaderRepo   repo.ImageUploaderRepo
}

func InitCreatePostSagaManager(
	or repo.OutboxRepo,
	pr repo.PostRepo,
	sr repo.SagaInstanceRepo,
	tr repo.TransactionRepo,
	iur repo.ImageUploaderRepo,
) *CreatePostSagaManager {
	return &CreatePostSagaManager{
		outboxRepo:        or,
		postRepo:          pr,
		sagaInstanceRepo:  sr,
		transactionRepo:   tr,
		imageUploaderRepo: iur,
	}
}

func (m *CreatePostSagaManager) NewCreatePostSagaManager(state *createPostSagaState) *CreatePostSagaManager {

	m.createPostSagaState = state

	m.FSM = fsm.NewFSM(
		"init",
		fsm.Events{
			// {Name: "UploadImage", Src: []string{"Init"}, Dst: "UploadingImage"},
			// {Name: "UploadImage", Src: []string{"Init"}, Dst: "UploadingImage"},
			// {Name: "", Src: []string{"UploadingImage"}, Dst: "CreatingCreateRoom"},
			{Name: "CreateRoom", Src: []string{"init"}, Dst: "CreatingRoom"},
			{Name: "RejectPost", Src: []string{"CreatingRoom"}, Dst: "PostRejected"},
			{Name: "ApprovePost", Src: []string{"CreatingRoom"}, Dst: "PostApproved"},
		},
		fsm.Callbacks{
			// "UploadImage": func(e *fsm.Event) { s.uploadImage(e) },
			"CreateRoom":  func(e *fsm.Event) { m.createRoom(e) },
			"RejectPost":  func(e *fsm.Event) { m.rejectPost(e) },
			"ApprovePost": func(e *fsm.Event) { m.approvePost(e) },
			// "enter_PostRejected": func(e *fsm.Event) { m.rejectPost(e) },
			// "enter_state": func(e *fsm.Event) { s.enterState(e) },
		},
	)

	m.FSM.SetState(m.createPostSagaState.currentState)

	return m
}

func (m *CreatePostSagaManager) createRoom(e *fsm.Event) {
	ctx, ok := e.Args[0].(context.Context)
	if !ok {
		e.Cancel(errors.New("missing context"))
	}

	// 遷移先のステートを入れる
	sagaIn, err := m.createPostSagaState.convSagaInstance(e.Dst)
	if err != nil {
		e.Cancel(err)
	}

	eventData, err := protojson.Marshal(&pb.CreateRoom{
		SagaId: m.createPostSagaState.sagaID,
		PostId: m.createPostSagaState.post.ID,
		UserId: m.createPostSagaState.post.UserID,
	})

	now := time.Now()

	event := &models.Outbox{
		ID:        uuid.New().String(),
		EventType: "create.room",
		EventData: eventData,
		Channel:   "create.room",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// sagaインスタンスの永続化とイベントの発行を同じトランザクション内でやる
	ctx, err = m.transactionRepo.BeginTx(ctx)
	if err != nil {
		e.Cancel(err)
	}

	defer func() {
		if recover() != nil {
			m.transactionRepo.Roolback(ctx)
		}
	}()

	if err := m.outboxRepo.CreateOutbox(ctx, event); err != nil {
		m.transactionRepo.Roolback(ctx)
		e.Cancel(err)
		return
	}

	if err := m.sagaInstanceRepo.CreateSagaInstance(ctx, sagaIn); err != nil {
		m.transactionRepo.Roolback(ctx)
		e.Cancel(err)
		return
	}

	ctx, err = m.transactionRepo.Commit(ctx)
	if err != nil {
		e.Cancel(err)
		return
	}

	m.createPostSagaState.currentState = e.Dst
}

func (m *CreatePostSagaManager) rejectPost(e *fsm.Event) {
	ctx, ok := e.Args[0].(context.Context)
	if !ok {
		e.Cancel(errors.New("missing context"))
	}
	jsonPost, err := json.Marshal(m.createPostSagaState.post)
	if err != nil {
		e.Cancel(err)
	}

	now := time.Now()
	event := &models.Outbox{
		ID:        uuid.New().String(),
		EventType: "post.rejected",
		EventData: jsonPost,
		Channel:   "post.rejected",
		CreatedAt: now,
		UpdatedAt: now,
	}

	sagaIn := &models.SagaInstance{
		ID:           m.createPostSagaState.sagaID,
		SagaType:     m.createPostSagaState.sagaType,
		SagaData:     jsonPost,
		CurrentState: e.Dst,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// createPostSagaFailedとサガイベントも発行する
	ctx, err = m.transactionRepo.BeginTx(ctx)
	if err != nil {
		e.Cancel(err)
	}

	defer func() {
		if recover() != nil {
			m.transactionRepo.Roolback(ctx)
		}
	}()

	if err := m.postRepo.DeletePost(ctx, m.createPostSagaState.post.ID); err != nil {
		m.transactionRepo.Roolback(ctx)
		e.Cancel(err)
		return
	}

	if err := m.outboxRepo.CreateOutbox(ctx, event); err != nil {
		m.transactionRepo.Roolback(ctx)
		e.Cancel(err)
		return
	}

	if err := m.sagaInstanceRepo.UpdateSagaInstance(ctx, sagaIn); err != nil {
		m.transactionRepo.Roolback(ctx)
		e.Cancel(err)
		return
	}

	ctx, err = m.transactionRepo.Commit(ctx)
	if err != nil {
		e.Cancel(err)
		return
	}

}

func (m *CreatePostSagaManager) approvePost(e *fsm.Event) {
	ctx, ok := e.Args[0].(context.Context)
	if !ok {
		e.Cancel(errors.New("missing context"))
	}
	jsonPost, err := json.Marshal(m.createPostSagaState.post)
	if err != nil {
		e.Cancel(err)
	}

	now := time.Now()
	event := &models.Outbox{
		ID:        uuid.New().String(),
		EventType: "post.approved",
		EventData: jsonPost,
		Channel:   "post.approved",
		CreatedAt: now,
		UpdatedAt: now,
	}

	sagaIn := &models.SagaInstance{
		ID:           m.createPostSagaState.sagaID,
		SagaType:     m.createPostSagaState.sagaType,
		SagaData:     jsonPost,
		CurrentState: e.Dst,
		UpdatedAt:    time.Now(),
	}

	ctx, err = m.transactionRepo.BeginTx(ctx)
	if err != nil {
		e.Cancel(err)
	}

	defer func() {
		if recover() != nil {
			m.transactionRepo.Roolback(ctx)
		}
	}()

	if err := m.outboxRepo.CreateOutbox(ctx, event); err != nil {
		m.transactionRepo.Roolback(ctx)
		e.Cancel(err)
		return
	}

	if err := m.sagaInstanceRepo.UpdateSagaInstance(ctx, sagaIn); err != nil {
		m.transactionRepo.Roolback(ctx)
		e.Cancel(err)
		return
	}

	ctx, err = m.transactionRepo.Commit(ctx)
	if err != nil {
		e.Cancel(err)
		return
	}

	m.createPostSagaState.currentState = e.Dst
}
