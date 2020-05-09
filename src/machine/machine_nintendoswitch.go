// +build nintendoswitch

package machine

type PinMode uint8

// Dummy machine package that calls out to external functions.

var (
	SPI0  = SPI{0}
	I2C0  = I2C{0}
	UART0 = UART{0}
)

const (
	PinInput PinMode = iota
	PinOutput
	PinInputPullup
	PinInputPulldown
)

func (p Pin) Configure(config PinConfig) {
	gpioConfigure(p, config)
}

func (p Pin) Set(value bool) {
	gpioSet(p, value)
}

func (p Pin) Get() bool {
	return gpioGet(p)
}

func gpioConfigure(pin Pin, config PinConfig) {

}

func gpioSet(pin Pin, value bool) {

}

func gpioGet(pin Pin) bool {
	return false
}

type SPI struct {
	Bus uint8
}

type SPIConfig struct {
	Frequency uint32
	SCK       Pin
	MOSI      Pin
	MISO      Pin
	Mode      uint8
}

func (spi SPI) Configure(config SPIConfig) {

}

func (spi SPI) Transfer(w byte) (byte, error) {
	return 0, nil
}

func spiConfigure(bus uint8, sck Pin, mosi Pin, miso Pin) {

}

func spiTransfer(bus uint8, w uint8) uint8 {
	return 0
}

// InitADC enables support for ADC peripherals.
func InitADC() {
	// Nothing to do here.
}

// Configure configures an ADC pin to be able to be used to read data.
func (adc ADC) Configure() {
}

// Get reads the current analog value from this ADC peripheral.
func (adc ADC) Get() uint16 {
	return adcRead(adc.Pin)
}

func adcRead(pin Pin) uint16 {
	return 0
}

// InitPWM enables support for PWM peripherals.
func InitPWM() {
	// Nothing to do here.
}

// Configure configures a PWM pin for output.
func (pwm PWM) Configure() {
}

// Set turns on the duty cycle for a PWM pin using the provided value.
func (pwm PWM) Set(value uint16) {

}

func pwmSet(pin Pin, value uint16) {

}

// I2C is a generic implementation of the Inter-IC communication protocol.
type I2C struct {
	Bus uint8
}

// I2CConfig is used to store config info for I2C.
type I2CConfig struct {
	Frequency uint32
	SCL       Pin
	SDA       Pin
}

// Configure is intended to setup the I2C interface.
func (i2c I2C) Configure(config I2CConfig) {

}

// Tx does a single I2C transaction at the specified address.
func (i2c I2C) Tx(addr uint16, w, r []byte) error {
	return nil
}

func i2cConfigure(bus uint8, scl Pin, sda Pin) {

}

func i2cTransfer(bus uint8, w *byte, wlen int, r *byte, rlen int) int {
	return 0
}

type UART struct {
	Bus uint8
}

type UARTConfig struct {
	BaudRate uint32
	TX       Pin
	RX       Pin
}

// Configure the UART.
func (uart UART) Configure(config UARTConfig) {
	uartConfigure(uart.Bus, config.TX, config.RX)
}

// Read from the UART.
func (uart UART) Read(data []byte) (n int, err error) {
	return uartRead(uart.Bus, &data[0], len(data)), nil
}

// Write to the UART.
func (uart UART) Write(data []byte) (n int, err error) {
	return uartWrite(uart.Bus, &data[0], len(data)), nil
}

// Buffered returns the number of bytes currently stored in the RX buffer.
func (uart UART) Buffered() int {
	return 0
}

// ReadByte reads a single byte from the UART.
func (uart UART) ReadByte() (byte, error) {
	return 0, nil
}

// WriteByte writes a single byte to the UART.
func (uart UART) WriteByte(b byte) error {
	return nil
}

func uartConfigure(bus uint8, tx Pin, rx Pin) {}

func uartRead(bus uint8, buf *byte, bufLen int) int { return 0 }

func uartWrite(bus uint8, buf *byte, bufLen int) int { return 0 }
