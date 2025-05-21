import api from "@/api/api.ts";

interface TransferRequest {
    senderWalletId: number;
    receiverWalletId: number;
    amount: number;
}

export default class TransferService {
    static async makeTransfer(payload: TransferRequest): Promise<void> {
        await api.post("/api/v1/transfers", payload);
    }
}