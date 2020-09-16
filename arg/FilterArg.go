package arg


type FilterArg struct {
	help    bool
	ports   []string
	ip      string
	speedup int
	scan    string
	file 	bool
}

func (f *FilterArg) File() bool {
	return f.file
}

func (f *FilterArg) SetFile(file bool) {
	f.file = file
}

func (f *FilterArg) Scan() string {
	return f.scan
}

func (f *FilterArg) SetScan(scan string) {
	f.scan = scan
}

func (f *FilterArg) Speedup() int {
	return f.speedup
}

func (f *FilterArg) SetSpeedup(speedup int) {
	f.speedup = speedup
}

func (f *FilterArg) Ip() string {
	return f.ip
}

func (f *FilterArg) SetIp(ip string) {
	f.ip = ip
}

func (f *FilterArg) Help() bool {
	return f.help
}

func (f *FilterArg) SetHelp(help bool) {
	f.help = help
}

func (f *FilterArg) Ports() []string {
	return f.ports
}

func (f *FilterArg) SetPorts(ports []string) {
	f.ports = ports
}
