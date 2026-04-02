package models

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

type CompanyNotificationSettings struct {
	WebhookType string          `json:"webhook_type"`
	EndpointUrl string          `json:"endpoint"`
	MetaKeys    json.RawMessage `json:"metakeys"`
}
type CompanyMember struct {
	UserID     string            `json:"user_id"`
	Permission CompanyPermission `json:"permission"`
}

type Company struct {
	ID string `json:"id"`

	DisplayName string `json:"display_name"`
	Holder      string `json:"holder"` // Immutable

	LimitNodes int  `json:"limit_nodes"`
	Enabled    bool `json:"enabled"`

	memberMu sync.RWMutex
	Members  map[string]*CompanyMember `json:"members,omitempty"`

	NotificationEndpoints []*CompanyNotificationSettings `json:"notification_endpoints"`

	notificationSettingsMu sync.RWMutex
	NotificationThresholds map[NotificationType]int `json:"notification_thresholds"`

	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

func NewCompany(displayName, holder string) (*Company, error) {
	if holder == "" {
		return nil, errors.New("holder is empty")
	}

	if len(holder) > 26 {
		return nil, errors.New("holder is too long, maybe malformed")
	}

	if displayName == "" {
		return nil, errors.New("display name is empty")
	}

	if len(displayName) > 255 {
		return nil, errors.New("display name is too long")
	}

	defaultNotificationThresholds := make(map[NotificationType]int)

	for _, notificationType := range NotificationTypes() {
		defaultNotificationThresholds[notificationType] = 75
	}

	return &Company{
		ID:                     ulid.Make().String(),
		DisplayName:            displayName,
		Holder:                 holder,
		NotificationEndpoints:  make([]*CompanyNotificationSettings, 0),
		NotificationThresholds: defaultNotificationThresholds,
		Enabled:                true,
		CreatedAt:              time.Now(),
	}, nil
}

func (c *Company) AddMember(userID string, permission int) {
	c.memberMu.Lock()
	defer c.memberMu.Unlock()

	c.Members[userID] = &CompanyMember{
		UserID:     userID,
		Permission: CompanyPermission(permission),
	}
}

func (c *Company) GetMember(userID string) *CompanyMember {
	c.memberMu.RLock()
	defer c.memberMu.RUnlock()

	return c.Members[userID]
}

func (c *Company) RemoveMember(userID string) {
	c.memberMu.Lock()
	defer c.memberMu.Unlock()

	delete(c.Members, userID)
}

func (c *Company) SetNotificationThreshold(notificationType NotificationType, threshold int) {
	c.notificationSettingsMu.Lock()
	defer c.notificationSettingsMu.Unlock()

	c.NotificationThresholds[notificationType] = threshold
}

func (c *Company) GetNotificationThreshold(notificationType NotificationType) int {
	c.notificationSettingsMu.RLock()
	defer c.notificationSettingsMu.RUnlock()

	return c.NotificationThresholds[notificationType]
}

func (c *Company) AddNotificationEndpoint(endpoint *CompanyNotificationSettings) {
	c.notificationSettingsMu.Lock()
	defer c.notificationSettingsMu.Unlock()

	c.NotificationEndpoints = append(c.NotificationEndpoints, endpoint)
}

func (c *Company) RemoveNotificationEndpoint(endpoint *CompanyNotificationSettings) {
	c.notificationSettingsMu.Lock()
	defer c.notificationSettingsMu.Unlock()

	for i, e := range c.NotificationEndpoints {
		if e.EndpointUrl == endpoint.EndpointUrl {
			c.NotificationEndpoints = append(c.NotificationEndpoints[:i], c.NotificationEndpoints[i+1:]...)
			break
		}
	}
}

func (c *Company) UpdateDisplayName(displayName string) {
	c.DisplayName = displayName
	c.LastUpdated = time.Now()
}

func (c *Company) UpdateLimitNodes(limitNodes int) {
	c.LimitNodes = limitNodes
	c.LastUpdated = time.Now()
}
