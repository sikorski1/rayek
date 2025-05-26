import { Maps } from "@/types/main";
import { url } from "@/utils/url";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
const fetchWallMatrix = async (mapTitle: string): Promise<Float64Array> => {
	try {
		const response = await axios.get(`${url}/maps/wallmatrix/${mapTitle}`, {
			responseType: "arraybuffer",
		});
		return new Float64Array(response.data);
	} catch (error) {
		console.error(`Error fetching wall matrix`, error);
		throw new Error("Failed to fetch wall matrix");
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

export const useWallMatrix = (mapTitle: string) => {
	return useQuery({
		queryKey: ["wallMatrix", mapTitle],
		queryFn: () => fetchWallMatrix(mapTitle),
		enabled: !!mapTitle,
	});
};

export const useRayLaunching = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: startRayLaunching,
		onSuccess: () => {
			console.log("done");
		},
	});
};
