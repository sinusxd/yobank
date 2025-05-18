import dollarLogo from './icons/dollar.png';
import euroLogo from './icons/euro.png';
import rubleLogo from './icons/ruble.png';
import yuanLogo from './icons/yuan.png';
import {Wallet} from "@/api/services/walletService.ts";

export function mapCurrencyToSymbol(currency: string): string {
    const map: Record<string, string> = {
        RUB: '₽',
        EUR: '€',
        USD: '$',
        CNY: '¥'
    }
    return map[currency] || currency
}

export function parseBalance(balance: number): number {
    return balance / 100
}

export function mapCurrencyToName(currency: string): string {
    const map: Record<string, string> = {
        RUB: 'Рубли',
        EUR: 'Евро',
        USD: 'Доллары',
        CNY: 'Юани'
    }
    return map[currency] || currency
}

export function mapCurrencyToLogo(currency: string): string {
    const map: Record<string, string> = {
        RUB: rubleLogo,
        EUR: euroLogo,
        USD: dollarLogo,
        CNY: yuanLogo
    }
    return map[currency] || ''
}


// Форматирование баланса
export const formatBalance = (balance?: number) => {
    if (balance === undefined || balance === null) return '0.00';
    return (balance / 100).toFixed(2);
};

export const convertToRub = (balance: number, rates: Record<string, number>, currency: string) => {
    if (!rates[currency] || rates['RUB'] === 0) return '-';
    return ((balance / 100) * rates[currency]).toFixed(2);
};

export function sumAllWalletsInRub(wallets: Wallet[] | null, rates: Record<string, number>): number {
    if (!wallets) return 0;
    return wallets.reduce((sum, wallet) => {
        const rate = rates[wallet.currency] || 1; // для RUB = 1
        return sum + (wallet.balance / 100) * rate;
    }, 0);
}




