# Thor Scraper - Dark Web CTI Aracı

**Thor Scraper**, Tor ağı üzerindeki (.onion) web sitelerini otomatik olarak izlemek, erişilebilirlik durumlarını kontrol etmek ve kanıt amaçlı ekran görüntüsü almak için geliştirilmiş, **Go (Golang)** tabanlı bir Siber Tehdit İstihbaratı (CTI) aracıdır.

Bu proje, manuel istihbarat toplama süreçlerini otomatize ederek analistlere hız ve güvenlik kazandırmayı amaçlar.

## Temel Özellikler

* ** Tam Gizlilik:** Tüm ağ trafiğini yerel **Tor SOCKS5 Proxy (127.0.0.1:9150)** üzerinden geçirerek gerçek IP adresinizi gizler.
* ** Görsel Kanıt:** Ziyaret edilen sitelerin anlık ekran görüntüsünü (.png) yüksek kalitede kaydeder.
* ** Otomatik Raporlama:** Taranan hedeflerin durumunu (Başarılı/Başarısız) ve yanıt sürelerini detaylı bir log dosyasına işler.
* ** Hata Toleransı:** Yanıt vermeyen veya kapanan siteler programı durdurmaz; hata yönetilir ve tarama devam eder.

## Gereksinimler

Projenin çalışması için sisteminizde aşağıdakilerin yüklü olması gerekir:

1.  **Go (Golang):** (Sürüm 1.20 ve üzeri önerilir)
2.  **Tor Browser:** Programın Tor ağına tünel açabilmesi için arka planda açık olması gerekir.

## Kurulum

Terminali açın ve aşağıdaki komutları sırasıyla uygulayın:

git clone https://github.com/KaanTuran28/CTI_Scrapper_2.git
cd CTI_Scrapper
go mod tidy

## Kullanım
1.  **Tor Bağlantısı:** Tor Browser'ı başlatın ve "Connect" butonuna basarak ağa bağlanın. Tarayıcıyı kapatmayın, simge durumuna küçültün.

2.  **Hedef Belirleme:** targets.yaml dosyasını açın ve izlemek istediğiniz .onion adreslerini ekleyin:
   
- name: "Örnek Forum"
  url: "[http://ornekadres...onion/]"

## Çıktılar
Program tamamlandığında proje klasöründe aşağıdaki veriler oluşur:

1.  **screenshots/:** Sitelerin kanıt niteliğindeki ekran görüntüleri.

2.  **logs/scan_report.log:** Tarama sonuçlarını içeren durum raporu.

3.  **logs/app.log:** Teknik çalışma kayıtları.
