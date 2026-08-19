/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package vc

import (
	"ashokshau/tgmusic/config"
	"ashokshau/tgmusic/src/core/db"
	"ashokshau/tgmusic/src/utils"
	"context"
	"fmt"
	"os"
	"time"

	td "github.com/AshokShau/gotdbot"
)

// sendLogger sends a formatted log message and uploads the track to the designated logger chat.
// It caches the Telegram file_id and message_id in MongoDB for instant future replay.
func sendLogger(client *td.Client, chatID int64, song *utils.CachedTrack) {
	if chatID == 0 || song == nil || config.LoggerId == 0 || chatID == config.LoggerId {
		return
	}

	caption := fmt.Sprintf(
		"<b>A song is playing</b> in <code>%d</code>\n\n‣ <b>Title:</b> <a href='%s'>%s</a>\n‣ <b>Duration:</b> %s\n‣ <b>Requested by:</b> %s\n‣ <b>Platform:</b> %s\n‣ <b>Is Video:</b> %t",
		chatID,
		song.URL,
		song.Name,
		utils.SecToMin(song.Duration),
		song.User,
		song.Platform,
		song.IsVideo,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// If already cached in Telegram cloud, just send text log
	if song.TrackID != "" && db.Instance != nil {
		existingFileID, _, _ := db.Instance.GetCachedFileID(ctx, song.TrackID, song.IsVideo)
		if existingFileID != "" {
			_, err := client.SendTextMessage(config.LoggerId, caption, &td.SendTextMessageOpts{DisableWebPagePreview: true, ParseMode: "HTML"})
			if err != nil {
				logger.Warn("Failed to send text log message", "error", err)
			}
			return
		}
	}

	var (
		msg *td.Message
		err error
	)

	formattedCaption := &td.FormattedText{Text: caption}
	parsedText, parseErr := client.ParseTextEntities(&td.TextParseModeHTML{}, caption)
	if parseErr == nil && parsedText != nil {
		formattedCaption = parsedText
	}

	// If local file exists, upload to Logger channel to obtain permanent Telegram file_id
	if song.FilePath != "" {
		if _, statErr := os.Stat(song.FilePath); statErr == nil {
			if song.IsVideo {
				msg, err = client.SendMessage(config.LoggerId, &td.InputMessageVideo{
					Video: &td.InputVideo{
						Video:    &td.InputFileLocal{Path: song.FilePath},
						Duration: int32(song.Duration),
					},
					Caption: formattedCaption,
				}, nil)
			} else {
				msg, err = client.SendMessage(config.LoggerId, &td.InputMessageAudio{
					Audio: &td.InputAudio{
						Audio:    &td.InputFileLocal{Path: song.FilePath},
						Title:    song.Name,
						Duration: int32(song.Duration),
					},
					Caption: formattedCaption,
				}, nil)
			}
		}
	}

	// Fallback to text message if media upload wasn't executed
	if msg == nil {
		if song.Platform == utils.Telegram && song.TrackID != "" && db.Instance != nil {
			_ = db.Instance.SaveCachedFileID(ctx, song.TrackID, song.IsVideo, song.TrackID, 0)
		}
		_, err = client.SendTextMessage(config.LoggerId, caption, &td.SendTextMessageOpts{DisableWebPagePreview: true, ParseMode: "HTML"})
		if err != nil {
			logger.Warn("Failed to send fallback text log", "error", err)
		}
		return
	}

	if err != nil {
		logger.Warn("Failed to send logger media message", "error", err)
		return
	}

	// Extract Telegram file_id and message_id and save to MongoDB
	if msg != nil && song.TrackID != "" && db.Instance != nil {
		var (
			fileID    string
			messageID int32 = int32(msg.Id)
		)

		switch content := msg.Content.(type) {
		case *td.MessageAudio:
			if content.Audio != nil && content.Audio.Audio != nil && content.Audio.Audio.Remote != nil {
				fileID = content.Audio.Audio.Remote.Id
			}
		case *td.MessageVideo:
			if content.Video != nil && content.Video.Video != nil && content.Video.Video.Remote != nil {
				fileID = content.Video.Video.Remote.Id
			}
		case *td.MessageDocument:
			if content.Document != nil && content.Document.Document != nil && content.Document.Document.Remote != nil {
				fileID = content.Document.Document.Remote.Id
			}
		}

		if fileID != "" {
			_ = db.Instance.SaveCachedFileID(ctx, song.TrackID, song.IsVideo, fileID, messageID)
			logger.Info("Saved track to Telegram file_id cache", "trackID", song.TrackID, "fileID", fileID, "msgID", messageID)
		}
	}
}
