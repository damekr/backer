package status

import (
	log "github.com/sirupsen/logrus"
	//"github.com/shirou/gopsutil/cpu"
	//"github.com/shirou/gopsutil/load"
	"runtime"

	"github.com/shirou/gopsutil/mem"
)

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
