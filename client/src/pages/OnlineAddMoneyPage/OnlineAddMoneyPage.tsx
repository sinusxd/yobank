import { FC, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Page } from "@/components/Page.tsx";
import { Button, Cell, Input, List, Section, Title } from "@telegram-apps/telegram-ui";
import { Image } from "@telegram-apps/telegram-ui";
import { mapCurrencyToLogo, mapCurrencyToName, mapCurrencyToSymbol } from "@/utils/currency";
import WalletService from "@/api/services/walletService.ts";

const wallets = [
    { id: "0", currency: "RUB" },
    { id: "1", currency: "USD" },
    { id: "2", currency: "EUR" },
    { id: "3", currency: "CNY" },
];

export const OnlineAddMoneyPage: FC = () => {
    const [selectedCurrency, setSelectedCurrency] = useState<string | null>(null);
    const [amount, setAmount] = useState<string>("");
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    return (
        <Page back>
            {!selectedCurrency ? (
                <>
                    <Title
                        style={{
                            textAlign: "center",
                            marginTop: "100px",
                            marginBottom: "50px",
                        }}
                        level="1"
                        weight="1"
                    >
                        Выберите валюту для пополнения
                    </Title>
                    <Section style={{ margin: "15px" }}>
                        {wallets.map((wallet) => (
                            <Cell
                                key={wallet.id}
                                before={
                                    <Image
                                        src={mapCurrencyToLogo(wallet.currency)}
                                        style={{
                                            boxShadow: "none",
                                            backgroundColor: "transparent",
                                        }}
                                    />
                                }
                                onClick={() => setSelectedCurrency(wallet.currency)}
                            >
                                <Title level="3">
                                    {mapCurrencyToName(wallet.currency)}
                                </Title>
                            </Cell>
                        ))}
                    </Section>
                </>
            ) : (
                <>
                    <Title
                        style={{
                            textAlign: "center",
                            marginTop: "80px",
                            marginBottom: "30px",
                        }}
                        level="2"
                        weight="1"
                    >
                        Введите сумму в {mapCurrencyToSymbol(selectedCurrency)}
                    </Title>
                    <List>
                        <Input
                            placeholder={`Сумма в ${mapCurrencyToSymbol(selectedCurrency)}`}
                            value={
                                amount
                                    ? `${amount} ${mapCurrencyToSymbol(selectedCurrency)}`
                                    : ""
                            }
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
                            type="text"
                            inputMode="decimal"
                        />

                        <Button
                            style={{ marginTop: "20px" }}
                            size="l"
                            stretched
                            loading={loading}
                            disabled={!amount || loading}
                            onClick={async () => {
                                try {
                                    const numericAmount = parseFloat(
                                        amount.replace(",", ".")
                                    );
                                    if (isNaN(numericAmount) || numericAmount <= 0) {
                                        alert("Введите корректную сумму");
                                        return;
                                    }

                                    const amountInMinorUnits = Math.round(
                                        numericAmount * 100
                                    );

                                    setLoading(true);
                                    await WalletService.topUpWallet(
                                        selectedCurrency!,
                                        amountInMinorUnits
                                    );

                                    alert(
                                        `Кошелек пополнен на ${numericAmount.toFixed(
                                            2
                                        )} ${mapCurrencyToSymbol(selectedCurrency!)}`
                                    );

                                    navigate("/bank");
                                } catch (error) {
                                    console.error(error);
                                    alert("Не удалось пополнить кошелек");
                                } finally {
                                    setLoading(false);
                                }
                            }}
                        >
                            Пополнить {amount} {selectedCurrency}
                        </Button>
                    </List>
                </>
            )}
        </Page>
    );
};

export default OnlineAddMoneyPage;
