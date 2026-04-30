package notificator

import (
	"fmt"

	slack "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gtuk/discordwebhook"
	redisModels "zhacked.me/oxyl/shared/pkg/messenger/models"
	comm "zhacked.me/oxyl/shared/pkg/models"
)

func Send(setting *comm.CompanyNotificationSettings, event redisModels.ThresholdNotification, displayName string) error {
	switch setting.WebhookType {
	case comm.WebhookTypeSlack:
		return sendSlack(setting.Endpoint, *setting.Channel, event, displayName)
	case comm.WebhookTypeDiscord:
		return sendDiscord(setting.Endpoint, event, displayName)
	default:

		return fmt.Errorf("unknown value %v", setting.WebhookType)
	}
}

func sendDiscord(endpoint string, event redisModels.ThresholdNotification, displayName string) error {
	var title, colorHex string

	if event.Resolved {
		title = "Sobreconsumo resuelto"
		colorHex = "2238030"
	} else {
		title = "Sobreconsumo detectado"
		colorHex = "15623364"
	}

	fields := []discordwebhook.Field{
		{
			Name:  new("Agent"),
			Value: new(displayName),
		},
		{
			Name:  new("¿Por qué?"),
			Value: new(event.TriggerReason.Stringified()),
		},
	}

	if !event.Resolved {
		fields = append(fields, discordwebhook.Field{
			Name:  new("Valor"),
			Value: new(event.TriggerValue),
		})
	}

	embed := discordwebhook.Embed{
		Title:  new(title),
		Color:  new(colorHex),
		Fields: &fields,
	}

	msg := discordwebhook.Message{
		Embeds: &[]discordwebhook.Embed{embed},
	}

	return discordwebhook.SendMessage(endpoint, msg)
}

func sendSlack(endpoint, channel string, event redisModels.ThresholdNotification, displayName string) error {
	var title, color string

	if event.Resolved {
		title = "Sobreconsumo resuelto"
		color = "#22c55e"
	} else {
		title = "Sobre consumo detectado"
		color = "#ef4444"
	}

	fields := []*slack.Field{
		{
			Title: "Agent",
			Value: displayName,
			Short: true,
		},
		{
			Title: "Tipo",
			Value: event.TriggerReason.Stringified(),
			Short: true,
		},
	}

	if !event.Resolved {
		fields = append(fields, &slack.Field{
			Title: "Value",
			Value: event.TriggerValue,
			Short: true,
		})
	}

	attachment := slack.Attachment{
		Title:  &title,
		Color:  &color,
		Fields: fields,
	}

	payload := slack.Payload{
		Channel:     channel,
		Attachments: []slack.Attachment{attachment},
	}

	errs := slack.Send(endpoint, "", payload)
	if len(errs) > 0 {
		return fmt.Errorf("slack delivery failed: %w", errs[0])
	}

	return nil
}
