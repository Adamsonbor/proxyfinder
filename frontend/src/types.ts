export type StatusId = number
export type CountryId = number

export type Proxy = {
	ip: string
	port: string
	country_id: CountryId
	status_id: StatusId
	protocol: string
	response_time: number
	created_at: string
	updated_at: string
}

export type Country = {
	id: CountryId
	name: string
	code: string
}

export type Status = {
	id: StatusId
	name: string
}

export type ProxyRow = Proxy & {
	country_name: string
	country_code: string
	status: string
	created_at_formatted: string
	updated_at_formatted: string
}
