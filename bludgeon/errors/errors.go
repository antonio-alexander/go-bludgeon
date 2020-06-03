package bludgeonerrors

import (
	"fmt"
	// errors "github.com/pkg/errors"
)

type blunder struct {
	remote error
	meta   error
}

func (b *blunder) IsRemote(err error) bool { return b.remote != nil }
func (b *blunder) IsMeta(err error) bool   { return b.meta != nil }
func (b *blunder) Remote(err error) error  { return b.remote }
func (b *blunder) Meta(err error) error    { return b.meta }
func (b *blunder) Error() (err string) {
	//check if remote error
	if b.remote != nil {
		err = err + "Remote: " + b.remote.Error()
	}
	//check if meta error
	if b.meta != nil {
		err = err + "Remote: " + b.meta.Error()
	}

	return
}

func Meta(err error) error {
	type meta interface {
		Meta(err error) error
	}

	if err != nil {
		if meta, ok := err.(meta); ok {
			return meta.Meta(err)
		}
	}

	return err
}

func Metaf(err error, format string, args ...interface{}) error {
	type meta interface {
		Meta(err error) error
	}

	if err != nil {
		if meta, ok := err.(meta); ok {
			return meta.Meta(fmt.Errorf(format, args...))
		}
	} else {
		return &blunder{
			meta: fmt.Errorf(format, args...),
		}
	}

	return err
}

func Remote(err error) error {
	type remote interface {
		Remote(err error) error
	}

	if err != nil {
		if remote, ok := err.(remote); ok {
			return remote.Remote(err)
		}
	}

	return err
}

func Remotef(err error, format string, args ...interface{}) error {
	type remote interface {
		Remote(err error) error
	}

	if err != nil {
		if remote, ok := err.(remote); ok {
			return remote.Remote(fmt.Errorf(format, args...))
		}
	} else {
		return &blunder{
			remote: fmt.Errorf(format, args...),
		}
	}

	return err
}
