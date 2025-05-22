import {FC, useEffect, useState} from "react";
import {Page} from "@/components/Page";
import {Button, Input, List, Section, Select, Title} from "@telegram-apps/telegram-ui";
import QRCode from "react-qr-code";
import WalletService, {Wallet} from "@/api/services/walletService";
import {mapCurrencyToSymbol} from "@/utils/currency";

export const AddMoneyQrPage: FC = () => {
    const [wallets, setWallets] = useState<Wallet[]>([]);
    const [selectedWalletId, setSelectedWalletId] = useState<number | null>(null);
    const [amount, setAmount] = useState<string>("");
    const [qrValue, setQrValue] = useState("");

    const selectedWallet = wallets.find(w => w.id === selectedWalletId);

    useEffect(() => {
        const load = async () => {
            const {data} = await WalletService.getUserWallets();
            setWallets(data);
            const rubWallet = data.find(w => w.currency === "RUB");
            if (rubWallet) setSelectedWalletId(rubWallet.id);
        };
        load();
    }, []);

    const generateQR = () => {
        if (!selectedWallet || !amount) return;
        const value = `yobank://transfer?method=card&target=${selectedWallet.number}&amount=${amount}`;
        setQrValue(value);
    };

    return (
        <Page back>
            <List>
                <Title level="2" weight="2" style={{textAlign: "center", margin: "24px 0"}}>
                    QR для пополнения
                </Title>

                <Select
                    value={selectedWalletId?.toString() || ""}
                    onChange={(e) => setSelectedWalletId(Number(e.target.value))}
                    style={{marginBottom: 12}}
                >
                    <option disabled value="">Выберите счёт</option>
                    {wallets.map(wallet => (
                        <option key={wallet.id} value={wallet.id}>
                            {wallet.currency} ({wallet.number})
                        </option>
                    ))}
                </Select>

                <Input
                    placeholder={`Введите сумму (${mapCurrencyToSymbol(selectedWallet?.currency || "")})`}
                    type="text"
                    inputMode="decimal"
                    value={amount ? `${amount} ${mapCurrencyToSymbol(selectedWallet?.currency || "")}` : ""}
                    onChange={(e) => {
                        const raw = e.target.value.replace(/[^\d.,]/g, "");
                        const normalized = raw.replace(",", ".");
                        if (/^\d*([.]?\d{0,2})?$/.test(normalized)) {
                            setAmount(normalized);
                            requestAnimationFrame(() => {
                                const el = e.target as HTMLInputElement;
                                const pos = normalized.length;
                                el.setSelectionRange(pos, pos);
                            });
                        }
                    }}
                />

                <Button
                    style={{marginTop: 12}}
                    onClick={generateQR}
                    stretched
                    size="l"
                    disabled={!selectedWallet || !amount}
                >
                    Сгенерировать QR
                </Button>

                {qrValue && (
                    <Section
                        header={'Сканируйте QR-код для пополнения'}
                        style={{textAlign: "center", marginTop: 24}}>
                        <QRCode value={qrValue} size={200}/>
                    </Section>
                )}
            </List>
        </Page>
    );
};
