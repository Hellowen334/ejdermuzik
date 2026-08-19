/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package dl

import (
	"ashokshau/tgmusic/config"
	"ashokshau/tgmusic/src/core/db"
	"ashokshau/tgmusic/src/utils"
	"context"
	"fmt"
	"log/slog"
	"time"

	td "github.com/AshokShau/gotdbot"
)

func DownloadCachedTrack(cached *utils.CachedTrack, bot *td.Client) (string, error) {
	if cached.Platform == utils.DirectLink {
		return cached.URL, nil
	}

	if cached.Platform == utils.Telegram {
		return downloadTelegramFile(cached, bot)
	}

	// Step 1: Check Telegram file_id cache in MongoDB
	if cached.TrackID != "" && db.Instance != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		fileID, msgID, err := db.Instance.GetCachedFileID(ctx, cached.TrackID, cached.IsVideo)
		cancel()

		if err == nil && fileID != "" {
			localPath, err := downloadTelegramFileID(fileID, bot)
			if err == nil && localPath != "" {
				slog.Info("[Cache] Track retrieved instantly from Telegram file_id cache", "trackID", cached.TrackID, "fileID", fileID)
				return localPath, nil
			}

			// Fallback: If file_id expired, attempt to fetch message from Logger channel
			if msgID != 0 && config.LoggerId != 0 {
				msg, err := bot.GetMessage(config.LoggerId, int64(msgID))
				if err == nil && msg != nil {
					file, err := msg.Download(bot, 1, 0, 0, true)
					if err == nil && file != nil && file.Local != nil && file.Local.Path != "" {
						if file.Remote != nil && file.Remote.Id != "" {
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							_ = db.Instance.SaveCachedFileID(ctx, cached.TrackID, cached.IsVideo, file.Remote.Id, msgID)
							cancel()
						}
						slog.Info("[Cache] Track retrieved from Logger channel message", "trackID", cached.TrackID, "msgID", msgID)
						return file.Local.Path, nil
					}
				}
			}
		}
	}

	// Step 2: Download via API / Wrapper / yt-dlp if not cached in Telegram
	dlBot := bot
	if DlBot != nil {
		dlBot = DlBot
	}

	return downloadViaWrapper(cached, dlBot)
}

func downloadTelegramFileID(fileID string, bot *td.Client) (string, error) {
	file, err := bot.GetRemoteFile(fileID, nil)
	if err != nil {
		return "", err
	}

	download, err := file.Download(bot, 0, 0, 1, &td.DownloadFileOpts{Synchronous: true})
	if err != nil {
		return "", err
	}

	if download == nil || download.Local == nil {
		return "", fmt.Errorf("failed to download file from Telegram remote file_id")
	}

	return download.Local.Path, nil
}

func downloadViaWrapper(cached *utils.CachedTrack, dlBot *td.Client) (string, error) {
	wrapper := NewDownloaderWrapper(cached.URL)
	if !wrapper.IsValid() {
		return "", fmt.Errorf("invalid cached URL: %s", cached.URL)
	}

	track, err := wrapper.GetTrack()
	if err != nil {
		return "", fmt.Errorf("get track info: %w", err)
	}

	path, err := wrapper.DownloadTrack(track, cached.IsVideo)
	if err != nil {
		return "", err
	}

	if utils.TelegramMessageRegex.MatchString(path) {
		return downloadFromTelegramMessage(dlBot, path)
	}

	return path, nil
}

func downloadTelegramFile(cached *utils.CachedTrack, bot *td.Client) (string, error) {
	file, err := bot.GetRemoteFile(cached.TrackID, nil)
	if err != nil {
		return "", err
	}

	download, err := file.Download(bot, 0, 0, 1, &td.DownloadFileOpts{Synchronous: true})
	if err != nil {
		return "", err
	}

	return download.Local.Path, nil
}

func downloadFromTelegramMessage(bot *td.Client, msgURL string) (string, error) {
	msg, err := utils.GetMessage(bot, msgURL)
	if err != nil {
		return "", fmt.Errorf("get telegram message: %w", err)
	}

	file, err := msg.Download(bot, 1, 0, 0, true)
	if err != nil {
		return "", err
	}

	if file == nil || file.Local == nil {
		return "", fmt.Errorf("failed to download file from Telegram message")
	}

	return file.Local.Path, nil
}
