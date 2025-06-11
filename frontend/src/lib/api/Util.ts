import { API_URL } from "./consts";

function generateFileUrl(collection: string, id: string, filename: string): string {
    return `${API_URL}/api/files/${collection}/${id}/${filename}`;
}

export { generateFileUrl };
