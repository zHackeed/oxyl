package provider

import (
	"context"
	"sync"

	"zhacked.me/oxyl/service/notifications/internal/storage"
	comm "zhacked.me/oxyl/shared/pkg/models"
)

type NotificationSettingsProvider struct {
	mu       sync.RWMutex
	settings map[string][]*comm.CompanyNotificationSettings

	storage *storage.NotificationSettingStorage
}

func NewNotificationSettingsProvider(storage *storage.NotificationSettingStorage) *NotificationSettingsProvider {
	return &NotificationSettingsProvider{
		settings: make(map[string][]*comm.CompanyNotificationSettings),
		storage:  storage,
	}
}

func (p *NotificationSettingsProvider) Load(ctx context.Context) error {
	settings, err := p.storage.GetAll(ctx)
	if err != nil {
		return err
	}

	p.settings = settings
	return nil
}

func (p *NotificationSettingsProvider) Get(companyID string) ([]*comm.CompanyNotificationSettings, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	s, ok := p.settings[companyID]
	return s, ok
}

func (p *NotificationSettingsProvider) Add(ctx context.Context, companyID string) error {
	settings, err := p.storage.GetByCompany(ctx, companyID)
	if err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.settings[companyID] = settings
	return nil
}

func (p *NotificationSettingsProvider) New(ctx context.Context, id string) error {
	setting, err := p.storage.GetById(ctx, id)
	if err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.settings[setting.Holder] = append(p.settings[setting.Holder], setting)
	return nil
}

func (p *NotificationSettingsProvider) Remove(companyID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.settings, companyID)
}

func (p *NotificationSettingsProvider) RemoveSetting(companyID, settingID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	settings := p.settings[companyID]
	for i, s := range settings {
		if s.ID == settingID {
			p.settings[companyID] = append(settings[:i], settings[i+1:]...)
			return
		}
	}
}
