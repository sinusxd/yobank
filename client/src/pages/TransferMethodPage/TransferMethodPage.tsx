import { FC, useEffect, useState } from "react";
import {
    Avatar,
    Button,
    Cell,
    Input,
    List,
    Section,
    Select,
    Spinner,
    Text,
    Title
} from "@telegram-apps/telegram-ui";
import {
    Icon28UserOutline,
    Icon28MailOutline,
    Icon28MoneyTransferOutline
} from "@vkontakte/icons";
import UserService, { User } from "@/api/services/userService";
import WalletService, { Wallet } from "@/api/services/walletService";
import TransferService from "@/api/services/transferService";
import { Page } from "@/components/Page";
import { useNavigate, useSearchParams } from "react-router-dom";
import { mapCurrencyToSymbol } from "@/utils/currency";

type TransferMethod = "email" | "username" | "card";
type Step = "choose" | "form";

const mapMethodToLabel = (method: TransferMethod) => {
    switch (method) {
        case "email": return "Email";
        case "username": return "Имя пользователя";
        case "card": return "Номер счёта";
    }
};

const METHOD_OPTIONS = [
    {
        key: "card",
        label: "По номеру счёта",
        subtitle: "Отправка по номеру кошелька",
        icon: <Icon28MoneyTransferOutline />
    },
    {
        key: "username",
        label: "По имени пользователя",
        subtitle: "Внутри приложения",
        icon: <Icon28UserOutline />
    },
    {
        key: "email",
        label: "По email",
        subtitle: "На зарегистрированный email",
        icon: <Icon28MailOutline />
    }
] as const;

export const TransferMethodPage: FC = () => {
    const [step, setStep] = useState<Step>("choose");
    const [method, setMethod] = useState<TransferMethod>("username");
    const [wallets, setWallets] = useState<Wallet[]>([]);
    const [selectedWalletId, setSelectedWalletId] = useState<number | null>(null);
    const [amount, setAmount] = useState<string>("");
    const [target, setTarget] = useState("");
    const [receiver, setReceiver] = useState<User | null>(null);
    const [receiverWallets, setReceiverWallets] = useState<Wallet[]>([]);
    const [receiverWalletId, setReceiverWalletId] = useState<number | null>(null);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();
    const [params] = useSearchParams();

    const senderWallet = wallets.find(w => w.id === selectedWalletId);

    useEffect(() => {
        const m = params.get("method") as TransferMethod | null;
        const t = params.get("target");
        const a = params.get("amount");
        if (m && t) {
            setMethod(m);
            setTarget(t);
            if (a) setAmount(a);
            setStep("form");
        }
    }, []);

    useEffect(() => {
        WalletService.getUserWallets().then(res => setWallets(res.data));
    }, []);

    useEffect(() => {
        if (step !== "form" || target.length < 3) {
            setReceiver(null);
            setReceiverWallets([]);
            setReceiverWalletId(null);
            return;
        }

        setLoading(true);

        let userPromise: Promise<User>;
        switch (method) {
            case "email":
                userPromise = UserService.getByEmail(target);
                break;
            case "username":
                userPromise = UserService.getByUsername(target);
                break;
            case "card":
                userPromise = UserService.getByWalletNumber(target);
                break;
            default:
                userPromise = Promise.reject("Неверный метод");
        }

        userPromise
            .then(user => {
                setReceiver(user);
                return WalletService.getWalletsByUserId(user.id);
            })
            .then(wds => {
                setReceiverWallets(wds);
                if (method === "card") {
                    const byNum = wds.find(w => w.number === target);
                    setReceiverWalletId(byNum ? byNum.id : null);
                }
            })
            .catch(() => {
                setReceiver(null);
                setReceiverWallets([]);
                setReceiverWalletId(null);
            })
            .finally(() => setLoading(false));
    }, [step, target, method]);

    useEffect(() => {
        if (!receiverWalletId || wallets.length === 0) return;
        const recv = receiverWallets.find(w => w.id === receiverWalletId);
        if (!recv) return;
        const match = wallets.find(w => w.currency === recv.currency);
        if (match) setSelectedWalletId(match.id);
    }, [receiverWalletId, wallets, receiverWallets]);

    const handleTransfer = async () => {
        if (!selectedWalletId || !receiverWalletId || !amount) return;
        const num = parseFloat(amount.replace(",", "."));
        if (isNaN(num) || num <= 0) {
            alert("Введите корректную сумму");
            return;
        }
        try {
            await TransferService.makeTransfer({
                senderWalletId: selectedWalletId,
                receiverWalletId,
                amount: Math.round(num * 100),
            });
            alert("✅ Перевод выполнен");
            navigate("/bank");
        } catch (err: any) {
            alert(`Ошибка: ${err?.response?.data?.message || err.message}`);
        }
    };

    if (step === "choose") {
        return (
            <Page back>
                <List>
                    <Title level="1" weight="1" style={{ textAlign: "center", margin: "40px 0 20px" }}>
                        Куда перевести?
                    </Title>
                    {METHOD_OPTIONS.map(opt => (
                        <Cell
                            key={opt.key}
                            before={opt.icon}
                            subtitle={<Text color="secondary">{opt.subtitle}</Text>}
                            onClick={() => {
                                setMethod(opt.key);
                                setStep("form");
                            }}
                        >
                            {opt.label}
                        </Cell>
                    ))}
                </List>
            </Page>
        );
    }

    return (
        <Page back>
            <List>
                <Title level="1" weight="1" style={{ textAlign: "center", margin: "40px 0 20px" }}>
                    Перевод {mapMethodToLabel(method).toLowerCase()}
                </Title>

                <Section header="Счёт для списания">
                    <Select
                        value={selectedWalletId?.toString() || ""}
                        onChange={e => setSelectedWalletId(Number(e.target.value))}
                    >
                        <option disabled value="">Выберите счёт</option>
                        {wallets.map(w =>
                            <option key={w.id} value={w.id}>
                                {w.currency} ({w.number}) • {(w.balance/100).toFixed(2)} {mapCurrencyToSymbol(w.currency)}
                            </option>
                        )}
                    </Select>
                </Section>

                <Input
                    placeholder={mapMethodToLabel(method)}
                    value={target}
                    onChange={e => setTarget(e.target.value)}
                    style={{ margin: "12px 0" }}
                />

                {loading && <Spinner size="l" />}

                {receiver && (
                    <Cell
                        before={<Avatar src={receiver.avatarUrl || `https://avatars.githubusercontent.com/u/${receiver.id % 1000000}?v=4`} />}
                        subtitle={<Text color="secondary">{receiver.email || receiver.telegramUsername}</Text>}
                    >
                        {receiver.telegramFirstName || receiver.username}
                    </Cell>
                )}

                {receiverWallets.length > 0 && (
                    <Section header="Счёт получателя">
                        <Select
                            value={receiverWalletId?.toString() || ""}
                            onChange={e => setReceiverWalletId(Number(e.target.value))}
                        >
                            <option disabled value="">Выберите счёт</option>
                            {receiverWallets.map(w =>
                                <option key={w.id} value={w.id}>
                                    {w.currency} ({w.number})
                                </option>
                            )}
                        </Select>
                    </Section>
                )}

                <Input
                    placeholder={`Сумма в ${mapCurrencyToSymbol(senderWallet?.currency || "")}`}
                    value={amount ? `${amount} ${mapCurrencyToSymbol(senderWallet?.currency || "")}` : ""}
                    onChange={e => {
                        const raw = e.target.value.replace(/[^\d.,]/g, "");
                        const norm = raw.replace(",", ".");
                        if (/^\d*([.]?\d{0,2})?$/.test(norm)) {
                            setAmount(norm);
                            requestAnimationFrame(() => {
                                const el = e.target as HTMLInputElement;
                                const pos = norm.length;
                                el.setSelectionRange(pos, pos);
                            });
                        }
                    }}
                    style={{ margin: "12px 0" }}
                />

                <Button
                    size="l"
                    stretched
                    onClick={handleTransfer}
                    disabled={!selectedWalletId || !receiverWalletId || !amount}
                >
                    Перевести
                </Button>

                <Button
                    size="l"
                    stretched
                    style={{ marginTop: 12 }}
                    onClick={() => setStep("choose")}
                >
                    Назад
                </Button>
            </List>
        </Page>
    );
};
