package service

import "master/entity"

func IsPrincipalDistributor(distributorID int) bool {
	return distributorID == 0
}

func NormalizeScopeSet(region, area, distributor string) entity.EmployeeDropdownScope {
	return entity.EmployeeDropdownScope{
		RegionScope:      entity.NormalizeDropdownScope(region),
		AreaScope:        entity.NormalizeDropdownScope(area),
		DistributorScope: entity.NormalizeDropdownScope(distributor),
	}
}
