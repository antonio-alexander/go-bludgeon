package tests

import (
	"context"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/logic"

	internal_kafka "github.com/antonio-alexander/go-bludgeon/pkg/kafka"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type Fixture struct {
	logic.Logic
	internal_kafka.Client
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateId() string {
	return uuid.Must(uuid.NewRandom()).String()
}

func NewFixture(items ...interface{}) *Fixture {
	f := &Fixture{}
	for _, item := range items {
		switch item := item.(type) {
		case logic.Logic:
			f.Logic = item
		case internal_kafka.Client:
			f.Client = item
		}
	}
	return f
}

func (f *Fixture) TestChangeHandler(changeTopic string) func(t *testing.T) {
	return func(t *testing.T) {
		var changeId string

		//generate dynamic constants
		ctx := context.TODO()
		changeReceived := make(chan struct{})
		serviceName := generateId()

		//subscribe for change
		handlerId, err := f.Subscribe(changeTopic, func(topic string, bytes []byte) {
			if len(bytes) == 0 {
				t.Logf("no bytes received")
				return
			}
			wrapper := &data.Wrapper{}
			if err := json.Unmarshal(bytes, wrapper); err != nil {
				t.Logf("error while unmarshalling json: %s", err)
				return
			}
			item, err := data.FromWrapper(wrapper)
			if err != nil {
				t.Logf("error during  FromWrapper: %s", err)
				return
			}
			switch v := item.(type) {
			case *data.Change:
				if changeId == v.Id {
					select {
					default:
						close(changeReceived)
						changeId = v.Id
						return
					case <-changeReceived:
					}
				}
			case *data.ChangeDigest:
				for _, change := range v.Changes {
					select {
					default:
						close(changeReceived)
						changeId = change.Id
						return
					case <-changeReceived:
					}
				}
			}
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, handlerId)

		//upsert change and validate change received
		dataId, version := generateId(), rand.Intn(1000)
		dataType, whenChanged := "employee", time.Now().UnixNano()
		changeCreated, err := f.ChangeUpsert(ctx, data.ChangePartial{
			DataId:          &dataId,
			DataVersion:     &version,
			DataType:        &dataType,
			DataServiceName: &serviceName,
			WhenChanged:     &whenChanged,
		})
		assert.Nil(t, err)
		changeId = changeCreated.Id
		defer func() {
			f.ChangesDelete(ctx, changeId)
		}()

		//validate that change received
		select {
		case <-time.After(10 * time.Second):
			assert.Fail(t, "unable to confirm change received")
		case <-changeReceived:
		}
	}
}
