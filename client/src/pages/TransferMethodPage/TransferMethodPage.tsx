import { FC, useEffect, useState } from "react";
import {
    Avatar,
    Button,
    Cell,
    Input,
    List,
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
import {useNavigate} from "react-router-dom";

type TransferMethod = "email" | "username" | "card";
type Step = "choose" | "form";

const mapMethodToLabel = (method: TransferMethod) => {
    switch (method) {
        case "email":
            return "По email";
        case "username":
            return "По имени пользователя";
        case "card":
            return "По номеру счёта";
        default:
            return "";
    }
}

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

    const senderWallet = wallets.find(w => w.id === selectedWalletId);

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

        const fetchReceiver = async () => {
            setLoading(true);
            try {
                let user: User;
                if (method === "email") user = await UserService.getByEmail(target);
                else if (method === "username") user = await UserService.getByUsername(target);
                else user = await UserService.getByWalletNumber(target);

                setReceiver(user);

                const userWallets = await WalletService.getWalletsByUserId(user.id);
                setReceiverWallets(userWallets);

                if (senderWallet) {
                    const matched = userWallets.find(w => w.currency === senderWallet.currency);
                    setReceiverWalletId(matched?.id || null);
                }
            } catch {
                setReceiver(null);
                setReceiverWallets([]);
                setReceiverWalletId(null);
            } finally {
                setLoading(false);
            }
        };

        fetchReceiver();
    }, [target, senderWallet, step, method]);

    const handleTransfer = async () => {
        if (!selectedWalletId || !receiverWalletId || !amount) return;
        await TransferService.makeTransfer({
            senderWalletId: selectedWalletId,
            receiverWalletId,
            amount: parseInt(amount, 10),
        });
        navigate("/bank")
    };

    if (step === "choose") {
        return (
            <Page back={true}>
                <List>
                    <Title style={{ textAlign: "center", marginTop: 40, marginBottom: 20 }} level="1" weight="1">
                        Куда перевести?
                    </Title>

                    {METHOD_OPTIONS.map(({ key, label, subtitle, icon }) => (
                        <Cell
                            key={key}
                            before={icon}
                            subtitle={<Text color="secondary">{subtitle}</Text>}
                            onClick={() => {
                                setMethod(key);
                                setStep("form");
                            }}
                        >
                            {label}
                        </Cell>
                    ))}
                </List>
            </Page>
        );
    }

    return (
        <Page back={true}>
            <List>
                <Title style={{ textAlign: "center", marginTop: 40, marginBottom: 20 }} level="1" weight="1">
                    Перевод по {
                    method === "email" ? "email" :
                        method === "username" ? "имени пользователя" :
                            "номеру счёта"
                }
                </Title>

                <List>

                </List>

                <Cell subtitle={<Text color="secondary">Счёт для списания</Text>}>
                    <Select
                        value={selectedWalletId?.toString() || ""}
                        onChange={(e) => setSelectedWalletId(Number(e.target.value))}
                    >
                        <option disabled value="">Выберите счёт</option>
                        {wallets.map(w => (
                            <option key={w.id} value={w.id}>
                                {w.currency} ({w.number})
                            </option>
                        ))}
                    </Select>
                </Cell>

                <Input
                    placeholder={mapMethodToLabel(method)}
                    value={target}
                    onChange={(e) => setTarget(e.target.value)}
                />

                {loading && <Spinner size="l" />}
                {receiver && (
                    <Cell
                        before={<Avatar  fallbackIcon={<span>😕</span>} />}
                        subtitle={<Text color="secondary">{receiver.email || receiver.telegramUsername}</Text>}
                    >
                        {receiver.telegramFirstName || receiver.username}
                    </Cell>
                )}

                {receiver && receiverWallets.length > 0 && (
                    <Cell subtitle={<Text color="secondary">Счёт получателя</Text>}>
                        <Select
                            value={receiverWalletId?.toString() || ""}
                            onChange={(e) => setReceiverWalletId(Number(e.target.value))}
                        >
                            <option disabled value="">Выберите счёт</option>
                            {receiverWallets
                                .filter(w => !senderWallet || w.currency === senderWallet.currency)
                                .map(w => (
                                    <option key={w.id} value={w.id}>
                                        {w.currency} ({w.number})
                                    </option>
                                ))}
                        </Select>
                    </Cell>
                )}

                <Input
                    placeholder="Введите сумму"
                    type="number"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
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
                    onClick={() => setStep("choose")}
                >
                    Назад к выбору способа
                </Button>

            </List>
        </Page>
    );
};
