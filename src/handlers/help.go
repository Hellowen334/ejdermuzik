/*
 * TgMusicBot - Telegram Music Bot
 *  Copyright (c) 2025-2026 Ashok Shau
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/AshokShau/TgMusicBot
 */

package handlers

import (
	"ashokshau/tgmusic/config"
	"fmt"
	"strings"

	"ashokshau/tgmusic/src/core"

	td "github.com/AshokShau/gotdbot"
)

func getHelpCategories() map[string]struct {
	Title   string
	Content string
	Markup  *td.ReplyMarkupInlineKeyboard
} {
	return map[string]struct {
		Title   string
		Content string
		Markup  *td.ReplyMarkupInlineKeyboard
	}{
		"help_user": {
			Title: "Kullanıcı Komutları",
			Content: `<p>Gruptaki tüm üyeler tarafından kullanılabilen komutlar.</p>

<details open>
  <summary>Oynatma</summary>
  <table bordered striped>
    <tr><th>Komut</th><th>Açıklama</th></tr>
    <tr><td><code>/play [şarkı]</code></td><td>YouTube, Spotify ve diğer platformlardan müzik çalın.</td></tr>
    <tr><td><code>/vplay [şarkı]</code></td><td>Sesli/görüntülü sohbette video oynatın.</td></tr>
    <tr><td><code>/fplay [şarkı]</code></td><td>Parçayı sırayı atlayarak hemen çalın.</td></tr>
    <tr><td><code>/fvplay [şarkı]</code></td><td>Videoyu sırayı atlayarak hemen oynatın.</td></tr>
  </table>
</details>

<details>
  <summary>Genel</summary>
  <table bordered striped>
    <tr><th>Komut</th><th>Açıklama</th></tr>
    <tr><td><code>/start</code></td><td>Botu başlatın veya aktifliğini kontrol edin.</td></tr>
    <tr><td><code>/help</code></td><td>Etkileşimli yardım menüsünü açın.</td></tr>
    <tr><td><code>/ping</code></td><td>Bot gecikmesini ve sistem durumunu gösterir.</td></tr>
    <tr><td><code>/privacy</code></td><td>Gizlilik politikasını görüntüleyin.</td></tr>
    <tr><td><code>/queue</code></td><td>Mevcut çalma sırasını gösterir.</td></tr>
  </table>
</details>`,
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_admin": {
			Title: "Yönetici Komutları",
			Content: `<p>Yöneticiler ve yetkili kullanıcılar tarafından kullanılabilen komutlar.</p>

<details open>
  <summary>Oynatma Kontrolleri</summary>
  <table bordered striped>
    <tr><th>Komut</th><th>Açıklama</th></tr>
    <tr><td><code>/skip</code></td><td>Mevcut çalan parçayı atla.</td></tr>
    <tr><td><code>/pause</code></td><td>Oynatmayı duraklat.</td></tr>
    <tr><td><code>/resume</code></td><td>Oynatmayı devam ettir.</td></tr>
    <tr><td><code>/seek [saniye]</code></td><td>Belirli bir konuma atla.</td></tr>
    <tr><td><code>/mute</code></td><td>Sesi sessize al.</td></tr>
    <tr><td><code>/unmute</code></td><td>Sesli sohbetin sesini aç.</td></tr>
  </table>
</details>

<details>
  <summary>Sıra & İzinler</summary>
  <table bordered striped>
    <tr><th>Komut</th><th>Açıklama</th></tr>
    <tr><td><code>/remove [sıra_no]</code></td><td>Sıradan bir parçayı kaldır.</td></tr>
    <tr><td><code>/loop [0-10]</code></td><td>Parçayı belirtilen sayıda tekrarla.</td></tr>
    <tr><td><code>/auth</code></td><td>Bir kullanıcıya yönetici yetkisi ver.</td></tr>
    <tr><td><code>/unauth</code></td><td>Kullanıcının yetkisini kaldır.</td></tr>
    <tr><td><code>/authlist</code></td><td>Yetkili kullanıcıları listele.</td></tr>
  </table>
</details>`,
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_devs": {
			Title: "Geliştirici Araçları",
			Content: `<p>Bot geliştiricileri ve yöneticileri için özel komutlar.</p>

<details open>
  <summary>Sistem</summary>
  <table bordered striped>
    <tr><th>Komut</th><th>Açıklama</th></tr>
    <tr><td><code>/stats</code></td><td>Sistem ve sunucu istatistiklerini gösterir.</td></tr>
    <tr><td><code>/av</code></td><td>Aktif sesli sohbetleri listeler.</td></tr>
    <tr><td><code>/clearass</code></td><td>Tüm asistan bağlantılarını temizler.</td></tr>
    <tr><td><code>/leaveall</code></td><td>Asistanı tüm sohbetlerden çıkarır.</td></tr>
    <tr><td><code>/logger</code></td><td>Loglama yapılandırmasını gösterir.</td></tr>
  </table>
</details>`,
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_owner": {
			Title: "Sahip Komutları",
			Content: `<p>Grup sahibi için yapılandırma seçenekleri.</p>

<details open>
  <summary>Sohbet Ayarları</summary>
  <table bordered striped>
    <tr><th>Komut</th><th>Açıklama</th></tr>
    <tr><td><code>/settings</code></td><td>Oynatma modu, yönetici modu, komut silme ve dil ayarlarını yönetin.</td></tr>
  </table>
</details>`,
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_playlist": {
			Title: "Çalma Listesi Komutları",
			Content: `<p>Kişisel çalma listelerinizi oluşturun ve yönetin.</p>

<details open>
  <summary>Çalma Listesi Yönetimi</summary>
  <table bordered striped>
    <tr><th>Komut</th><th>Açıklama</th></tr>
    <tr><td><code>/createplaylist</code></td><td>Yeni çalma listesi oluştur.</td></tr>
    <tr><td><code>/deleteplaylist</code></td><td>Çalma listesini sil.</td></tr>
    <tr><td><code>/addtoplaylist</code></td><td>Çalma listesine şarkı ekle.</td></tr>
    <tr><td><code>/removefromplaylist</code></td><td>Çalma listesinden şarkı çıkar.</td></tr>
    <tr><td><code>/playlistinfo</code></td><td>Çalma listesi bilgilerini göster.</td></tr>
    <tr><td><code>/myplaylists</code></td><td>Tüm çalma listelerini listele.</td></tr>
  </table>
</details>`,
			Markup: core.BackHelpMenuKeyboard(),
		},
		"help_autoplay": {
			Title: "Otomatik Çalma Komutları",
			Content: `<p>Önerilen parçalarla otomatik oynatmaya devam edin.</p>

<details open>
  <summary>Otomatik Çalma</summary>
  <table bordered striped>
    <tr><th>Komut</th><th>Açıklama</th></tr>
    <tr><td><code>/autoplay</code></td><td>Otomatik çalmayı aç/kapat. Açık olduğunda sıra bittiğinde otomatik önerilen parçalar eklenir.</td></tr>
  </table>
</details>`,
			Markup: core.BackHelpMenuKeyboard(),
		},
	}
}

func helpCallbackHandler(c *td.Client, cb *td.UpdateNewCallbackQuery) error {
	data := cb.DataString()

	user, err := c.GetUser(cb.SenderUserId)
	if err != nil {
		user = &td.User{FirstName: "Kullanıcı", Id: cb.SenderUserId}
	}

	helpCategories := getHelpCategories()

	if strings.Contains(data, "help_all") {
		_ = cb.Answer(c, 0, false, "Yardım menüsü açılıyor...", "")
		response := fmt.Sprintf(
			"<h3>Hoş Geldiniz, %s!</h3>\n"+
				"<p><b>%s</b> Telegram sesli ve görüntülü sohbetleriniz için hızlı ve gelişmiş bir müzik botudur.</p>\n\n"+
				"<p><b>Desteklenen platformlar:</b> YouTube, Spotify, Apple Music, SoundCloud, Deezer vb.</p>\n\n"+
				"<p>Aşağıdaki kategorilerden birini seçerek komutları inceleyebilirsiniz.</p>",
			user.FirstName,
			c.Me.FirstName,
		)

		richMessage := &td.InputRichMessage{
			Source: &td.RichMessageSourceHtml{
				Text: response,
			},
		}
		_, _ = c.EditMessageText(cb.ChatId, &td.InputMessageRichMessage{Message: richMessage}, cb.MessageId, &td.EditMessageTextOpts{ReplyMarkup: core.HelpMenuKeyboard()})
		return nil
	}

	if strings.Contains(data, "help_back") {
		_ = cb.Answer(c, 0, false, "Ana menüye dönülüyor...", "")

		response := fmt.Sprintf(
			"<img src=\"%s\"/>\n"+
				"<h3>Hoş Geldiniz, %s!</h3>\n"+
				"<p><b>%s</b> gruplarınızda yüksek kalitede müzik ve video oynatmanızı sağlar.</p>\n\n"+
				"<p><b>Desteklenen platformlar:</b> YouTube, Spotify, Apple Music, SoundCloud vb.</p>\n\n"+
				"<p>Botu grubunuza eklemek veya komutları görmek için aşağıdaki düğmeleri kullanın.</p>",
			config.StartImg,
			user.FirstName,
			c.Me.FirstName,
		)

		richMessage := &td.InputRichMessage{
			Source: &td.RichMessageSourceHtml{
				Text: response,
			},
		}
		_, _ = c.EditMessageText(cb.ChatId, &td.InputMessageRichMessage{Message: richMessage}, cb.MessageId, &td.EditMessageTextOpts{ReplyMarkup: core.AddMeMarkup(c.Me.Usernames.EditableUsername)})
		return nil
	}

	if category, ok := helpCategories[data]; ok {
		_ = cb.Answer(c, 0, false, category.Title, "")
		response := fmt.Sprintf("<h3>%s</h3>\n\n%s\n\n<i>Geri dönmek için aşağıdaki düğmeleri kullanın.</i>", category.Title, category.Content)
		richMessage := &td.InputRichMessage{
			Source: &td.RichMessageSourceHtml{
				Text: response,
			},
		}
		_, _ = c.EditMessageText(cb.ChatId, &td.InputMessageRichMessage{Message: richMessage}, cb.MessageId, &td.EditMessageTextOpts{ReplyMarkup: category.Markup})
		return nil
	}

	_ = cb.Answer(c, 0, true, "Bilinmeyen yardım kategorisi.", "")
	return nil
}
