import api from "../api";
import {AxiosResponse} from "axios";

export interface Wallet {
    id: number;
    userId: number;
    number: string;
    balance: number;
    currency: string;
    status: string;
    createdAt: string;
}

export default class WalletService {
    /**
     * Получение информации о кошельке пользователя
     */
    static async getUserWallets(): Promise<AxiosResponse<Wallet[]>> {
        return api.get<Wallet[]>('/api/v1/wallet');
    }
    
    /**
     * Проверка наличия и создание кошелька, если его нет
     */
    static async initWalletIfNeeded(): Promise<Wallet> {
        const response = await api.post<Wallet>('/api/v1/wallet/init');
        return response.data;
    }
} 