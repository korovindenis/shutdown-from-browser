package usecase

import (
	"context"
	"fmt"
	_ "syscall"
	"time"

	"github.com/korovindenis/shutdown-from-browser/v2/internal/domain/entity"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate mockery --name=storage
type storage interface {
	SetPoff(pc entity.MyPc) error
	GetTimePoff() (string, error)
	GetModePoff() (string, error)
}

type ComputerUsecase struct {
	computerStorage storage
	logger          *zap.Logger
}

func New(cStorage storage, log *zap.Logger) *ComputerUsecase {
	return &ComputerUsecase{
		computerStorage: cStorage,
		logger:          log,
	}
}

func (cu *ComputerUsecase) GetTimePowerOff() (string, error) {
	time, err := cu.computerStorage.GetTimePoff()
	if err != nil {
		return "", errors.Wrap(err, "GetTimePowerOff err")
	}

	return time, nil
}

func (cu *ComputerUsecase) GetModePowerOff() (string, error) {
	time, err := cu.computerStorage.GetModePoff()
	if err != nil {
		return "", errors.Wrap(err, "GetModePowerOff err")
	}

	return time, nil
}

func (cu *ComputerUsecase) SetPowerOff(pc entity.MyPc) error {
	if err := validator(pc); err != nil {
		return errors.Wrap(err, "validator SetPowerOff")
	}

	if err := cu.computerStorage.SetPoff(pc); err != nil {
		return errors.Wrap(err, "SetPoff SetPowerOff")
	}
	return nil
}

func (cu *ComputerUsecase) IsNeedPowerOff(ctx context.Context, logslevel uint8) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// TODO: ERR!
			modePoff, _ := cu.GetModePowerOff()
			timePoff, _ := cu.GetTimePowerOff()

			time.Sleep(5 * time.Second)
			v, _ := time.Parse(time.RFC3339, timePoff)
			timeRemaining := getTimeRemaining(v)

			if modePoff != "" {
				if timeRemaining.Total <= 0 {
					if logslevel > 0 {
						cu.logger.Info("Run: " + modePoff)
					}
					callMode := syscall.LINUX_REBOOT_CMD_POWER_OFF
					if modePoff == "reboot" {
						callMode = syscall.LINUX_REBOOT_CMD_RESTART
					}
					// BYE
					if err := syscall.Reboot(callMode); err != nil {
						panic(err)
					}
				}
				if logslevel > 1 {
					cu.logger.Info(fmt.Sprintf("Time for %s - %d:%d:%d\n", modePoff, timeRemaining.Hours, timeRemaining.Minutes, timeRemaining.Seconds))
				}
			}
		}
	}
}

func getTimeRemaining(t time.Time) entity.Countdown {
	currentTime := time.Now().UTC()
	difference := t.Sub(currentTime)

	total := int(difference.Seconds())
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60
	seconds := int(total % 60)

	return entity.Countdown{
		Total:   total,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
}

func validator(pc entity.MyPc) error {
	// validate input
	if pc.ModePowerOff != "" && pc.ModePowerOff != "shutdown" && pc.ModePowerOff != "reboot" {
		return errors.New("'mode' error validator")
	}

	// validate time (format and time no more 24h.)
	const customDateTimeFormat = "2006-01-02T15:04:05.000Z"
	if _, err := time.Parse(customDateTimeFormat, pc.TimePowerOff); err != nil {
		return errors.Wrap(err, "'timestamp' error validator")
	}

	return nil
}
