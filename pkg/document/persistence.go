package document

type Persistence interface {
	Save(item Document) error
	FindAll() ([]Document, error)
	FindByID(itemId string) (Document, error)
	Delete(itemId string) error
}

type Manager struct {
	pers Persistence
}

func NewManager(pers Persistence) *Manager {
	return &Manager{pers: pers}
}

func (m *Manager) Create(name, description string) (Document, error) {
	doc, err := NewDocument(
		WithName(name),
		WithDescription(description),
	)
	if err != nil {
		return Document{}, err
	}
	if err := m.pers.Save(*doc); err != nil {
		return Document{}, err
	}
	return *doc, nil
}

func (m *Manager) Get(id string) (Document, error) {
	return m.pers.FindByID(id)
}

func (m *Manager) Delete(id string) error {
	return m.pers.Delete(id)
}

func (m *Manager) FindAll() ([]Document, error) {
	return m.pers.FindAll()
}
