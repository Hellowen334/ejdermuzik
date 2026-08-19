/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package handlers

import (
	"ashokshau/tgmusic/src/core"
	"ashokshau/tgmusic/src/core/cache"
	"fmt"
	"time"

	td "github.com/AshokShau/gotdbot"
)

func handleVoiceChatMessage(c *td.Client, update *td.UpdateNewMessage) error {
	m := update.Message
	chatID := m.ChatId

	if m.IsGroup() {
		text := fmt.Sprintf(
			"<blockquote>⚠️ <b>Bu grup (%d) henüz bir süpergrup değil.</b>\nLütfen bu grubu bir süpergruba dönüştürün ve beni yönetici yapın.</blockquote>",
			chatID,
		)

		_, _ = c.SendTextMessage(chatID, text, &td.SendTextMessageOpts{
			ReplyMarkup:           core.AddMeMarkup(c.Me.Usernames.EditableUsername),
			DisableWebPagePreview: true,
			ParseMode:             "HTML",
		})

		time.Sleep(1 * time.Second)
		_ = c.LeaveChat(chatID)
		return nil
	}

	go storeChatToDB(chatID)

	if m.Content == nil {
		return nil
	}
	var message string
	switch m.Content.(type) {
	case *td.MessageVideoChatStarted:
		cache.ChatCache.ClearChat(chatID)
		message = "<blockquote>🎙️ <b>Sesli sohbet başlatıldı!</b>\nMüzik dinlemek için <code>/play [şarkı adı]</code> yazın.</blockquote>"
	case *td.MessageVideoChatEnded:
		cache.ChatCache.ClearChat(chatID)
		message = "<blockquote>🎧 <b>Sesli sohbet sona erdi!</b>\nTüm sıralar temizlendi.</blockquote>"
	default:
		return nil
	}

	_, _ = c.SendTextMessage(chatID, message, &td.SendTextMessageOpts{ParseMode: "HTML"})
	return td.EndGroups
}
