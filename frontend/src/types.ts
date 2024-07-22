export type StatusId = number
export type CountryId = number

export type Proxy = {
	Ip: string
	Port: string
	CountryId: CountryId
	StatusId: StatusId
	Protocol: string
	ResponseTime: number
	CreatedAt: string
	UpdatedAt: string
}

export type Country = {
	Id: CountryId
	Name: string
	Code: string
}

export type Status = {
	Id: StatusId
	Name: string
}

export type ProxyRow = Proxy & {
	CountryName: string
	CountryCode: string
	Status: string
	CreatedAtFormatted: string
	UpdatedAtFormatted: string
}
