export interface Favorits {
	id: number
	user_id: number
	proxy_id: number
	created_at: string
	updated_at: string
}

export interface Proxy {
	id: number
	ip: string
	port: string
	protocol: string
	response_time: number
	created_at: string
	updated_at: string
	status_id: number
	country_id: number
	status: Status
	country: Country
}

export type Country = {
	id: number
	name: string
	code: string
}

export type Status = {
	id: number
	name: string
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
