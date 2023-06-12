package service

import (
	"fmt"
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
	fmt.Println("A ", mode, " would have occurred")

	return nil
}
