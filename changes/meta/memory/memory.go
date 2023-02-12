package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/changes/meta"
	"github.com/antonio-alexander/go-bludgeon/internal"
	"github.com/antonio-alexander/go-bludgeon/internal/logger"

	"github.com/pkg/errors"
)

type memory struct {
	sync.RWMutex
	logger.Logger
	changesMux          sync.RWMutex
	changes             map[string]*data.Change
	registrationsMux    sync.RWMutex
	registrations       map[string]struct{}
	registrationChanges map[string]map[string]struct{}
}

func New() interface {
	meta.Serializer
	meta.Change
	meta.Registration
	meta.RegistrationChange
	internal.Initializer
	internal.Configurer
	internal.Parameterizer
} {
	return &memory{
		Logger:              logger.NewNullLogger(),
		changes:             make(map[string]*data.Change),
		registrations:       make(map[string]struct{}),
		registrationChanges: make(map[string]map[string]struct{}),
	}
}

func (m *memory) validateChange(c data.ChangePartial, create bool, ids ...string) error {
	var id string

	if len(ids) > 0 {
		id = strings.ToLower(ids[0])
	}
	if c.DataId == nil || *c.DataId == "" {
		return meta.ErrChangeNotWritten
	}
	for _, change := range m.changes {
		if change.Id == id {
			continue
		}
		if c.DataId != nil && *c.DataId == change.DataId &&
			c.DataServiceName != nil && *c.DataServiceName == change.DataServiceName &&
			c.DataType != nil && *c.DataType == change.DataType &&
			c.DataAction != nil && *c.DataAction == change.DataAction &&
			c.DataVersion != nil && *c.DataVersion == change.DataVersion {
			return meta.ErrChangeConflictWrite
		}
	}
	return nil
}

func (m *memory) validateRegistration(registrationId string) error {
	if registrationId == "" {
		return meta.ErrRegistrationNotWritten
	}
	return nil
}

func (m *memory) validateRegistrationChange(changeId string) error {
	if changeId == "" {
		return meta.ErrRegistrationChangeNotWritten
	}
	return nil
}

func (m *memory) findChangeIdsToDelete(changeIds ...string) ([]string, error) {
	var changeIdsToDelete []string

	for _, changeId := range changeIds {
		found := false
		for registrationId := range m.registrationChanges {
			if _, ok := m.registrationChanges[registrationId][changeId]; ok {
				found = true
				break
			}
		}
		if !found {
			changeIdsToDelete = append(changeIdsToDelete, changeId)
		}
	}
	return changeIdsToDelete, nil
}

func (m *memory) SetUtilities(parameters ...interface{}) {
	for _, p := range parameters {
		switch p := p.(type) {
		case logger.Logger:
			m.Logger = p
		}
	}
}

func (m *memory) SetParameters(parameters ...interface{}) {
	//
}

func (m *memory) Configure(...interface{}) error {
	//KIM: thsi is a dummy place holder to satisfy the
	// internal.Configurable interface
	return nil
}

func (m *memory) Initialize() error {
	//KIM: thsi is a dummy place holder to satisfy the
	// internal.Initializer interface
	return nil
}

func (m *memory) Shutdown() {
	m.Lock()
	defer m.Unlock()
	m.changes = make(map[string]*data.Change)
	m.Debug(logAlias + "shutdown")
}

func (m *memory) Serialize() (*meta.SerializedData, error) {
	m.Lock()
	defer m.Unlock()
	serializedData := &meta.SerializedData{
		Changes: make(map[string]data.Change),
	}
	for id, employee := range m.changes {
		serializedData.Changes[id] = *employee
	}
	return serializedData, nil
}

func (m *memory) Deserialize(serializedData *meta.SerializedData) error {
	m.Lock()
	defer m.Unlock()
	if serializedData == nil {
		return errors.New("serialized data is nil")
	}
	m.changes = make(map[string]*data.Change)
	for id, employee := range serializedData.Changes {
		m.changes[id] = &employee
	}
	return nil
}

func (m *memory) ChangeCreate(ctx context.Context, c data.ChangePartial) (*data.Change, error) {
	m.changesMux.Lock()
	defer m.changesMux.Unlock()
	if err := m.validateChange(c, true); err != nil {
		return nil, err
	}
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	change := &data.Change{
		Id: id,
	}
	if c.DataId != nil {
		change.DataId = *c.DataId
	}
	if c.DataVersion != nil {
		change.DataVersion = *c.DataVersion
	}
	if c.DataType != nil {
		change.DataType = *c.DataType
	}
	if c.DataServiceName != nil {
		change.DataServiceName = *c.DataServiceName
	}
	if c.DataAction != nil {
		change.DataAction = *c.DataAction
	}
	if c.WhenChanged != nil {
		change.WhenChanged = *c.WhenChanged
	}
	if c.ChangedBy != nil {
		change.ChangedBy = *c.ChangedBy
	}
	m.changes[id] = change
	m.Debug(logAlias+"created change: %s", change.Id)
	return copyChange(change), nil
}

func (m *memory) ChangeRead(ctx context.Context, changeId string) (*data.Change, error) {
	m.changesMux.RLock()
	defer m.changesMux.RUnlock()
	change, ok := m.changes[changeId]
	if !ok {
		return nil, meta.ErrChangeNotFound
	}
	m.Debug(logAlias+"read change: %s", changeId)
	return copyChange(change), nil
}

func (m *memory) ChangesDelete(ctx context.Context, changeIds ...string) error {
	m.changesMux.Lock()
	defer m.changesMux.Unlock()
	m.registrationsMux.Lock()
	defer m.registrationsMux.Unlock()

	for _, registrationChanges := range m.registrationChanges {
		for changeId := range registrationChanges {
			for _, changeIdToDelete := range changeIds {
				if changeIdToDelete == changeId {
					return meta.ErrChangeNotDeletedConflict
				}
			}
		}
	}
	for _, changeId := range changeIds {
		m.Debug(logAlias+"deleted change: %s", changeId)
		delete(m.changes, changeId)
	}
	return nil
}

func (m *memory) ChangesRead(ctx context.Context, search data.ChangeSearch) ([]*data.Change, error) {
	m.changesMux.RLock()
	defer m.changesMux.RUnlock()
	var versions map[string]int
	var latestVersion bool

	if search.LatestVersion != nil && *search.LatestVersion {
		latestVersion = true
		versions = make(map[string]int)
	}
	searchFx := func(c *data.Change) bool {
		//KIM: this is an inclusive search and is computationally expensive
		if len(search.ChangeIds) > 0 {
			found := false
			for _, id := range search.ChangeIds {
				if c.Id == id {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		if len(search.DataIds) > 0 {
			found := false
			for _, dataId := range search.DataIds {
				if c.DataId == dataId {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		if len(search.Types) > 0 {
			found := false
			for _, dataType := range search.Types {
				if c.DataType == dataType {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		if len(search.ServiceNames) > 0 {
			found := false
			for _, service := range search.ServiceNames {
				if c.DataServiceName == service {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		switch {
		case search.Since != nil:
			if c.WhenChanged < *search.Since {
				return false
			}
		}
		if latestVersion {
			if c.DataVersion > versions[c.Id] {
				versions[c.Id] = c.DataVersion
			}
		}
		return true
	}
	var changes []*data.Change
	for _, change := range m.changes {
		if searchFx(change) {
			switch {
			default:
				changes = append(changes, copyChange(change))
			case search.LatestVersion != nil && *search.LatestVersion &&
				change.DataVersion == versions[change.Id]:
				changes = append(changes, copyChange(change))
			}
		}
	}
	return changes, nil
}

func (m *memory) RegistrationUpsert(ctx context.Context, registrationId string) error {
	m.registrationsMux.Lock()
	defer m.registrationsMux.Unlock()

	if err := m.validateRegistration(registrationId); err != nil {
		return err
	}
	if _, ok := m.registrations[registrationId]; !ok {
		m.registrations[registrationId] = struct{}{}
	}
	if _, ok := m.registrationChanges[registrationId]; !ok {
		m.registrationChanges[registrationId] = make(map[string]struct{})
	}
	m.Debug(logAlias+"upserted registration: %s", registrationId)
	return nil
}

func (m *memory) RegistrationDelete(ctx context.Context, registrationId string) error {
	m.registrationsMux.Lock()
	defer m.registrationsMux.Unlock()
	_, ok := m.registrations[registrationId]
	if !ok {
		return meta.ErrRegistrationNotFound
	}
	delete(m.registrations, registrationId)
	delete(m.registrationChanges, registrationId)
	m.Debug(logAlias+"deleted registration: %s", registrationId)
	return nil
}

func (m *memory) RegistrationChangeUpsert(ctx context.Context, changeId string) error {
	m.registrationsMux.Lock()
	defer m.registrationsMux.Unlock()

	if err := m.validateRegistrationChange(changeId); err != nil {
		return err
	}
	if len(m.registrationChanges) == 0 {
		return nil
	}
	for registrationId := range m.registrationChanges {
		m.registrationChanges[registrationId][changeId] = struct{}{}
	}
	m.Debug(logAlias+"upserted registration change: %s", changeId)
	return nil
}

func (m *memory) RegistrationChangesRead(ctx context.Context, registrationId string) ([]string, error) {
	m.registrationsMux.RLock()
	defer m.registrationsMux.RUnlock()
	var changeIds []string

	_, ok := m.registrationChanges[registrationId]
	if !ok {
		return nil, meta.ErrRegistrationNotFound
	}
	for changeId := range m.registrationChanges[registrationId] {
		changeIds = append(changeIds, changeId)
	}
	return changeIds, nil
}

func (m *memory) RegistrationChangeAcknowledge(ctx context.Context, registrationId string, changeIds ...string) ([]string, error) {
	m.registrationsMux.Lock()
	defer m.registrationsMux.Unlock()
	if _, ok := m.registrations[registrationId]; !ok {
		return nil, meta.ErrRegistrationNotFound
	}
	for _, changeId := range changeIds {
		m.Debug(logAlias+"acknowledged change %s for %s", changeId, registrationId)
		delete(m.registrationChanges[registrationId], changeId)
	}
	return m.findChangeIdsToDelete(changeIds...)
}
