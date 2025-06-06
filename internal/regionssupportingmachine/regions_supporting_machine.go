package regionssupportingmachine

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/kyma-project/kyma-environment-broker/internal/provider"

	"github.com/kyma-project/kyma-environment-broker/internal/utils"
)

type RegionsSupportingMachine map[string]map[string][]string

func ReadRegionsSupportingMachineFromFile(filename string) (RegionsSupportingMachine, error) {
	var regionsSupportingMachineWithZones RegionsSupportingMachine
	err := utils.UnmarshalYamlFile(filename, &regionsSupportingMachineWithZones)
	if err != nil {
		return RegionsSupportingMachine{}, fmt.Errorf("while unmarshalling a file with regions supporting machine extended with zone mapping: %w", err)
	}
	return regionsSupportingMachineWithZones, nil
}

func (r RegionsSupportingMachine) IsSupported(region string, machineType string) bool {
	for machineFamily, regions := range r {
		if strings.HasPrefix(machineType, machineFamily) {
			if _, exists := regions[region]; exists {
				return true
			}
			return false
		}
	}

	return true
}

func (r RegionsSupportingMachine) SupportedRegions(machineType string) []string {
	for machineFamily, regionsMap := range r {
		if strings.HasPrefix(machineType, machineFamily) {
			regions := make([]string, 0, len(regionsMap))
			for region := range regionsMap {
				regions = append(regions, region)
			}
			sort.Strings(regions)
			return regions
		}
	}
	return []string{}
}

func (r RegionsSupportingMachine) AvailableZones(machineType, region, providerType string) ([]string, error) {
	for machineFamily, regionsMap := range r {
		if strings.HasPrefix(machineType, machineFamily) {
			zones := regionsMap[region]
			if len(zones) == 0 {
				return []string{}, nil
			}
			rand.Shuffle(len(zones), func(i, j int) { zones[i], zones[j] = zones[j], zones[i] })
			if len(zones) > 3 {
				zones = zones[:3]
			}

			var formattedZones []string
			for _, zone := range zones {
				formattedZones = append(formattedZones, provider.FullZoneName(providerType, region, zone))
			}
			return formattedZones, nil
		}
	}

	return []string{}, nil
}
