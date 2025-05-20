import { Maps } from "@/types/main";
import { url } from "@/utils/url";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
const fetchMaps = async () => {
	try {
		const response = await axios.get(`${url}/maps/`);
		return response.data;
	} catch (error) {
		console.error(`Error fetching maps`, error);
		throw new Error(`Failed to fetch maps`);
	}
};

export const useGetMaps = () => {
	return useQuery<Maps, Error>({ queryKey: ["maps"], queryFn: fetchMaps });
};
