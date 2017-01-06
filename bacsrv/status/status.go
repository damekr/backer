package status

import (
	log "github.com/Sirupsen/logrus"
	//"github.com/shirou/gopsutil/cpu"
	//"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"os"
	"runtime"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type Status struct {
	Cpu   int                    `json:"numberOfCpu"`
	Mem   *mem.VirtualMemoryStat `json:"amountOfRam"`
	Cores int                    `json:"cpuCores"`
	//Load  *load.AvgStat          `json:"load"`
}

func GetNumberOfCpu() int {
	ncpu := runtime.NumCPU()
	log.Debug("Number of CPUs: ", ncpu)
	return ncpu

}

func GetAllMemoryInformation() *mem.VirtualMemoryStat {
	ram, _ := mem.VirtualMemory()
	log.Debug("Total Memory: ", ram.Total)
	return ram
}

func GetNumberOfCores() int {
	// TODO
	return 0

}

//func (s *Status) GetServerLoad() {
//	log.Debug("Current load: ", load.AvgStat.String())

//}

func GetSeverStatus() *Status {
	status := Status{
		Cpu:   GetNumberOfCpu(),
		Mem:   GetAllMemoryInformation(),
		Cores: GetNumberOfCores(),
	}
	return &status

}