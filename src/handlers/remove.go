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
	"sort"
	"strconv"
	"strings"

	"ashokshau/tgmusic/src/core/cache"

	td "github.com/AshokShau/gotdbot"
)

// removeHandler handles the /remove command.
func removeHandler(c *td.Client, m *td.Message) error {
	if !adminMode(c, m) {
		return td.EndGroups
	}

	chatID := m.ChatId

	if !cache.ChatCache.IsActive(chatID) {
		_, _ = m.ReplyText(c, "<blockquote>⚠️ Bot sesli sohbette yayın yapmıyor.</blockquote>", &td.SendTextMessageOpts{ParseMode: "HTML"})
		return nil
	}

	queue := cache.ChatCache.GetQueue(chatID)
	if len(queue) == 0 {
		_, _ = m.ReplyText(c, "<blockquote>🎵 Oynatma sırası boş.</blockquote>", &td.SendTextMessageOpts{ParseMode: "HTML"})
		return nil
	}

	args := Args(m)
	if args == "" {
		_, _ = m.ReplyText(c, "<blockquote><b>Kullanım:</b> <code>/remove [şarkı numarası veya aralık]</code>\n\nÖrnekler:\n- <code>/remove 1</code> (#1 numaralı şarkıyı siler)\n- <code>/remove 1-5</code> (1 ile 5 arasındaki şarkıları siler)\n- <code>/remove 1,3,5</code> (1, 3 ve 5. şarkıları siler)</blockquote>", replyOpts)
		return nil
	}

	var tracksToRemove []int
	for _, part := range strings.Split(args, ",") {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				continue
			}
			start, err1 := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err1 == nil && err2 == nil {
				if start > end {
					start, end = end, start
				}

				if start < 1 {
					start = 1
				}
				if end >= len(queue) {
					end = len(queue)
				}
				for i := start; i <= end; i++ {
					tracksToRemove = append(tracksToRemove, i)
				}
			}
		} else {
			if val, err := strconv.Atoi(part); err == nil {
				tracksToRemove = append(tracksToRemove, val)
			}
		}
	}

	if len(tracksToRemove) == 0 {
		_, _ = m.ReplyText(c, "<blockquote>⚠️ Lütfen geçerli bir şarkı numarası veya aralık girin.</blockquote>", replyOpts)
		return nil
	}

	uniqueTracks := make(map[int]bool)
	var sortedTracks []int
	for _, t := range tracksToRemove {
		if t > 0 && t <= len(queue) && !uniqueTracks[t] {
			uniqueTracks[t] = true
			sortedTracks = append(sortedTracks, t)
		}
	}

	if len(sortedTracks) == 0 {
		_, _ = m.ReplyText(c, "<blockquote>⚠️ Çıkarılacak geçerli şarkı bulunamadı.</blockquote>", replyOpts)
		return nil
	}

	sort.Slice(sortedTracks, func(i, j int) bool {
		return sortedTracks[i] > sortedTracks[j]
	})

	for _, t := range sortedTracks {
		cache.ChatCache.RemoveTrack(chatID, t)
	}

	var err error
	if len(sortedTracks) == 1 {
		_, err = m.ReplyText(c, fmt.Sprintf("<blockquote>🗑️ <b>#%d numaralı şarkı %s tarafından sıradan çıkarıldı.</b></blockquote>", sortedTracks[0], firstName(c, m)), replyOpts)
	} else {
		_, err = m.ReplyText(c, fmt.Sprintf("<blockquote>🗑️ <b>%d parça %s tarafından sıradan çıkarıldı.</b></blockquote>", len(sortedTracks), firstName(c, m)), replyOpts)
	}

	return err
}
