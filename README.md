# Thor Scraper — Dark Web CTI Tool

![Go](https://img.shields.io/badge/Go-1.20%2B-00ADD8)
![Tor](https://img.shields.io/badge/network-Tor-7D4698)
![License](https://img.shields.io/badge/license-MIT-green)

<p align="center"><b><a href="#english">English</a></b> · <b><a href="#türkçe">Türkçe</a></b></p>

---

## English

**Thor Scraper** is a Go-based Cyber Threat Intelligence (CTI) tool for monitoring Tor (`.onion`) sites: it checks whether targets are reachable and takes screenshots as evidence. It automates manual intelligence-gathering so analysts work faster and more safely.

### Key features

- **Anonymity:** routes all traffic through a local **Tor SOCKS5 proxy (127.0.0.1:9150)**, so your real IP is hidden.
- **Visual evidence:** saves a high-quality `.png` screenshot of each visited site.
- **Automatic reporting:** logs each target's status (success/failure) and response time to a detailed log file.
- **Fault tolerance:** unresponsive or down sites don't stop the run — the error is handled and scanning continues.

### Requirements

- Go 1.20+
- Tor Browser running in the background (to tunnel into the Tor network)

### Setup

```bash
git clone https://github.com/KaanTuran28/CTI_Scrapper_2.git
cd CTI_Scrapper_2
go mod tidy
```

### Usage

1. **Tor connection:** start Tor Browser and click "Connect". Leave it running (minimised).
2. **Targets:** open `targets.yaml` and add the `.onion` addresses you want to monitor:

   ```yaml
   - name: "Example Forum"
     url: "http://exampleaddress...onion/"
   ```

3. Run the tool.

### Output

- `screenshots/` — evidence screenshots of the sites.
- `logs/scan_report.log` — scan results and status report.
- `logs/app.log` — technical runtime logs.

> For authorised CTI research and defensive monitoring only.

### License

MIT — see [LICENSE](./LICENSE).

---

## Türkçe

**Thor Scraper**, Tor ağı üzerindeki (`.onion`) web sitelerini otomatik izlemek, erişilebilirlik durumlarını kontrol etmek ve kanıt amaçlı ekran görüntüsü almak için geliştirilmiş, **Go** tabanlı bir Siber Tehdit İstihbaratı (CTI) aracıdır. Manuel istihbarat toplama süreçlerini otomatize ederek analistlere hız ve güvenlik kazandırır.

### Temel özellikler

- **Tam gizlilik:** tüm ağ trafiğini yerel **Tor SOCKS5 Proxy (127.0.0.1:9150)** üzerinden geçirerek gerçek IP adresinizi gizler.
- **Görsel kanıt:** ziyaret edilen sitelerin yüksek kaliteli ekran görüntüsünü (`.png`) kaydeder.
- **Otomatik raporlama:** taranan hedeflerin durumunu (başarılı/başarısız) ve yanıt sürelerini ayrıntılı bir log dosyasına işler.
- **Hata toleransı:** yanıt vermeyen veya kapanan siteler programı durdurmaz; hata yönetilir ve tarama devam eder.

### Gereksinimler

- Go 1.20+
- Arka planda açık Tor Browser (Tor ağına tünel açmak için)

### Kurulum

```bash
git clone https://github.com/KaanTuran28/CTI_Scrapper_2.git
cd CTI_Scrapper_2
go mod tidy
```

### Kullanım

1. **Tor bağlantısı:** Tor Browser'ı başlatın ve "Connect" ile ağa bağlanın. Tarayıcıyı kapatmayın, simge durumuna küçültün.
2. **Hedef belirleme:** `targets.yaml` dosyasını açın ve izlemek istediğiniz `.onion` adreslerini ekleyin:

   ```yaml
   - name: "Örnek Forum"
     url: "http://ornekadres...onion/"
   ```

3. Programı çalıştırın.

### Çıktılar

- `screenshots/` — sitelerin kanıt niteliğindeki ekran görüntüleri.
- `logs/scan_report.log` — tarama sonuçlarını içeren durum raporu.
- `logs/app.log` — teknik çalışma kayıtları.

> Yalnızca yetkili CTI araştırması ve savunma amaçlı izleme için.

### Lisans

MIT — bkz. [LICENSE](./LICENSE).
