/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package handlers

import (
	"context"
	"fmt"
	"html"
	"log/slog"
	"strings"
	"time"

	"ashokshau/tgmusic/src/core"
	"ashokshau/tgmusic/src/core/cache"
	"ashokshau/tgmusic/src/core/db"
	"ashokshau/tgmusic/src/utils"
	"ashokshau/tgmusic/src/vc"

	td "github.com/AshokShau/gotdbot"
)

func playCallbackHandler(c *td.Client, cb *td.UpdateNewCallbackQuery) error {
	data := cb.DataString()
	if !adminModeCB(c, cb) {
		return td.EndGroups
	}

	chatID := cb.ChatId
	user, err := c.GetUser(cb.SenderUserId)
	if err != nil {
		user = &td.User{FirstName: "Kullanıcı", Id: cb.SenderUserId}
	}

	if !cache.ChatCache.IsActive(chatID) {
		text := "<blockquote>⚠️ Aktif bir yayın bulunmuyor.</blockquote>"
		_ = cb.Answer(c, 0, false, "Aktif bir yayın bulunmuyor.", "")
		_, _ = cb.EditMessageText(c, text, &td.EditTextMessageOpts{ReplyMarkup: core.ControlButtons(""), ParseMode: "HTML", DisableWebPagePreview: true})
		return nil
	}

	currentTrack := cache.ChatCache.GetPlayingTrack(chatID)
	if currentTrack == nil {
		_ = cb.Answer(c, 0, false, "Aktif bir yayın bulunmuyor.", "")
		_, _ = cb.EditMessageText(c, "<blockquote>⚠️ Aktif bir yayın bulunmuyor.</blockquote>", &td.EditTextMessageOpts{ReplyMarkup: core.ControlButtons(""), ParseMode: "HTML", DisableWebPagePreview: true})
		return nil
	}

	buildTrackMessage := func(status, emoji string) string {
		escURL := html.EscapeString(currentTrack.URL)
		escName := html.EscapeString(currentTrack.Name)
		escUser := html.EscapeString(currentTrack.User)
		return fmt.Sprintf("%s <b>%s</b>\n\n<b>Parça:</b> <a href='%s'>%s</a>\n<b>Süre:</b> %s\n<b>İsteyen:</b> %s",
			emoji, status,
			escURL, escName,
			utils.SecToMin(currentTrack.Duration),
			escUser,
		)
	}

	switch {
	case strings.Contains(data, "play_skip"):
		if err := vc.Calls.PlayNext(c, chatID); err != nil {
			_ = cb.Answer(c, 0, false, "Şarkı atlanamadı.", "")
			_, _ = cb.EditMessageText(c, "<blockquote>⚠️ Şarkı atlanamadı.</blockquote>", &td.EditTextMessageOpts{ReplyMarkup: core.ControlButtons(""), ParseMode: "HTML", DisableWebPagePreview: true})
			return nil
		}
		_ = cb.Answer(c, 0, false, "Şarkı atlandı ⏭️", "")
		_ = c.DeleteMessages(chatID, []int64{cb.MessageId}, &td.DeleteMessagesOpts{Revoke: true})
		return nil

	case strings.Contains(data, "play_stop"):
		if err := vc.Calls.Stop(chatID, false); err != nil {
			_ = cb.Answer(c, 0, false, "Oynatma durdurulamadı.", "")
			_, _ = cb.EditMessageText(c, "<blockquote>⚠️ Oynatma durdurulamadı.</blockquote>", &td.EditTextMessageOpts{ReplyMarkup: core.ControlButtons(""), ParseMode: "HTML", DisableWebPagePreview: true})
			return nil
		}

		msg := fmt.Sprintf("<blockquote>⏹️ <b>Oynatma durduruldu.</b>\nDurduran: %s</blockquote>", html.EscapeString(user.FirstName))
		_ = cb.Answer(c, 0, false, "Oynatma durduruldu.", "")
		_, err := cb.EditMessageText(c, msg, &td.EditTextMessageOpts{ReplyMarkup: core.ControlButtons(""), ParseMode: "HTML", DisableWebPagePreview: true})
		return err

	case strings.Contains(data, "play_pause"):
		if _, err = vc.Calls.Pause(chatID); err != nil {
			_ = cb.Answer(c, 0, false, "Oynatma duraklatılamadı.", "")
			_, _ = cb.EditMessageText(c, "<blockquote>⚠️ Oynatma duraklatılamadı.</blockquote>", &td.EditTextMessageOpts{ReplyMarkup: core.ControlButtons(""), ParseMode: "HTML", DisableWebPagePreview: true})
			return nil
		}
		_ = cb.Answer(c, 0, false, "Oynatma duraklatıldı ⏸️", "")
		text := buildTrackMessage("Duraklatıldı", "⏸") + fmt.Sprintf("\n\nDuraklatan: %s", html.EscapeString(user.FirstName))
		_, _ = cb.EditMessageText(c, text, &td.EditTextMessageOpts{ReplyMarkup: core.ControlButtons("pause"), ParseMode: "HTML", DisableWebPagePreview: true})
		return nil

	case strings.Contains(data, "play_resume"):
		if _, err := vc.Calls.Resume(chatID); err != nil {
			_ = cb.Answer(c, 0, false, "Oynatma devam ettirilemedi.", "")
			_, _ = cb.EditMessageText(c, "<blockquote>⚠️ Oynatma devam ettirilemedi.</blockquote>", &td.EditTextMessageOpts{ReplyMarkup: core.ControlButtons("pause"), ParseMode: "HTML", DisableWebPagePreview: true})
			return nil
		}
		_ = cb.Answer(c, 0, false, "Oynatma devam ediyor ▶️", "")
		text := buildTrackMessage("Şimdi Oynatılıyor", "▶") + fmt.Sprintf("\n\nDevam ettiren: %s", html.EscapeString(user.FirstName))
		_, _ = cb.EditMessageText(c, text, &td.EditTextMessageOpts{ReplyMarkup: core.ControlButtons("resume"), ParseMode: "HTML", DisableWebPagePreview: true})
		return nil

	case strings.Contains(data, "play_toggle_autoplay"):
		state := cache.ChatCache.GetAutoplay(chatID)
		newState := !state
		cache.ChatCache.SetAutoplay(chatID, newState)
		var statusText string
		if newState {
			statusText = "Otomatik Çalma AÇILDI ✅"
		} else {
			statusText = "Otomatik Çalma KAPATILDI ❌"
		}
		_ = cb.Answer(c, 0, true, statusText, "")
		return nil

	case strings.Contains(data, "play_download"):
		if currentTrack == nil {
			_ = cb.Answer(c, 0, true, "❌ Şu an çalan bir parça yok.", "")
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var fileID string
		if currentTrack.TrackID != "" && db.Instance != nil {
			fileID, _, _ = db.Instance.GetCachedFileID(ctx, currentTrack.TrackID, false)
		}

		caption := fmt.Sprintf("🎵 <b><a href='%s'>%s</a></b>\n\n📥 <i>İndirme tamamlandı! İyi dinlemeler.</i>", html.EscapeString(currentTrack.URL), html.EscapeString(currentTrack.Name))
		formattedCaption := &td.FormattedText{Text: caption}
		if parsedText, parseErr := c.ParseTextEntities(&td.TextParseModeHTML{}, caption); parseErr == nil && parsedText != nil {
			formattedCaption = parsedText
		}

		// Try 1: Send via Telegram cloud file_id if available
		if fileID != "" {
			_, err := c.SendMessage(cb.SenderUserId, &td.InputMessageAudio{
				Audio: &td.InputAudio{
					Audio:    &td.InputFileRemote{Id: fileID},
					Title:    currentTrack.Name,
					Duration: int32(currentTrack.Duration),
				},
				Caption: formattedCaption,
			}, nil)

			if err == nil {
				_ = cb.Answer(c, 0, true, "📥 Şarkı özel mesajınıza (DM) MP3 olarak gönderildi!", "")
				return nil
			}
		}

		// Try 2: Send via local file path if file exists on server
		if currentTrack.FilePath != "" {
			if _, statErr := os.Stat(currentTrack.FilePath); statErr == nil {
				_, err := c.SendMessage(cb.SenderUserId, &td.InputMessageAudio{
					Audio: &td.InputAudio{
						Audio:    &td.InputFileLocal{Path: currentTrack.FilePath},
						Title:    currentTrack.Name,
						Duration: int32(currentTrack.Duration),
					},
					Caption: formattedCaption,
				}, nil)

				if err == nil {
					_ = cb.Answer(c, 0, true, "📥 Şarkı özel mesajınıza (DM) MP3 olarak gönderildi!", "")
					return nil
				}
			}
		}

		// Try 3: Download track locally and send
		localPath, err := dl.DownloadCachedTrack(currentTrack, c)
		if err == nil && localPath != "" {
			_, err = c.SendMessage(cb.SenderUserId, &td.InputMessageAudio{
				Audio: &td.InputAudio{
					Audio:    &td.InputFileLocal{Path: localPath},
					Title:    currentTrack.Name,
					Duration: int32(currentTrack.Duration),
				},
				Caption: formattedCaption,
			}, nil)

			if err == nil {
				_ = cb.Answer(c, 0, true, "📥 Şarkı özel mesajınıza (DM) MP3 olarak gönderildi!", "")
				return nil
			}
		}

		botUser, _ := c.GetMe()
		botUsername := "bot"
		if botUser != nil && botUser.Usernames != nil {
			botUsername = botUser.Usernames.EditableUsername
		}

		_ = cb.Answer(c, 0, true, fmt.Sprintf("📥 Şarkıyı DM'den alabilmek için lütfen önce botla özel sohbet başlatın: @%s", botUsername), "")
		return nil

	case strings.Contains(data, "play_add_to_list"):
		playlists, err := db.Instance.GetUserPlaylists(cb.SenderUserId)
		if err != nil {
			_ = cb.Answer(c, 0, true, "❌ Çalma listeleri alınamadı.", "")
			return nil
		}

		var playlistID string
		if len(playlists) == 0 {
			playlistID, err = db.Instance.CreatePlaylist("Favorilerim", cb.SenderUserId)
			if err != nil {
				_ = cb.Answer(c, 0, true, "❌ Çalma listesi oluşturulamadı.", "")
				return nil
			}
		} else {
			playlistID = playlists[0].ID
		}

		song := db.Song{
			URL:      currentTrack.URL,
			Name:     currentTrack.Name,
			TrackID:  currentTrack.TrackID,
			Duration: currentTrack.Duration,
			Platform: currentTrack.Platform,
		}

		err = db.Instance.AddSongToPlaylist(playlistID, song)
		if err != nil {
			_ = cb.Answer(c, 0, true, "❌ Parça listeye eklenemedi.", "")
			return nil
		}

		playlist, err := db.Instance.GetPlaylist(playlistID)
		pName := "Çalma Listem"
		if err == nil && playlist != nil {
			pName = playlist.Name
		}

		notifyMsg := fmt.Sprintf("<blockquote>⭐ <b>Parça Çalma Listesine Eklendi!</b>\n\n🎵 <b>Şarkı:</b> %s\n📁 <b>Liste:</b> %s</blockquote>", html.EscapeString(song.Name), html.EscapeString(pName))
		_, _ = c.SendTextMessage(chatID, notifyMsg, &td.SendTextMessageOpts{ParseMode: "HTML"})

		_ = cb.Answer(c, 0, true, fmt.Sprintf("✅ \"%s\" parçası '%s' listenize eklendi!", song.Name, pName), "")
		return nil

	case strings.HasPrefix(data, "play_now_"):
		trackID := strings.TrimPrefix(data, "play_now_")
		if ok := cache.ChatCache.MoveTrackToFront(chatID, trackID); !ok {
			_ = cb.Answer(c, 0, false, "Parça sırada bulunamadı.", "")
			return nil
		}
		if err := vc.Calls.PlayNext(c, chatID); err != nil {
			_ = cb.Answer(c, 0, false, "Parça oynatılamadı.", "")
			return nil
		}
		_ = cb.Answer(c, 0, false, "Şimdi oynatılıyor ▶️", "")
		_ = c.DeleteMessages(chatID, []int64{cb.MessageId}, &td.DeleteMessagesOpts{Revoke: true})
		return nil
	}

	text := buildTrackMessage("Şimdi Oynatılıyor", "▶")
	_, _ = cb.EditMessageText(c, text, &td.EditTextMessageOpts{ReplyMarkup: core.ControlButtons("resume"), ParseMode: "HTML", DisableWebPagePreview: true})
	return nil
}

func vcPlayHandler(c *td.Client, cb *td.UpdateNewCallbackQuery) error {
	data := cb.DataString()

	if strings.Contains(data, "vcplay_close") {
		_ = cb.Answer(c, 0, false, "Panel kapatıldı.", "")
		_ = c.DeleteMessages(cb.ChatId, []int64{cb.MessageId}, &td.DeleteMessagesOpts{Revoke: true})
		return nil
	}

	slog.Info("Received vcplay callback", "arg1", data)
	return nil
}
