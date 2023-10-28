package dto

const url = "/api/v1/users/my-address"

type GetDetailAddressRequest struct {
	AddressId string `json:"address_id"`
	AuthorizationHeader
}

type GetDetailAddressResponse struct {
	Id                 string `json:"id"`
	ContactName        string `json:"contactName"`
	Phone              string `json:"phone"`
	DetailAddress      string `json:"detailAddress"`
	ZipCode            string `json:"zipCode"`
	CityOrProvinceId   string `json:"cityOrProvinceId"`
	CityOrProvinceName string `json:"cityOrProvinceName"`
	DistrictId         string `json:"districtId"`
	DistrictName       string `json:"districtName"`
	WardId             string `json:"wardId"`
	WardName           string `json:"wardName"`
	CountryId          int    `json:"countryId"`
	CountryName        string `json:"countryName"`
}

func (GetDetailAddressRequest) URL() string {
	return url
}
