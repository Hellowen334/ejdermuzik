/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package handlers

import (
	"ashokshau/tgmusic/src/core/cache"
	"fmt"

	td "github.com/AshokShau/gotdbot"
)

func autoplayHandler(c *td.Client, m *td.Message) error {
	if !adminMode(c, m) {
		return td.EndGroups
	}

	chatID := m.ChatId

	if cache.ChatCache.GetPlayingTrack(chatID) == nil {
		_, err := m.ReplyText(c, "<blockquote>⚠️ Sesli sohbette çalan bir yayın yok.</blockquote>", &td.SendTextMessageOpts{ParseMode: "HTML"})
		return err
	}

	state := cache.ChatCache.GetAutoplay(chatID)
	text := "<blockquote>❤️‍🔥 <b>Otomatik Çalma Yönetimi</b>\n\nOtomatik çalma aktif olduğunda, sıra bittiğinde bot YouTube önerilerinden yeni parçaları otomatik olarak oynatır.</blockquote>"
	button := autoplayButton(state)

	_, err := m.ReplyText(c, text, &td.SendTextMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: button,
	})
	return err
}

func autoplayCallbackHandler(c *td.Client, cb *td.UpdateNewCallbackQuery) error {
	if !adminModeCB(c, cb) {
		return nil
	}

	chatID := cb.ChatId
	if cache.ChatCache.GetPlayingTrack(chatID) == nil {
		_ = cb.Answer(c, 0, true, "Sesli sohbette çalan bir yayın yok.", "")
		return nil
	}

	state := cache.ChatCache.GetAutoplay(chatID)
	newState := !state
	cache.ChatCache.SetAutoplay(chatID, newState)

	text := "<blockquote>❤️‍🔥 <b>Otomatik Çalma Yönetimi</b>\n\nOtomatik çalma aktif olduğunda, sıra bittiğinde bot YouTube önerilerinden yeni parçaları otomatik olarak oynatır.</blockquote>"
	button := autoplayButton(newState)

	_, err := cb.EditMessageText(c, text, &td.EditTextMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: button,
	})
	if err != nil {
		c.Logger.Warn("Failed to edit autoplay message", "error", err)
	}

	var status string
	if newState {
		status = "açıldı ✅"
	} else {
		status = "kapatıldı ❌"
	}
	_ = cb.Answer(c, 0, false, fmt.Sprintf("Otomatik Çalma %s.", status), "")

	return nil
}

func autoplayButton(state bool) *td.ReplyMarkupInlineKeyboard {
	var text string
	if state {
		text = "❤️‍🔥 Otomatik Çalma: AÇIK | ✅"
	} else {
		text = "❤️‍🔥 Otomatik Çalma: KAPALI | ❌"
	}

	return &td.ReplyMarkupInlineKeyboard{
		Rows: [][]td.InlineKeyboardButton{
			{
				{
					Text: text,
					Type: &td.InlineKeyboardButtonTypeCallback{
						Data: []byte("autoplay_toggle"),
					},
					Style: td.ButtonStylePrimary{},
				},
			},
		},
	}
}
