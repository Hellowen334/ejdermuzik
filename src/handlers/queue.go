/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package handlers

import (
	"ashokshau/tgmusic/src/utils"
	"fmt"
	"math"
	"strconv"
	"strings"

	"ashokshau/tgmusic/src/core/cache"
	"ashokshau/tgmusic/src/vc"

	td "github.com/AshokShau/gotdbot"
)

// queueHandler displays the current playback queue with detailed information.
func queueHandler(c *td.Client, m *td.Message) error {
	if !adminMode(c, m) {
		return td.EndGroups
	}

	chatID := m.ChatId

	chat, err := c.GetChat(chatID)
	if err != nil {
		_, _ = m.ReplyText(c, "<blockquote>⚠️ Sohbet bilgisi alınamadı.</blockquote>", &td.SendTextMessageOpts{ParseMode: "HTML"})
		return nil
	}

	queue := cache.ChatCache.GetQueue(chatID)
	if len(queue) == 0 {
		_, _ = m.ReplyText(c, "<blockquote>🎵 <b>Oynatma sırası boş.</b></blockquote>", &td.SendTextMessageOpts{ParseMode: "HTML"})
		return nil
	}

	if !cache.ChatCache.IsActive(chatID) {
		_, _ = m.ReplyText(c, "<blockquote>⚠️ Bot sesli sohbette yayın yapmıyor.</blockquote>", &td.SendTextMessageOpts{ParseMode: "HTML"})
		return nil
	}

	current := queue[0]
	playedTime, _ := vc.Calls.PlayedTime(chatID)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("<blockquote>📜 <b>%s Oynatma Sırası</b>\n\n", chat.Title))

	b.WriteString("▶️ <b>Şimdi Çalıyor:</b>\n")
	b.WriteString(fmt.Sprintf("• <b>Şarkı:</b> <code>%s</code>\n", truncate(current.Name, 45)))
	b.WriteString(fmt.Sprintf("• <b>İsteyen:</b> %s\n", current.User))
	b.WriteString(fmt.Sprintf("• <b>Süre:</b> %s dk\n", utils.SecToMin(current.Duration)))
	b.WriteString("• <b>Döngü:</b> ")
	if current.Loop > 0 {
		b.WriteString("Açık ✅\n")
	} else {
		b.WriteString("Kapalı ❌\n")
	}
	b.WriteString("• <b>İlerleme:</b> ")
	if playedTime > 0 && playedTime < math.MaxInt {
		b.WriteString(utils.SecToMin(int(playedTime)))
	} else {
		b.WriteString("0:00")
	}
	b.WriteString(" dk\n")

	if len(queue) > 1 {
		b.WriteString(fmt.Sprintf("\n⏭️ <b>Sıradakiler (%d):</b>\n", len(queue)-1))

		for i, song := range queue[1:] {
			if i >= 14 {
				break
			}
			b.WriteString(strconv.Itoa(i + 1))
			b.WriteString(". <code>")
			b.WriteString(truncate(song.Name, 45))
			b.WriteString("</code> | ")
			b.WriteString(utils.SecToMin(song.Duration))
			b.WriteString(" dk\n")
		}

		if len(queue) > 15 {
			b.WriteString(fmt.Sprintf("...ve %d parça daha\n", len(queue)-15))
		}
	}

	b.WriteString(fmt.Sprintf("\n📊 <b>Toplam:</b> %d parça</blockquote>", len(queue)))

	text := b.String()
	if len(text) > 4096 {
		var sb strings.Builder
		progress := "0:00"
		if playedTime > 0 && playedTime < math.MaxInt {
			progress = utils.SecToMin(int(playedTime))
		}
		sb.WriteString(fmt.Sprintf(
			"<blockquote>📜 <b>%s Oynatma Sırası</b>\n\n▶️ <b>Şimdi Çalıyor:</b>\n• <code>%s</code>\n• %s/%s dk\n\n📊 <b>Toplam:</b> %d parça</blockquote>",
			chat.Title,
			truncate(current.Name, 45),
			progress,
			utils.SecToMin(current.Duration),
			len(queue),
		))
		text = sb.String()
	}

	_, err = m.ReplyText(c, text, replyOpts)
	return err
}
