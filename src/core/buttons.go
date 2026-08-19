/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package core

import (
	"ashokshau/tgmusic/config"
	"ashokshau/tgmusic/src/lang"
	"ashokshau/tgmusic/src/utils"
	"fmt"

	"github.com/AshokShau/gotdbot"
)

func cb(text, data string, style gotdbot.ButtonStyle) gotdbot.InlineKeyboardButton {
	return gotdbot.InlineKeyboardButton{
		Text: text,
		Type: &gotdbot.InlineKeyboardButtonTypeCallback{
			Data: []byte(data),
		},
		Style: style,
	}
}

func userId(text string, userId int64, style gotdbot.ButtonStyle) gotdbot.InlineKeyboardButton {
	return gotdbot.InlineKeyboardButton{
		Text:  text,
		Type:  &gotdbot.InlineKeyboardButtonTypeUser{UserId: userId},
		Style: style,
	}
}

func url(text, link string, style gotdbot.ButtonStyle) gotdbot.InlineKeyboardButton {
	return gotdbot.InlineKeyboardButton{
		Text: text,
		Type: &gotdbot.InlineKeyboardButtonTypeUrl{
			Url: link,
		},
		Style: style,
	}
}

var CloseBtn = cb("❤️‍🩹 Kapat ❤️‍🩹", "vcplay_close", gotdbot.ButtonStyleDanger{})
var HomeBtn = cb("🌹 Ana Menü 🌹", "help_back", gotdbot.ButtonStylePrimary{})
var HelpBtn = cb("😇 Yardım & Komutlar", "help_all", gotdbot.ButtonStylePrimary{})
var UserBtn = cb("🎧 Kullanıcı Komutları", "help_user", gotdbot.ButtonStylePrimary{})
var AdminBtn = cb("⚙️ Yönetici Komutları", "help_admin", gotdbot.ButtonStylePrimary{})
var OwnerBtn = cb("👑 Sahip Komutları", "help_owner", gotdbot.ButtonStylePrimary{})
var DevsBtn = cb("🛠️ Geliştirici Komutları", "help_devs", gotdbot.ButtonStylePrimary{})
var PlaylistBtn = cb("🎵 Çalma Listesi Komutları", "help_playlist", gotdbot.ButtonStylePrimary{})
var AutoplayBtn = cb("❤️‍🔥 Otomatik Çalma", "help_autoplay", gotdbot.ButtonStylePrimary{})

var channelBtn = url("📢 Duyurular", config.SupportChannel, gotdbot.ButtonStylePrimary{})
var groupBtn = url("💬 Destek", config.SupportGroup, gotdbot.ButtonStylePrimary{})

func SupportKeyboard() *gotdbot.ReplyMarkupInlineKeyboard {
	return &gotdbot.ReplyMarkupInlineKeyboard{
		Rows: [][]gotdbot.InlineKeyboardButton{
			{channelBtn, groupBtn},
			{CloseBtn},
		},
	}
}

func SupportBtn() *gotdbot.ReplyMarkupInlineKeyboard {
	return &gotdbot.ReplyMarkupInlineKeyboard{
		Rows: [][]gotdbot.InlineKeyboardButton{
			{channelBtn, groupBtn},
		},
	}
}

func SettingsKeyboard(playMode, adminMode string, cmdDelete bool, language string) *gotdbot.ReplyMarkupInlineKeyboard {
	playText := "Herkes"
	if playMode == utils.Admins {
		playText = "Yöneticiler"
	}

	deleteText := "Kapalı"
	if cmdDelete {
		deleteText = "Açık"
	}

	adminText := "Herkes"
	if adminMode == utils.Admins {
		adminText = "Yöneticiler"
	}

	langText := "Türkçe"
	if language != "tr" && language != "" {
		langText = lang.GetLangDisplayName(language)
	}

	return &gotdbot.ReplyMarkupInlineKeyboard{
		Rows: [][]gotdbot.InlineKeyboardButton{
			{
				cb("🎵 Oynatma Modu ➜", "settings_main", gotdbot.ButtonStylePrimary{}),
				cb(playText, "settings_play", gotdbot.ButtonStylePrimary{}),
			},
			{
				cb("🗑️ Komut Silme ➜", "settings_main", gotdbot.ButtonStylePrimary{}),
				cb(deleteText, "settings_delete", gotdbot.ButtonStylePrimary{}),
			},
			{
				cb("🛡️ Yönetici Modu ➜", "settings_main", gotdbot.ButtonStylePrimary{}),
				cb(adminText, "settings_admin", gotdbot.ButtonStylePrimary{}),
			},
			{
				cb("🌐 Dil ➜", "settings_main", gotdbot.ButtonStylePrimary{}),
				cb(langText, "settings_lang", gotdbot.ButtonStylePrimary{}),
			},
			{CloseBtn},
		},
	}
}

func HelpMenuKeyboard() *gotdbot.ReplyMarkupInlineKeyboard {
	return &gotdbot.ReplyMarkupInlineKeyboard{
		Rows: [][]gotdbot.InlineKeyboardButton{
			{UserBtn, AdminBtn},
			{OwnerBtn, DevsBtn},
			{PlaylistBtn, AutoplayBtn},
			{HomeBtn, CloseBtn},
		},
	}
}

func BackHelpMenuKeyboard() *gotdbot.ReplyMarkupInlineKeyboard {
	return &gotdbot.ReplyMarkupInlineKeyboard{
		Rows: [][]gotdbot.InlineKeyboardButton{
			{HelpBtn, HomeBtn},
			{CloseBtn},
		},
	}
}

func ControlButtons(mode string) *gotdbot.ReplyMarkupInlineKeyboard {
	skipBtn := cb("⏭️ Atla", "play_skip", gotdbot.ButtonStylePrimary{})
	stopBtn := cb("⏹️ Durdur", "play_stop", gotdbot.ButtonStyleDanger{})
	pauseBtn := cb("⏸️ Durakla", "play_pause", gotdbot.ButtonStylePrimary{})
	resumeBtn := cb("▶️ Devam", "play_resume", gotdbot.ButtonStylePrimary{})
	downloadBtn := cb("📥 İndir", "play_download", gotdbot.ButtonStylePrimary{})
	autoplayBtn := cb("❤️‍🔥 Otomatik Çal", "play_toggle_autoplay", gotdbot.ButtonStylePrimary{})
	addToPlaylistBtn := cb("⭐ Çalma Listesine Ekle", "play_add_to_list", gotdbot.ButtonStylePrimary{})
	closeBtn := cb("❤️‍🩹 Kapat ❤️‍🩹", "vcplay_close", gotdbot.ButtonStyleDanger{})

	switch mode {
	case "play":
		return &gotdbot.ReplyMarkupInlineKeyboard{
			Rows: [][]gotdbot.InlineKeyboardButton{
				{pauseBtn, stopBtn, skipBtn},
				{downloadBtn, autoplayBtn},
				{addToPlaylistBtn},
				{closeBtn},
			},
		}

	case "pause":
		return &gotdbot.ReplyMarkupInlineKeyboard{
			Rows: [][]gotdbot.InlineKeyboardButton{
				{resumeBtn, stopBtn, skipBtn},
				{downloadBtn, autoplayBtn},
				{addToPlaylistBtn},
				{closeBtn},
			},
		}

	case "resume":
		return &gotdbot.ReplyMarkupInlineKeyboard{
			Rows: [][]gotdbot.InlineKeyboardButton{
				{pauseBtn, stopBtn, skipBtn},
				{downloadBtn, autoplayBtn},
				{addToPlaylistBtn},
				{closeBtn},
			},
		}

	default:
		return &gotdbot.ReplyMarkupInlineKeyboard{
			Rows: [][]gotdbot.InlineKeyboardButton{
				{closeBtn},
			},
		}
	}
}

func AddMeMarkup(username string) *gotdbot.ReplyMarkupInlineKeyboard {
	addMeBtn := url(
		"✨ Beni Grubuna Ekle ✨",
		fmt.Sprintf("https://t.me/%s?startgroup=true", username),
		gotdbot.ButtonStylePrimary{},
	)

	return &gotdbot.ReplyMarkupInlineKeyboard{
		Rows: [][]gotdbot.InlineKeyboardButton{
			{addMeBtn},
			{HelpBtn},
			{channelBtn, groupBtn},
		},
	}
}

func DownloadKeyboard(langCode, trackID string) *gotdbot.ReplyMarkupInlineKeyboard {
	audioBtn := cb("🎵 İndir Audio", fmt.Sprintf("dl_a_%s", trackID), gotdbot.ButtonStylePrimary{})
	videoBtn := cb("📹 İndir Video", fmt.Sprintf("dl_v_%s", trackID), gotdbot.ButtonStylePrimary{})

	return &gotdbot.ReplyMarkupInlineKeyboard{
		Rows: [][]gotdbot.InlineKeyboardButton{
			{audioBtn, videoBtn},
			{CloseBtn},
		},
	}
}

func PlayNowButton(trackID string) gotdbot.InlineKeyboardButton {
	return cb("▶️ Şimdi Oynat", fmt.Sprintf("play_now_%s", trackID), gotdbot.ButtonStyleDanger{})
}

func QueueMarkup(trackID string) *gotdbot.ReplyMarkupInlineKeyboard {
	return &gotdbot.ReplyMarkupInlineKeyboard{
		Rows: [][]gotdbot.InlineKeyboardButton{
			{PlayNowButton(trackID), CloseBtn},
		},
	}
}
