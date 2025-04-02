// gracefuler предоставляет фнкционал для реализации graceful shutdown.
package graceful

import (
	"os"
	"os/signal"
)

var interrupt chan os.Signal

func init() {
	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
}

// WaitShutdown блокирует текущую горутину до получения сигнала [os.Interrupt].
func WaitShutdown() {
	<-interrupt
}

// Shutdown записывает сигнал [os.Interrupt] во внутренний канал interrupt.
func Shutdown() {
	interrupt <- os.Interrupt
}
