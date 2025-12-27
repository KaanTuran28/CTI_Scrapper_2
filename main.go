package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v3"
)

const (
	TorProxy    = "socks5://127.0.0.1:9150"
	LogsDir     = "logs"
	ScreensDir  = "screenshots"
	AppLog      = "logs/app.log"
	ReportLog   = "logs/scan_report.log"
	TargetsFile = "targets.yaml"
)

const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorCyan   = "\033[36m"
	ColorPurple = "\033[35m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
)

type Target struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func main() {
	setupFolders()
	logger("[INFO] === Yeni oturum başlatıldı ===")
	logger("[INFO] Program başlatıldı")

	fmt.Printf("%s i %s Tor bağlantısı kontrol ediliyor (Port: 9150)...\n", ColorCyan, ColorReset)
	logger("[INFO] Tor port kontrolü başlatıldı (Port: 9150)")

	if checkTorPort() {
		fmt.Printf("%s ✓ Tor portu açık.%s\n", ColorGreen, ColorReset)
		logger("[SUCCESS] Tor portu açık: 127.0.0.1:9150")
	} else {
		fmt.Printf("%s X Tor bağlantısı başarısız! Lütfen Tor Browser'ı açın.%s\n", ColorRed, ColorReset)
		logger("[ERROR] Tor portuna erişilemedi")
		return
	}

	fmt.Printf("%s i %s IP maskeleme kontrol ediliyor (check.torproject.org)...\n", ColorCyan, ColorReset)
	maskedIP := getTorIP()
	if maskedIP != "" {
		fmt.Printf("%s ✓ Tor Ağına Bağlısınız! Görünen IP: %s%s\n", ColorGreen, maskedIP, ColorReset)
		logger(fmt.Sprintf("[SUCCESS] IP Maskeleme Başarılı. Çıkış IP: %s", maskedIP))
	} else {
		fmt.Printf("%s ! IP adresi doğrulanamadı (Zaman aşımı). Devam ediliyor...%s\n", ColorYellow, ColorReset)
		logger("[WARNING] IP doğrulama zaman aşımına uğradı.")
	}

	targets, err := loadTargets()
	if err != nil {
		log.Fatal(err)
	}
	logger(fmt.Sprintf("[INFO] targets.yaml okundu - %d forum yüklendi", len(targets)))

	fmt.Printf("%s i %s Scan raporu oluşturuluyor...\n", ColorCyan, ColorReset)
	createInitialReport(targets, maskedIP)
	fmt.Printf("%s ✓ Scan raporu oluşturuldu: %s%s\n", ColorGreen, ReportLog, ColorReset)
	logger(fmt.Sprintf("[SUCCESS] Scan raporu oluşturuldu: %s", ReportLog))

	reader := bufio.NewReader(os.Stdin)
	for {
		printMenu(targets)
		fmt.Print("\nSelect an option: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		var selection int
		_, err := fmt.Sscanf(input, "%d", &selection)

		if err != nil {
			fmt.Println("Lütfen bir sayı girin!")
			continue
		}

		if selection == 0 {
			fmt.Println("i Exiting...")
			logger("[INFO] === Oturum sonlandırıldı ===")
			break
		} else if selection == len(targets)+1 {
			logger("[INFO] 'Scrape all forums' seçildi")
			for _, t := range targets {
				scrapeTarget(t)
			}
		} else if selection > 0 && selection <= len(targets) {
			target := targets[selection-1]
			scrapeTarget(target)
		} else {
			fmt.Println("Geçersiz seçim!")
		}
	}
}

func scrapeTarget(t Target) {
	logger(fmt.Sprintf("[INFO] Scanning forum: %s", t.Name))
	fmt.Printf("\n%s[*] %s taranıyor...%s\n", ColorCyan, t.Name, ColorReset)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(TorProxy),
		chromedp.Flag("headless", true),
		chromedp.WindowSize(1920, 1080),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 90*time.Second) 
	defer cancel()

	var buf []byte
	start := time.Now()

	err := chromedp.Run(ctx,
		chromedp.Navigate(t.URL),
		chromedp.Sleep(10*time.Second), 
		chromedp.CaptureScreenshot(&buf),
	)

	duration := time.Since(start)

	if err != nil {
		errMsg := fmt.Sprintf("%v", err)
		logger(fmt.Sprintf("[ERROR] Scraping failed for %s: %v", t.Name, err))
		updateScanReport(t.Name, "FAIL", "Erişilemedi / Zaman Aşımı") 
		fmt.Printf("%s[!] Hata: %s%s\n", ColorRed, errMsg, ColorReset)
	} else {
		logger(fmt.Sprintf("[SUCCESS] Screenshot captured for %s (%d bytes)", t.Name, len(buf)))
		
		filename := fmt.Sprintf("%s/%s_%d.png", ScreensDir, strings.ReplaceAll(t.Name, " ", "_"), time.Now().Unix())
		if err := ioutil.WriteFile(filename, buf, 0644); err == nil {
			logger(fmt.Sprintf("[SUCCESS] Screenshot saved: %s", filename))
			
			updateScanReport(t.Name, "SUCCESS", fmt.Sprintf("Aktif (%.2fs) - Kanıt: %s", duration.Seconds(), filename))
			
			fmt.Printf("%s[+] %s başarıyla kaydedildi! (%.2fs)%s\n", ColorGreen, t.Name, duration.Seconds(), ColorReset)
		}
	}
}

func getTorIP() string {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(TorProxy),
		chromedp.Flag("headless", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://check.torproject.org/"),
		chromedp.Text("body", &body),
	)
	if err != nil {
		return ""
	}
	
	if strings.Contains(body, "Congratulations") {
		return "Gizli (Tor Exit Node)"
	}
	return "Tespit Edilemedi"
}

func printMenu(targets []Target) {
	fmt.Printf("\n%s--- Dark Web Forum Scraper ---%s\n\n", ColorPurple, ColorReset)
	for i, t := range targets {
		fmt.Printf("%s%d.%s %s\n", ColorCyan, i+1, ColorReset, t.Name)
	}
	fmt.Printf("%s%d.%s Scrape all forums\n", ColorCyan, len(targets)+1, ColorReset)
	fmt.Printf("%s0.%s Exit\n", ColorRed, ColorReset)
}

func createInitialReport(targets []Target, ip string) {
	f, _ := os.Create(ReportLog)
	defer f.Close()

	header := `
--------------------------------------------------
        DARK WEB FORUM SCAN RAPORU
--------------------------------------------------
Tarih: %s
Tor Port: 9150
Çıkış IP Durumu: %s
Toplam Hedef: %d
--------------------------------------------------
HEDEFLER VE BAŞLANGIÇ DURUMLARI:
`
	fmt.Fprintf(f, header, time.Now().Format("2006-12-01 15:04:05"), ip, len(targets))

	for i, t := range targets {
		entry := fmt.Sprintf("[%d] %s\n\tURL: %s\n\tDurum: BEKLİYOR\n", i+1, t.Name, t.URL)
		f.WriteString(entry)
	}
	f.WriteString("\n--------------------------------------------------\nTARAMA SONUÇLARI :\n--------------------------------------------------\n")
}

// Tarama sonucunu rapora ekler (Append)
func updateScanReport(name, status, detail string) {
	f, err := os.OpenFile(ReportLog, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("15:04:05")
	entry := fmt.Sprintf("[%s] %s -> %s\n\tDetay: %s\n--------------------------------------------------\n", timestamp, name, status, detail)
	f.WriteString(entry)
}

func checkTorPort() bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:9150", 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func loadTargets() ([]Target, error) {
	data, err := ioutil.ReadFile(TargetsFile)
	if err != nil {
		return nil, err
	}
	var targets []Target
	err = yaml.Unmarshal(data, &targets)
	return targets, err
}

func setupFolders() {
	os.Mkdir(LogsDir, 0755)
	os.Mkdir(ScreensDir, 0755)
}

func logger(msg string) {
	f, _ := os.OpenFile(AppLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	timestamp := time.Now().Format("2006-12-01 15:04:05")
	f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, msg))
}