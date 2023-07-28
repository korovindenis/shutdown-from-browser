package usecase

import (
	"log"
	"os"
)

var syscall Syscall

func Init() {
	syscall = Syscall{}
}

type Syscall struct {
	LINUX_REBOOT_CMD_POWER_OFF byte
	LINUX_REBOOT_CMD_RESTART   byte
}

func (s Syscall) Reboot(mode byte) error {
	log.Println("A mode '", mode, "' would have occurred")

	os.Exit(0)

	return nil
}
