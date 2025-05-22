// services/transferService.ts
import api from "@/api/api.ts";

export interface Transfer {
    id: number;
    amount: number;
    currency: string;
    toUsername?: string;
    fromUsername?: string;
    createdAt: string;
    senderWalletId: number;
    receiverWalletId: number;
}

interface TransferRequest {
    senderWalletId: number;
    receiverWalletId: number;
    amount: number;
}

export default class TransferService {
    static async makeTransfer(payload: TransferRequest): Promise<void> {
        await api.post("/api/v1/transfers", payload);
    }

    static async getTransferHistory(walletId: number): Promise<Transfer[]> {
        const response = await api.get(`/api/v1/transfers/wallet/${walletId}`);
        return response.data;
    }

    static async getReceiverUsername(walletId: number): Promise<{username: string, avatarUrl: string}> {
        const res = await api.get(`/api/v1/transfers/username/${walletId}`);
        return {username: res.data.username, avatarUrl: res.data.avatarUrl};
    }
}
