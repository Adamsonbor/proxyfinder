import { useState } from "react"

export class QueryBuilder {
	private _filters: object = {}
	private _sorts: object = {}

	constructor() { }

	public getFilter(key: string): any {
		return this._filters[key]
	}

	public getSort(key: string): any {
		return this._sorts[key]
	}

	public setFilter(key: string, value: any) {
		this._filters[key] = value
	}

	public setFilters(filters: object) {
		this._filters = filters
	}

	public setSort(sort: string, order: string) {
		this._sorts[sort] = order
	}

	public setSorts(sorts: object) {
		this._sorts = sorts
	}

	public delFilter(key: string) {
		delete this._filters[key]
	}

	public delSort(key: string) {
		delete this._sorts[key]
	}

	public clearFilters() {
		this._filters = {}
	}

	public clearSorts() {
		this._sorts = {}
	}

	private buildFilters(): string {
		let query = ""
		if (Object.keys(this._filters).length === 0) {
			return query
		}
		for (let key in this._filters) {
			query += `${key}=${this._filters[key]}&`
		}
		return query
	}

	private buildSorts(): string {
		let query = ""
		if (Object.keys(this._sorts).length === 0) {
			return query
		}

		let sorts = ""
		let orders = ""

		let i = 0
		for (let key in this._sorts) {
			if (i > 0) {
				sorts += ","
				orders += ","
			}
			sorts += `${key}`
			orders += `${this._sorts[key]}`
			i++
		}

		query += `sort_by=${sorts}&sort_order=${orders}`
		return query
	}

	public build(): string {
		if (Object.keys(this._filters).length === 0 && Object.keys(this._sorts).length === 0) {
			return ""
		}
		let query = "?"
		query += this.buildFilters()
		query += this.buildSorts()

		return query
	}
}

export function useQuery() {
	const [queryBuilder, _] = useState<QueryBuilder>(new QueryBuilder())

	return queryBuilder
}
