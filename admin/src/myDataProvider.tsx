import simpleRestProvider from "ra-data-simple-rest";

const proxyfinderDataProvider = simpleRestProvider(import.meta.env.VITE_API_SERVER_URL)

proxyfinderDataProvider.getList = async (resource: string, params: any) => {
	const { page, perPage } = params.pagination;
	const url = `
		${import.meta.env.VITE_API_URL}/api/v2/${resource}?\
		page=${page}&\
		perPage=${perPage}`
	return fetch(url)
		.then(res => res.json())
		.then(jsonData => { console.log(jsonData); jsonData["total"] = jsonData.data.length; return jsonData })
		.catch(error => console.log(error));
};

export { proxyfinderDataProvider };
