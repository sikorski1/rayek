import { Maps } from "@/types/main";
import { url } from "@/utils/url";
import { useMutation, useQuery } from "@tanstack/react-query";
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
const fetchMapById = async (mapTitle: string) => {
	try {
		const response = await axios.get(`${url}/maps/${mapTitle}`);
		return response.data;
	} catch (error) {
		console.error(`Error fetching maps`, error);
		throw new Error(`Failed to fetch maps`);
	}
};

const startRayLaunching = async ({ mapTitle, configData }: { mapTitle: string; configData: any }) => {
	console.log(configData);
	try {
		const response = await axios.post(
			url + `/maps/rayLaunch/${mapTitle}`,
			{ ...configData },
			{
				headers: {
					"Content-Type": "application/json",
				},
			}
		);
		return response.data;
	} catch (error) {
		console.error(`Error fetching wall matrix`, error);
		throw new Error("Failed to fetch wall matrix");
	}
};
export const useGetMaps = () => {
	return useQuery<Maps[], Error>({ queryKey: ["maps"], queryFn: fetchMaps });
};

export const useGetMapById = (mapTitle: string) => {
	return useQuery({ queryKey: ["map", mapTitle], queryFn: () => fetchMapById(mapTitle) });
};


export const useRayLaunching = (handleOnSuccess:(data:any) => void) => {

	return useMutation({
		mutationFn: startRayLaunching,
		onSuccess: (data) => {
			handleOnSuccess(data)
		},
	});
};
