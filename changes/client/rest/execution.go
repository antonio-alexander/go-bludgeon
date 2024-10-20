package rest

import "github.com/antonio-alexander/go-bludgeon/changes/data"

func changePartialToChange(changePartial data.ChangePartial) *data.Change {
	change := new(data.Change)
	if changedBy := changePartial.ChangedBy; changedBy != nil {
		change.ChangedBy = *changedBy
	}
	if dataAction := changePartial.DataAction; dataAction != nil {
		change.DataAction = *dataAction
	}
	if dataId := changePartial.DataId; dataId != nil {
		change.DataId = *dataId
	}
	if dataServiceName := changePartial.DataServiceName; dataServiceName != nil {
		change.DataServiceName = *dataServiceName
	}
	if dataType := changePartial.DataType; dataType != nil {
		change.DataType = *dataType
	}
	if dataVersion := changePartial.DataVersion; dataVersion != nil {
		change.DataVersion = *dataVersion
	}
	if whenChanged := changePartial.WhenChanged; whenChanged != nil {
		change.WhenChanged = *whenChanged
	}
	return change
}
