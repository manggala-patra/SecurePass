package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode"
)

// ============================================================
// STRUCT
// ============================================================

type Akun struct {
	NamaLayanan    string
	Email          string
	KataSandi      string
	TanggalUpdate  time.Time
}

// ============================================================
// DATA GLOBAL
// ============================================================

var daftarAkun []Akun
var scanner = bufio.NewScanner(os.Stdin)

// ============================================================
// UTILITAS
// ============================================================

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func bacaInput(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func tekanEnterUntukLanjut() {
	fmt.Print("\n[ Tekan ENTER untuk melanjutkan... ]")
	scanner.Scan()
}

func cetakGaris(char string, panjang int) {
	fmt.Println(strings.Repeat(char, panjang))
}

func cetakJudul(judul string) {
	clearScreen()
	cetakGaris("=", 60)
	fmt.Printf("  🔐 SECUREPASS - %s\n", judul)
	cetakGaris("=", 60)
	fmt.Println()
}

// ============================================================
// KEKUATAN KATA SANDI
// ============================================================

func nilaiKekuatanKataSandi(ks string) (string, string) {
	var panjang, hurufBesar, hurufKecil, angka, simbol int
	panjang = len(ks)
	for _, c := range ks {
		switch {
		case unicode.IsUpper(c):
			hurufBesar++
		case unicode.IsLower(c):
			hurufKecil++
		case unicode.IsDigit(c):
			angka++
		default:
			simbol++
		}
	}

	skor := 0
	if panjang >= 8 {
		skor++
	}
	if panjang >= 12 {
		skor++
	}
	if hurufBesar > 0 {
		skor++
	}
	if hurufKecil > 0 {
		skor++
	}
	if angka > 0 {
		skor++
	}
	if simbol > 0 {
		skor++
	}

	switch {
	case skor <= 2:
		return "Lemah", "🔴"
	case skor <= 4:
		return "Sedang", "🟡"
	default:
		return "Kuat", "🟢"
	}
}

// ============================================================
// TAMPIL TABEL AKUN
// ============================================================

func tampilTabelAkun(data []Akun) {
	if len(data) == 0 {
		fmt.Println("  ⚠️  Tidak ada data akun.")
		return
	}
	cetakGaris("-", 80)
	fmt.Printf("  %-4s %-20s %-25s %-12s %-10s\n", "No", "Nama Layanan", "Email", "Kekuatan", "Diperbarui")
	cetakGaris("-", 80)
	for i, a := range data {
		kekuatan, ikon := nilaiKekuatanKataSandi(a.KataSandi)
		fmt.Printf("  %-4d %-20s %-25s %s %-8s %-10s\n",
			i+1,
			a.NamaLayanan,
			a.Email,
			ikon,
			kekuatan,
			a.TanggalUpdate.Format("02/01/2006"),
		)
	}
	cetakGaris("-", 80)
}

// ============================================================
// TAMBAH AKUN
// ============================================================

func tambahAkun() {
	cetakJudul("TAMBAH AKUN BARU")

	nama := bacaInput("  Nama Layanan  : ")
	if nama == "" {
		fmt.Println("\n  ❌ Nama layanan tidak boleh kosong.")
		tekanEnterUntukLanjut()
		return
	}

	// Cek duplikat
	for _, a := range daftarAkun {
		if strings.EqualFold(a.NamaLayanan, nama) {
			fmt.Println("\n  ❌ Akun dengan nama layanan tersebut sudah ada.")
			tekanEnterUntukLanjut()
			return
		}
	}

	email := bacaInput("  Email          : ")
	ks := bacaInput("  Kata Sandi     : ")

	if email == "" || ks == "" {
		fmt.Println("\n  ❌ Email dan kata sandi tidak boleh kosong.")
		tekanEnterUntukLanjut()
		return
	}

	akun := Akun{
		NamaLayanan:   nama,
		Email:         email,
		KataSandi:     ks,
		TanggalUpdate: time.Now(),
	}
	daftarAkun = append(daftarAkun, akun)

	kekuatan, ikon := nilaiKekuatanKataSandi(ks)
	fmt.Printf("\n  ✅ Akun '%s' berhasil ditambahkan!\n", nama)
	fmt.Printf("  Kekuatan kata sandi: %s %s\n", ikon, kekuatan)
	tekanEnterUntukLanjut()
}

// ============================================================
// LIHAT SEMUA AKUN
// ============================================================

func lihatSemuaAkun() {
	cetakJudul("DAFTAR SEMUA AKUN")
	tampilTabelAkun(daftarAkun)
	fmt.Printf("\n  Total akun tersimpan: %d\n", len(daftarAkun))
	tekanEnterUntukLanjut()
}

// ============================================================
// UBAH AKUN
// ============================================================

func ubahAkun() {
	cetakJudul("UBAH DATA AKUN")
	tampilTabelAkun(daftarAkun)
	if len(daftarAkun) == 0 {
		tekanEnterUntukLanjut()
		return
	}

	inputNama := bacaInput("\n  Masukkan nama layanan yang ingin diubah: ")
	idx := -1
	for i, a := range daftarAkun {
		if strings.EqualFold(a.NamaLayanan, inputNama) {
			idx = i
			break
		}
	}

	if idx == -1 {
		fmt.Println("\n  ❌ Akun tidak ditemukan.")
		tekanEnterUntukLanjut()
		return
	}

	fmt.Printf("\n  Akun ditemukan: %s (%s)\n", daftarAkun[idx].NamaLayanan, daftarAkun[idx].Email)
	fmt.Println("  (Kosongkan field jika tidak ingin mengubah)")

	emailBaru := bacaInput("  Email baru     : ")
	ksBaru := bacaInput("  Kata sandi baru: ")

	if emailBaru != "" {
		daftarAkun[idx].Email = emailBaru
	}
	if ksBaru != "" {
		daftarAkun[idx].KataSandi = ksBaru
	}
	daftarAkun[idx].TanggalUpdate = time.Now()

	fmt.Printf("\n  ✅ Akun '%s' berhasil diperbarui!\n", daftarAkun[idx].NamaLayanan)
	tekanEnterUntukLanjut()
}

// ============================================================
// HAPUS AKUN
// ============================================================

func hapusAkun() {
	cetakJudul("HAPUS AKUN")
	tampilTabelAkun(daftarAkun)
	if len(daftarAkun) == 0 {
		tekanEnterUntukLanjut()
		return
	}

	inputNama := bacaInput("\n  Masukkan nama layanan yang ingin dihapus: ")
	idx := -1
	for i, a := range daftarAkun {
		if strings.EqualFold(a.NamaLayanan, inputNama) {
			idx = i
			break
		}
	}

	if idx == -1 {
		fmt.Println("\n  ❌ Akun tidak ditemukan.")
		tekanEnterUntukLanjut()
		return
	}

	konfirmasi := bacaInput(fmt.Sprintf("\n  Yakin ingin menghapus akun '%s'? (y/t): ", daftarAkun[idx].NamaLayanan))
	if strings.ToLower(konfirmasi) != "y" {
		fmt.Println("\n  ❎ Penghapusan dibatalkan.")
		tekanEnterUntukLanjut()
		return
	}

	namaHapus := daftarAkun[idx].NamaLayanan
	daftarAkun = append(daftarAkun[:idx], daftarAkun[idx+1:]...)
	fmt.Printf("\n  ✅ Akun '%s' berhasil dihapus!\n", namaHapus)
	tekanEnterUntukLanjut()
}

// ============================================================
// PENCARIAN - SEQUENTIAL SEARCH
// ============================================================

func sequentialSearch(keyword string) []Akun {
	var hasil []Akun
	keyword = strings.ToLower(keyword)
	for _, a := range daftarAkun {
		if strings.Contains(strings.ToLower(a.NamaLayanan), keyword) {
			hasil = append(hasil, a)
		}
	}
	return hasil
}

// ============================================================
// PENCARIAN - BINARY SEARCH
// ============================================================

func binarySearch(keyword string) *Akun {
	// Binary search mensyaratkan data terurut → buat salinan terurut
	sorted := make([]Akun, len(daftarAkun))
	copy(sorted, daftarAkun)
	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].NamaLayanan) < strings.ToLower(sorted[j].NamaLayanan)
	})

	keyword = strings.ToLower(keyword)
	low, high := 0, len(sorted)-1
	for low <= high {
		mid := (low + high) / 2
		midVal := strings.ToLower(sorted[mid].NamaLayanan)
		if midVal == keyword {
			return &sorted[mid]
		} else if midVal < keyword {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return nil
}

func cariAkun() {
	cetakJudul("CARI AKUN")

	fmt.Println("  Metode pencarian:")
	fmt.Println("  1. Sequential Search (cari berdasarkan kata kunci sebagian)")
	fmt.Println("  2. Binary Search     (cari berdasarkan nama layanan tepat)")
	fmt.Println()
	pilihan := bacaInput("  Pilih metode [1/2]: ")

	switch pilihan {
	case "1":
		keyword := bacaInput("  Masukkan kata kunci nama layanan: ")
		fmt.Printf("\n  🔍 Sequential Search untuk '%s':\n\n", keyword)
		hasil := sequentialSearch(keyword)
		if len(hasil) == 0 {
			fmt.Println("  ❌ Tidak ada akun yang cocok.")
		} else {
			tampilTabelAkun(hasil)
			fmt.Printf("\n  Ditemukan %d akun.\n", len(hasil))
		}

	case "2":
		keyword := bacaInput("  Masukkan nama layanan (tepat): ")
		fmt.Printf("\n  🔍 Binary Search untuk '%s':\n\n", keyword)
		hasil := binarySearch(keyword)
		if hasil == nil {
			fmt.Println("  ❌ Akun tidak ditemukan.")
		} else {
			tampilTabelAkun([]Akun{*hasil})
			kekuatan, ikon := nilaiKekuatanKataSandi(hasil.KataSandi)
			fmt.Printf("\n  Kata sandi: %s\n", hasil.KataSandi)
			fmt.Printf("  Kekuatan  : %s %s\n", ikon, kekuatan)
		}

	default:
		fmt.Println("\n  ❌ Pilihan tidak valid.")
	}

	tekanEnterUntukLanjut()
}

// ============================================================
// PENGURUTAN - SELECTION SORT (berdasarkan nama layanan)
// ============================================================

func selectionSort() {
	n := len(daftarAkun)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if strings.ToLower(daftarAkun[j].NamaLayanan) < strings.ToLower(daftarAkun[minIdx].NamaLayanan) {
				minIdx = j
			}
		}
		daftarAkun[i], daftarAkun[minIdx] = daftarAkun[minIdx], daftarAkun[i]
	}
}

// ============================================================
// PENGURUTAN - INSERTION SORT (berdasarkan waktu input / TanggalUpdate)
// ============================================================

func insertionSort() {
	n := len(daftarAkun)
	for i := 1; i < n; i++ {
		kunci := daftarAkun[i]
		j := i - 1
		for j >= 0 && daftarAkun[j].TanggalUpdate.After(kunci.TanggalUpdate) {
			daftarAkun[j+1] = daftarAkun[j]
			j--
		}
		daftarAkun[j+1] = kunci
	}
}

func urutkanAkun() {
	cetakJudul("URUTKAN AKUN")

	fmt.Println("  Metode pengurutan:")
	fmt.Println("  1. Selection Sort — Alfabetis berdasarkan Nama Layanan")
	fmt.Println("  2. Insertion Sort — Berdasarkan Waktu Input (terlama → terbaru)")
	fmt.Println()
	pilihan := bacaInput("  Pilih metode [1/2]: ")

	switch pilihan {
	case "1":
		selectionSort()
		fmt.Println("\n  ✅ Data diurutkan secara alfabetis (Selection Sort).\n")
		tampilTabelAkun(daftarAkun)
	case "2":
		insertionSort()
		fmt.Println("\n  ✅ Data diurutkan berdasarkan waktu input (Insertion Sort).\n")
		tampilTabelAkun(daftarAkun)
	default:
		fmt.Println("\n  ❌ Pilihan tidak valid.")
	}

	tekanEnterUntukLanjut()
}

// ============================================================
// STATISTIK
// ============================================================

func tampilStatistik() {
	cetakJudul("STATISTIK AKUN")

	total := len(daftarAkun)
	fmt.Printf("  Total akun tersimpan  : %d\n\n", total)

	if total == 0 {
		fmt.Println("  Belum ada data akun.")
		tekanEnterUntukLanjut()
		return
	}

	lemah, sedang, kuat := 0, 0, 0
	for _, a := range daftarAkun {
		kekuatan, _ := nilaiKekuatanKataSandi(a.KataSandi)
		switch kekuatan {
		case "Lemah":
			lemah++
		case "Sedang":
			sedang++
		case "Kuat":
			kuat++
		}
	}

	cetakGaris("-", 40)
	fmt.Println("  Klasifikasi Kekuatan Kata Sandi:")
	cetakGaris("-", 40)
	fmt.Printf("  🔴 Lemah  : %d akun\n", lemah)
	fmt.Printf("  🟡 Sedang : %d akun\n", sedang)
	fmt.Printf("  🟢 Kuat   : %d akun\n", kuat)
	cetakGaris("-", 40)

	// Bar chart sederhana
	fmt.Println("\n  Distribusi Kekuatan:")
	if total > 0 {
		barLemah := int(float64(lemah) / float64(total) * 30)
		barSedang := int(float64(sedang) / float64(total) * 30)
		barKuat := int(float64(kuat) / float64(total) * 30)
		fmt.Printf("  Lemah  [%-30s] %d%%\n", strings.Repeat("█", barLemah), lemah*100/total)
		fmt.Printf("  Sedang [%-30s] %d%%\n", strings.Repeat("█", barSedang), sedang*100/total)
		fmt.Printf("  Kuat   [%-30s] %d%%\n", strings.Repeat("█", barKuat), kuat*100/total)
	}

	tekanEnterUntukLanjut()
}

// ============================================================
// LIHAT DETAIL KATA SANDI
// ============================================================

func lihatKataSandi() {
	cetakJudul("LIHAT KATA SANDI")
	tampilTabelAkun(daftarAkun)
	if len(daftarAkun) == 0 {
		tekanEnterUntukLanjut()
		return
	}

	inputNama := bacaInput("\n  Masukkan nama layanan untuk melihat kata sandi: ")
	for _, a := range daftarAkun {
		if strings.EqualFold(a.NamaLayanan, inputNama) {
			fmt.Printf("\n  Layanan    : %s\n", a.NamaLayanan)
			fmt.Printf("  Email      : %s\n", a.Email)
			fmt.Printf("  Kata Sandi : %s\n", a.KataSandi)
			kekuatan, ikon := nilaiKekuatanKataSandi(a.KataSandi)
			fmt.Printf("  Kekuatan   : %s %s\n", ikon, kekuatan)
			fmt.Printf("  Diperbarui : %s\n", a.TanggalUpdate.Format("02 January 2006 15:04:05"))
			tekanEnterUntukLanjut()
			return
		}
	}
	fmt.Println("\n  ❌ Akun tidak ditemukan.")
	tekanEnterUntukLanjut()
}

// ============================================================
// DATA CONTOH
// ============================================================

func muatDataContoh() {
	waktuBase := time.Now()
	daftarAkun = []Akun{
		{"Google", "user@gmail.com", "G00gl3@Secure!", waktuBase.Add(-72 * time.Hour)},
		{"Instagram", "user@gmail.com", "insta123", waktuBase.Add(-48 * time.Hour)},
		{"Github", "dev@email.com", "Dev#Code$2024", waktuBase.Add(-24 * time.Hour)},
		{"Netflix", "user@gmail.com", "nflx456", waktuBase.Add(-12 * time.Hour)},
		{"Tokopedia", "user@gmail.com", "Toko!Pedia#99", waktuBase.Add(-6 * time.Hour)},
	}
}

// ============================================================
// MENU UTAMA
// ============================================================

func menuUtama() {
	for {
		cetakJudul("PENGELOLA KATA SANDI PRIBADI")
		fmt.Printf("  Total Akun: %d\n\n", len(daftarAkun))
		cetakGaris("-", 40)
		fmt.Println("  MENU UTAMA")
		cetakGaris("-", 40)
		fmt.Println("  1. Tambah Akun Baru")
		fmt.Println("  2. Lihat Semua Akun")
		fmt.Println("  3. Ubah Data Akun")
		fmt.Println("  4. Hapus Akun")
		fmt.Println("  5. Cari Akun")
		fmt.Println("  6. Urutkan Akun")
		fmt.Println("  7. Lihat Kata Sandi")
		fmt.Println("  8. Statistik")
		fmt.Println("  0. Keluar")
		cetakGaris("-", 40)

		pilihan := bacaInput("  Pilihan Anda: ")

		switch pilihan {
		case "1":
			tambahAkun()
		case "2":
			lihatSemuaAkun()
		case "3":
			ubahAkun()
		case "4":
			hapusAkun()
		case "5":
			cariAkun()
		case "6":
			urutkanAkun()
		case "7":
			lihatKataSandi()
		case "8":
			tampilStatistik()
		case "0":
			cetakJudul("TERIMA KASIH")
			fmt.Println("  Sampai jumpa! Data Anda aman bersama SecurePass. 🔐")
			fmt.Println()
			os.Exit(0)
		default:
			fmt.Println("\n  ❌ Pilihan tidak valid. Silakan coba lagi.")
			tekanEnterUntukLanjut()
		}
	}
}

// ============================================================
// MAIN
// ============================================================

func main() {
	muatDataContoh()
	menuUtama()
}