export type StatusId = number
export type CountryId = number

export interface Favorits {
	id: number
	user_id: number
	proxy_id: number
	created_at: string
	updated_at: string
}

export interface ProxyV2 {
	id: number
	ip: string
	port: string
	protocol: string
	response_time: number
	created_at: string
	updated_at: string
	status: Status
	country: Country
}

export type Proxy = {
	id: number
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
	created_at_formatted: Date
	updated_at_formatted: Date
}

export type User = {
	id: number
	name: string
	email: string
	phone: string
	photo_url: string
	date_of_birth: string
	created_at: string
	updated_at: string
}
