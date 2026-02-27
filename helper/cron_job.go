// helper/cron_job.go

package helper

import (
	"backend_go/internal/model"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type VoucherCronJob struct {
	DB *gorm.DB
}

func NewVoucherCronJob(db *gorm.DB) *VoucherCronJob {
	return &VoucherCronJob{DB: db}
}

// AutoExpireVouchers updates voucher status to "expired"
func (c *VoucherCronJob) AutoExpireVouchers() (int64, error) {
	now := time.Now()

	// =====================================================
	// PERBAIKAN: Gunakan waktu saat ini, bukan 00:00:00
	// =====================================================
	// Ini akan mengambil semua voucher yang:
	// - status = 'used'
	// - valid_until <= sekarang (sudah expired)

	fmt.Printf("[CRON] Checking vouchers expiring at: %s\n", now.Format("2006-01-02 15:04:05"))

	result := c.DB.Model(&model.Voucher{}).
		Where("status = ?", "used").
		Where("valid_until <= ?", now). // ← GANTI: sekarang, bukan today 00:00
		Update("status", "expired")

	if result.Error != nil {
		fmt.Printf("[CRON] Error: %v\n", result.Error)
		return 0, result.Error
	}

	affected := result.RowsAffected
	fmt.Printf("[CRON] Expired %d vouchers\n", affected)

	return affected, nil
}

// StartVoucherCron -versi testing (jalan setiap 1 menit)
func (c *VoucherCronJob) StartVoucherCronForTesting() *cron.Cron {
	scheduler := cron.New()

	// Setiap 1 menit - BAGUS UNTUK TESTING
	_, err := scheduler.AddFunc("@every 1m", func() {
		fmt.Println("[CRON] Running auto-expire vouchers (test)...")
		affected, _ := c.AutoExpireVouchers()
		if affected > 0 {
			fmt.Printf("[CRON] Expired %d vouchers\n", affected)
		}
	})

	if err != nil {
		fmt.Printf("[CRON] Failed to add cron job: %v\n", err)
		return nil
	}

	scheduler.Start()
	fmt.Println("[CRON] Voucher auto-expire started (every 1 minute for testing)")

	return scheduler
}

// StartVoucherCron - versi production (jalan setiap jam)
func (c *VoucherCronJob) StartVoucherCron() *cron.Cron {
	scheduler := cron.New()

	// Setiap jam pada menit ke-0 (00:00, 01:00, 02:00, dll)
	_, err := scheduler.AddFunc("0 0 * * *", func() {
		fmt.Println("[CRON] Running auto-expire vouchers...")
		affected, _ := c.AutoExpireVouchers()
		if affected > 0 {
			fmt.Printf("[CRON] Expired %d vouchers\n", affected)
		}
	})

	if err != nil {
		fmt.Printf("[CRON] Failed to add cron job: %v\n", err)
		return nil
	}

	scheduler.Start()
	fmt.Println("[CRON] Voucher auto-expire started (every hour)")

	return scheduler
}

func (c *VoucherCronJob) StopVoucherCron(scheduler *cron.Cron) {
	if scheduler != nil {
		scheduler.Stop()
		fmt.Println("[CRON] Scheduler stopped")
	}
}
