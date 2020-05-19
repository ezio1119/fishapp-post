package infrastructure

import (
	"log"

	"github.com/ezio1119/fishapp-post/conf"
	"github.com/ezio1119/fishapp-post/interfaces/controllers"
	"github.com/ezio1119/fishapp-post/interfaces/controllers/event"
	stan "github.com/nats-io/stan.go"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewNatsStreamingConn() (stan.Conn, error) {
	return stan.Connect(conf.C.Nats.ClusterID, conf.C.Nats.ClientID, stan.NatsURL(conf.C.Nats.URL))
}

func StartSubscribeCreatePostSagaReply(conn stan.Conn, c controllers.SagaReplyController) {
	conn.QueueSubscribe("create.post.saga.reply", conf.C.Nats.ClientID, func(m *stan.Msg) {
		e := &event.Event{}
		if err := protojson.Unmarshal(m.MsgProto.Data, e); err != nil {
			log.Fatal(err)
		}

		switch e.EventType {
		case "room.created":
			event.
			protojson.Unmarshal(e.EventData, m proto.Message)
			e.EventData
			c.RoomCreated(ctx, sagaID string)
		}

	}, stan.DurableName(conf.C.Nats.ClientID))
}
