import {FC, useState} from "react";
import {Page} from "@/components/Page.tsx";
import {Button, Cell, Input, List, Section, Title} from "@telegram-apps/telegram-ui";
import {Image} from "@telegram-apps/telegram-ui";
import {mapCurrencyToLogo, mapCurrencyToName, mapCurrencyToSymbol} from "@/utils/currency";

const wallets = [
    {id: '0', currency: 'RUB'},
    {id: '1', currency: 'USD'},
    {id: '2', currency: 'EUR'},
    {id: '3', currency: 'CNY'}
];

export const OnlineAddMoneyPage: FC = () => {
    const [selectedCurrency, setSelectedCurrency] = useState<string | null>(null);
    const [amount, setAmount] = useState<string>("");


    return (
        <Page back>
            {!selectedCurrency ? (
                <>
                    <Title
                        style={{
                            textAlign: "center",
                            marginTop: "100px",
                            marginBottom: "50px"
                        }}
                        level="1"
                        weight="1"
                    >
                        Выберите валюту для пополнения
                    </Title>
                    <Section
                        style={{margin:"15px"}}
                    >
                        {wallets.map(wallet => (
                            <Cell
                                key={wallet.id}
                                before={
                                    <Image
                                        src={mapCurrencyToLogo(wallet.currency)}
                                        style={{boxShadow: "none", backgroundColor: "transparent"}}
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
                            marginBottom: "30px"
                        }}
                        level="2"
                        weight="1"
                    >
                        Введите сумму в {mapCurrencyToSymbol(selectedCurrency)}
                    </Title>
                    <List>
                        <Input
                            placeholder={`Сумма в ${mapCurrencyToSymbol(selectedCurrency)}`}
                            value={amount ? `${amount} ${mapCurrencyToSymbol(selectedCurrency)}` : ""}
                            onChange={(e) => {
                                const raw = e.target.value.replace(/[^\d.,]/g, '');
                                const normalized = raw.replace(',', '.');
                                if (/^\d*([.]?\d{0,2})?$/.test(normalized)) {
                                    setAmount(normalized);

                                    // принудительно вернуть курсор до символа
                                    requestAnimationFrame(() => {
                                        const el = e.target as HTMLInputElement;
                                        const pos = normalized.length;
                                        el.setSelectionRange(pos, pos); // курсор перед символом
                                    });
                                }
                            }}

                            type="text"
                            inputMode="decimal"
                        />

                        <Button
                            style={{marginTop: "20px"}}
                            size="l"
                            stretched
                            disabled={!amount}
                            onClick={() => {
                                alert(`Пополнение на ${amount} ${mapCurrencyToSymbol(selectedCurrency)}`);
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
