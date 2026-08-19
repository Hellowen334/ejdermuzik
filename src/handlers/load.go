/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package handlers

import (
	"time"

	"github.com/AshokShau/gotdbot"
	"github.com/AshokShau/gotdbot/filters/callbackquery"
)

var startTime = time.Now()

// setBotCommands automatically registers commands with Telegram BotFather menu on startup.
func setBotCommands(c *gotdbot.Client) {
	commands := []*gotdbot.BotCommand{
		{Command: "oynat", Description: "🎵 Müzik oynatır (veya /play)"},
		{Command: "voynat", Description: "🎥 Video oynatır (veya /vplay)"},
		{Command: "duraklat", Description: "⏸️ Müziği duraklatır (veya /pause)"},
		{Command: "devam", Description: "▶️ Müziği devam ettirir (veya /resume)"},
		{Command: "atla", Description: "⏭️ Sonraki şarkıya geçer (veya /skip)"},
		{Command: "durdur", Description: "⏹️ Müziği durdurur (veya /stop)"},
		{Command: "sira", Description: "📜 Şarkı sırasını gösterir (veya /queue)"},
		{Command: "ayarlar", Description: "⚙️ Bot ayarlarını açar (veya /settings)"},
		{Command: "baslat", Description: "🤖 Botu başlatır (veya /start)"},
		{Command: "yardim", Description: "❓ Yardım menüsünü gösterir (veya /help)"},
		{Command: "ping", Description: "⚡ Bot durumunu gösterir (veya /ping)"},
	}

	go func() {
		time.Sleep(2 * time.Second)
		_, err := c.SetBotCommands(&gotdbot.BotCommandScopeDefault{}, "", commands)
		if err != nil {
			c.Logger.Warn("Failed to set bot commands", "error", err)
		} else {
			c.Logger.Info("Bot commands registered to Telegram successfully")
		}
	}()
}

// LoadModules loads all the handlers.
// It takes a telegram gotdbot.Client as input.
func LoadModules(c *gotdbot.Client) {
	setBotCommands(c)

	c.OnCommand("reload", reloadAdminCacheHandler)
	c.OnCommand("yenile", reloadAdminCacheHandler)

	c.OnCommand("authList", authListHandler)
	c.OnCommand("auths", authListHandler)
	c.OnCommand("yetkililer", authListHandler)
	c.OnCommand("auth", addAuthHandler)
	c.OnCommand("addAuth", addAuthHandler)
	c.OnCommand("yetkiver", addAuthHandler)
	c.OnCommand("removeAuth", removeAuthHandler)
	c.OnCommand("rmAuth", removeAuthHandler)
	c.OnCommand("yetkial", removeAuthHandler)

	c.OnCommand("broadcast", broadcastHandler)
	c.OnCommand("gCast", broadcastHandler)
	c.OnCommand("duyuru", broadcastHandler)
	c.OnCommand("stop_gcast", cancelBroadcastHandler)
	c.OnCommand("stop_broadcast", cancelBroadcastHandler)

	c.OnCommand("av", activeVcHandler)
	c.OnCommand("active_vc", activeVcHandler)
	c.OnCommand("aktifsesli", activeVcHandler)

	c.OnCommand("clearass", clearAssistantsHandler)
	c.OnCommand("clearAssistants", clearAssistantsHandler)
	c.OnCommand("leaveAll", leaveAllHandler)
	c.OnCommand("tumundencik", leaveAllHandler)

	c.OnCommand("logger", loggerHandler)
	c.OnCommand("privacy", privacyHandler)

	c.OnCommand("autoplay", autoplayHandler)
	c.OnCommand("otocal", autoplayHandler)
	c.OnCommand("otomatik", autoplayHandler)

	c.OnCommand("loop", loopHandler)
	c.OnCommand("dongu", loopHandler)
	c.OnCommand("döngü", loopHandler)

	c.OnCommand("pause", pauseHandler)
	c.OnCommand("duraklat", pauseHandler)

	c.OnCommand("resume", resumeHandler)
	c.OnCommand("devam", resumeHandler)
	c.OnCommand("devamet", resumeHandler)

	c.OnCommand("cplist", createPlaylistHandler)
	c.OnCommand("createplaylist", createPlaylistHandler)
	c.OnCommand("listeolustur", createPlaylistHandler)
	c.OnCommand("deleteplaylist", deletePlaylistHandler)
	c.OnCommand("listesil", deletePlaylistHandler)

	c.OnCommand("queue", queueHandler)
	c.OnCommand("sira", queueHandler)
	c.OnCommand("sıra", queueHandler)
	c.OnCommand("liste", queueHandler)

	c.OnCommand("seek", seekHandler)
	c.OnCommand("sar", seekHandler)
	c.OnCommand("sarma", seekHandler)

	c.OnCommand("sh", shellCommand)

	c.OnCommand("skip", skipHandler)
	c.OnCommand("atla", skipHandler)
	c.OnCommand("gec", skipHandler)
	c.OnCommand("geç", skipHandler)

	c.OnCommand("stop", stopHandler)
	c.OnCommand("end", stopHandler)
	c.OnCommand("durdur", stopHandler)
	c.OnCommand("bitir", stopHandler)
	c.OnCommand("sonlandir", stopHandler)

	c.OnCommand("start", startHandler)
	c.OnCommand("baslat", startHandler)
	c.OnCommand("başlat", startHandler)

	c.OnCommand("help", startHandler)
	c.OnCommand("yardim", startHandler)
	c.OnCommand("yardım", startHandler)
	c.OnCommand("komutlar", startHandler)

	c.OnCommand("ping", pingHandler)

	c.OnCommand("play", playHandler)
	c.OnCommand("p", playHandler)
	c.OnCommand("oynat", playHandler)
	c.OnCommand("cal", playHandler)
	c.OnCommand("çal", playHandler)
	c.OnCommand("sarki", playHandler)
	c.OnCommand("şarkı", playHandler)

	c.OnCommand("fplay", fPlayHandler)
	c.OnCommand("fp", fPlayHandler)
	c.OnCommand("zoynat", fPlayHandler)
	c.OnCommand("zcal", fPlayHandler)

	c.OnCommand("vplay", vPlayHandler)
	c.OnCommand("v", vPlayHandler)
	c.OnCommand("voynat", vPlayHandler)
	c.OnCommand("vcal", vPlayHandler)

	c.OnCommand("fvplay", fVPlayHandler)
	c.OnCommand("fvp", fVPlayHandler)
	c.OnCommand("zvoynat", fVPlayHandler)

	c.OnCommand("remove", removeHandler)
	c.OnCommand("cikar", removeHandler)

	c.OnCommand("mute", muteHandler)
	c.OnCommand("sessiz", muteHandler)
	c.OnCommand("sustur", muteHandler)

	c.OnCommand("unmute", unmuteHandler)
	c.OnCommand("sesac", unmuteHandler)
	c.OnCommand("sesli", unmuteHandler)

	c.OnCommand("settings", settingsHandler)
	c.OnCommand("ayarlar", settingsHandler)

	c.OnCommand("addtoplaylist", addToPlaylistHandler)
	c.OnCommand("addtoplist", addToPlaylistHandler)
	c.OnCommand("listeyeekle", addToPlaylistHandler)

	c.OnCommand("removefromplaylist", removeFromPlaylistHandler)
	c.OnCommand("rmplist", removeFromPlaylistHandler)
	c.OnCommand("listedencikar", removeFromPlaylistHandler)

	c.OnCommand("plistinfo", playlistInfoHandler)
	c.OnCommand("playlistinfo", playlistInfoHandler)
	c.OnCommand("listabilgi", playlistInfoHandler)

	c.OnCommand("myplaylists", myPlaylistsHandler)
	c.OnCommand("myplist", myPlaylistsHandler)
	c.OnCommand("listelerim", myPlaylistsHandler)

	c.OnCommand("stats", statsHandler)
	c.OnCommand("istatistik", statsHandler)

	c.OnUpdateNewCallbackQuery(helpCallbackHandler, callbackquery.Prefix("help_"))
	c.OnUpdateNewCallbackQuery(playCallbackHandler, callbackquery.Prefix("play_"))
	c.OnUpdateNewCallbackQuery(vcPlayHandler, callbackquery.Prefix("vcplay_"))
	c.OnUpdateNewCallbackQuery(settingsCallbackHandler, callbackquery.Prefix("settings_"))
	c.OnUpdateNewCallbackQuery(autoplayCallbackHandler, callbackquery.Equal("autoplay_toggle"))

	c.OnUpdateChatMember(handleParticipant, nil)
	c.OnUpdateNewMessage(handleVoiceChatMessage, nil)

	c.Logger.Debug("Handlers loaded successfully")
}
