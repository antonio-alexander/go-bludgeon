package bludgeon

type Remote interface {
	//
	TimerCreate() (timer Timer, err error)

	//
	TimerRead(id string) (timer Timer, err error)

	//
	TimerUpdate(timer Timer) (err error)

	//
	TimerDelete(id string) (err error)

	//
	TimeSliceCreate(id string) (timeSlice TimeSlice, err error)

	//
	TimeSliceRead(id string) (timeSlice TimeSlice, err error)

	//
	TimeSliceUpdate(timeSlice TimeSlice) (err error)

	//
	TimeSliceDelete(id string) (err error)
}
