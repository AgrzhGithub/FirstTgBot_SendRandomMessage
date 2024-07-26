package Telegram

import (
	"TelegramBot/clients/Telegram"
	"TelegramBot/events"
	err2 "TelegramBot/lib/err"
	"TelegramBot/storage"
	"errors"
)

type Processor struct {
	tg      *Telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUknownMetaType  = errors.New("Uknown meta type")
	ErrUknownEventType = errors.New("UknownEventType")
)

func New(client *Telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {

	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, err2.Wrap("can't get events", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))

	}
	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return err2.Wrap("can'r process message", ErrUknownEventType)

	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return err2.Wrap("can't process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return err2.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, err2.Wrap("can't get meta from event", ErrUknownMetaType)
	}
	return res, nil
}

func event(upd Telegram.Update) events.Event {

	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchType(upd Telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(upd Telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
