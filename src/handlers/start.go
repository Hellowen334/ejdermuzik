/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package handlers

import (
	"fmt"
	"runtime"
	"time"

	"ashokshau/tgmusic/src/core"
	"ashokshau/tgmusic/src/core/db"

	td "github.com/AshokShau/gotdbot"
)

// pingHandler handles the /ping command.
func pingHandler(c *td.Client, m *td.Message) error {

	start := time.Now()

	msg, err := m.ReplyText(c, "<blockquote>🔥 Ölçülüyor… lütfen bekleyin…</blockquote>", &td.SendTextMessageOpts{ParseMode: "HTML"})
	if err != nil {
		return err
	}

	latency := time.Since(start).Milliseconds()
	uptime := getFormattedDuration(time.Since(startTime))

	response := fmt.Sprintf(
		"<blockquote><b>📊 Sistem Performans Metrikleri</b></blockquote>\n\n"+
			"<blockquote expandable>⚡ <b>Bot Gecikmesi:</b> <code>%d ms</code>\n"+
			"⏳ <b>Çalışma Süresi:</b> <code>%s</code>\n"+
			"🔄 <b>Go Goroutine:</b> <code>%d</code></blockquote>",
		latency, uptime, runtime.NumGoroutine(),
	)

	_, err = msg.EditText(c, response, &td.EditTextMessageOpts{ParseMode: "HTML"})
	return err
}

// startHandler handles the /start command.
func startHandler(c *td.Client, m *td.Message) error {
	chatID := m.ChatId

	if m.IsPrivate() {
		go func(chatID int64) {
			_ = db.Instance.AddUser(chatID)
		}(chatID)

		text := fmt.Sprintf(
			"<blockquote>👀 <b>Merhaba %s, Hoş Geldin!</b> 🤗</blockquote>\n\n"+
				"<blockquote expandable>🫴 Ben <b>%s</b>\nEn hızlı ve en güçlü müzik + yayın botuyum!\n\n"+
				"💫 Grubunda ve kanalında kesintisiz, yüksek kalitede müzik dinlemek için beni grubuna ekleyebilirsin.</blockquote>\n\n"+
				"<blockquote>🔰 Tüm komutları görmek için <b>Yardım & Komutlar</b> butonuna tıkla!</blockquote>",
			firstName(c, m), c.Me.FirstName,
		)

		_, err := m.ReplyText(c, text, &td.SendTextMessageOpts{
			ParseMode:   "HTML",
			ReplyMarkup: core.AddMeMarkup(c.Me.Usernames.EditableUsername),
		})

		return err
	}

	go func(chatID int64) {
		_ = db.Instance.AddChat(chatID)
	}(chatID)

	uptime := getFormattedDuration(time.Since(startTime))
	htmlText := fmt.Sprintf(
		"<blockquote>🎵 <b>%s Aktivasyon Tamamlandı!</b> 🎉</blockquote>\n\n"+
			"<blockquote expandable>✅ Kurulumum tamamlandı ve müzik yayını yapmaya hazırım.\n\n"+
			"⏳ <b>Çalışma Süresi:</b> <code>%s</code>\n"+
			"🎶 Müzik başlatmak için <code>/play [şarkı adı]</code> yazın!</blockquote>",
		c.Me.FirstName,
		uptime,
	)

	_, err := m.ReplyText(c, htmlText, &td.SendTextMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: core.SupportBtn(),
	})

	return err
}
