import {AxiosResponse} from "axios";
import api from "@/api/api.ts";

export interface User {
    id: number;
    email: string | null;
    username: string;
    telegramId: number | null;
    telegramUsername: string | null;
    telegramFirstName: string | null;
    createdAt: string;
    updatedAt: string;
}

export default class UserService {
    static async getCurrentUser(): Promise<User> {
        const response: AxiosResponse<User> = await api.get("/api/v1/users/me");
        return response.data;
    }
}
