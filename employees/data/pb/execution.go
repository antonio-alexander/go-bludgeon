package pb

import "github.com/antonio-alexander/go-bludgeon/employees/data"

func FromEmployeePartial(e *data.EmployeePartial) *EmployeePartial {
	if e == nil {
		return nil
	}
	employeePartial := &EmployeePartial{}
	if e.FirstName != nil {
		employeePartial.FirstNameOneof = &EmployeePartial_FirstName{
			FirstName: *e.FirstName,
		}
	}
	if e.LastName != nil {
		employeePartial.LastNameOneof = &EmployeePartial_LastName{
			LastName: *e.LastName,
		}
	}
	if e.EmailAddress != nil {
		employeePartial.EmailAddressOneof = &EmployeePartial_EmailAddress{
			EmailAddress: *e.EmailAddress,
		}
	}
	return employeePartial
}

func ToEmployeePartial(e *EmployeePartial) *data.EmployeePartial {
	if e == nil {
		return nil
	}
	employeePartial := &data.EmployeePartial{}
	if e.FirstNameOneof != nil {
		s := e.GetFirstName()
		employeePartial.FirstName = &s
	}
	if e.LastNameOneof != nil {
		s := e.GetLastName()
		employeePartial.LastName = &s
	}
	if e.EmailAddressOneof != nil {
		s := e.GetEmailAddress()
		employeePartial.EmailAddress = &s
	}
	return employeePartial
}

func FromEmployee(e *data.Employee) *Employee {
	if e == nil {
		return nil
	}
	return &Employee{
		Id:            e.ID,
		FirstName:     e.FirstName,
		LastName:      e.LastName,
		EmailAddress:  e.EmailAddress,
		LastUdpated:   e.LastUpdated,
		LastUpdatedBy: e.LastUpdatedBy,
		Version:       int32(e.Version),
	}
}

func ToEmployee(e *Employee) *data.Employee {
	if e == nil {
		return nil
	}
	return &data.Employee{
		ID:            e.GetId(),
		FirstName:     e.GetFirstName(),
		LastName:      e.GetLastName(),
		EmailAddress:  e.GetEmailAddress(),
		LastUpdated:   e.GetLastUdpated(),
		LastUpdatedBy: e.GetLastUpdatedBy(),
		Version:       int(e.GetVersion()),
	}
}

func FromEmployees(e []*data.Employee) []*Employee {
	var employees []*Employee
	for _, e := range e {
		employees = append(employees, FromEmployee(e))
	}
	return employees
}

func ToEmployees(e []*Employee) []*data.Employee {
	var employees []*data.Employee
	for _, e := range e {
		employees = append(employees, ToEmployee(e))
	}
	return employees
}

func FromEmployeeSearch(e *EmployeeSearch) *data.EmployeeSearch {
	if e == nil {
		return nil
	}
	employeeSearch := &data.EmployeeSearch{
		IDs:            e.Ids,
		FirstNames:     e.FirstNames,
		LastNames:      e.LastNames,
		EmailAddresses: e.EmailAddresses,
	}
	if e.FirstNameOneof != nil {
		s := e.GetFirstName()
		employeeSearch.FirstName = &s
	}
	if e.LastNameOneof != nil {
		s := e.GetLastName()
		employeeSearch.LastName = &s
	}
	if e.EmailAddressOneof != nil {
		s := e.GetEmailAddress()
		employeeSearch.EmailAddress = &s
	}
	return employeeSearch
}

func ToEmployeeSearch(e *data.EmployeeSearch) *EmployeeSearch {
	if e == nil {
		return nil
	}
	employeeSearch := &EmployeeSearch{
		Ids:            e.IDs,
		FirstNames:     e.FirstNames,
		LastNames:      e.LastNames,
		EmailAddresses: e.EmailAddresses,
	}
	if e.FirstName != nil {
		employeeSearch.FirstNameOneof = &EmployeeSearch_FirstName{
			FirstName: *e.FirstName,
		}
	}
	if e.LastName != nil {
		employeeSearch.LastNameOneof = &EmployeeSearch_LastName{
			LastName: *e.LastName,
		}
	}
	if e.EmailAddress != nil {
		employeeSearch.EmailAddressOneof = &EmployeeSearch_EmailAddress{
			EmailAddress: *e.EmailAddress,
		}
	}
	return employeeSearch
}
