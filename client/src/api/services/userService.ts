import { AxiosResponse } from "axios";
import api from "@/api/api.ts";

export interface User {
    id: number;
    email: string | null;
    username: string;
    telegramId: number | null;
    telegramUsername: string | null;
    telegramFirstName: string | null;
    avatarUrl: string | null;
    createdAt: string;
    updatedAt: string;
}

export default class UserService {
    static async getCurrentUser(): Promise<User> {
        const response: AxiosResponse<User> = await api.get("/api/v1/users/me");
        return response.data;
    }

    static async getById(id: number): Promise<User> {
        const response: AxiosResponse<User> = await api.get(`/api/v1/users/id/${id}`);
        return response.data;
    }

    static async getByEmail(email: string): Promise<User> {
        const response: AxiosResponse<User> = await api.get(`/api/v1/users/email/${encodeURIComponent(email)}`);
        return response.data;
    }

    static async getByTelegramId(telegramId: number): Promise<User> {
        const response: AxiosResponse<User> = await api.get(`/api/v1/users/telegram/${telegramId}`);
        return response.data;
    }

    static async getByUsername(username: string): Promise<User> {
        const response: AxiosResponse<User> = await api.get(`/api/v1/users/username/${username}`);
        return response.data;
    }

    static async getByWalletNumber(walletNumber: string): Promise<User> {
        const response: AxiosResponse<User> = await api.get(`/api/v1/users/by-wallet/${walletNumber}`);
        return response.data;
    }

}
