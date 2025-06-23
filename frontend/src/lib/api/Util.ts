import { API_URL } from "./Consts";

function generateFileUrl(resourceType: string, id: string): string {
    return `${API_URL}/${resourceType}/${id}/image`;
}

export { generateFileUrl };
