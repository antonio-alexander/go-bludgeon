package queue

import (
	data "github.com/antonio-alexander/go-bludgeon/changes/data"

	goqueue "github.com/antonio-alexander/go-queue"
)

func ChangePartialConvertSingle(item interface{}) data.ChangePartial {
	switch v := item.(type) {
	default:
		return data.ChangePartial{}
	case data.ChangePartial:
		return v
	case []byte:
		changePartial := new(data.ChangePartial)
		if err := changePartial.UnmarshalBinary(v); err != nil {
			return data.ChangePartial{}
		}
		return *changePartial
	}
}

func ChangePartialConvertMultiple(items []interface{}) []data.ChangePartial {
	values := make([]data.ChangePartial, 0, len(items))
	for _, item := range items {
		values = append(values, ChangePartialConvertSingle(item))
	}
	return values
}

func ChangePartialPeek(queue goqueue.Peeker) (changePartials []data.ChangePartial) {
	return ChangePartialConvertMultiple(queue.Peek())
}

func ChangePartialFlush(queue goqueue.Dequeuer) (changePartials []data.ChangePartial) {
	return ChangePartialConvertMultiple(queue.Flush())
}

func ChangePartialEnqueueMultiple(queue goqueue.Enqueuer, changePartials []data.ChangePartial) (changePartialsRemaining []data.ChangePartial, overflow bool) {
	var items []interface{}

	for _, changePartial := range changePartials {
		items = append(items, changePartial)
	}
	itemsRemaining, overflow := queue.EnqueueMultiple(items)
	if overflow {
		return nil, true
	}
	return ChangePartialConvertMultiple(itemsRemaining), false
}
