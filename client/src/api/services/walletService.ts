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
    static async getUserWallets(): Promise<AxiosResponse<Wallet[]>> {
        return api.get<Wallet[]>('/api/v1/wallet');
    }

    static async initWalletIfNeeded(): Promise<Wallet> {
        const response = await api.post<Wallet>('/api/v1/wallet/init');
        return response.data;
    }

    static async topUpWallet(currency: string, amount: number): Promise<Wallet> {
        const response = await api.post<Wallet>('/api/v1/wallet/topup', {
            currency,
            amount
        });
        return response.data;
    }
}
