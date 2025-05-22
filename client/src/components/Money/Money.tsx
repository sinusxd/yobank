import {FC, useEffect, useState} from "react";
import {
    Avatar,
    Cell,
    Image,
    Info,
    InlineButtons,
    LargeTitle,
    List,
    Modal,
    Placeholder,
    Section,
    Skeleton,
    Text,
    Title,
} from "@telegram-apps/telegram-ui";
import {InlineButtonsItem} from "@telegram-apps/telegram-ui/dist/components/Blocks/InlineButtons/components/InlineButtonsItem/InlineButtonsItem";
import {Icon16AddCircle, Icon32SendCircle} from "@vkontakte/icons";
import {Icon24QR} from "@telegram-apps/telegram-ui/dist/icons/24/qr";
import {Icon28Close} from "@telegram-apps/telegram-ui/dist/icons/28/close";
import {initDataState as _initDataState, qrScanner} from "@telegram-apps/sdk-react";
import WalletService, {Wallet} from "@/api/services/walletService.ts";
import RateService from "@/api/services/rateService.ts";
import moneyGif from "./money.gif";
import {
    convertToRub,
    formatBalance,
    mapCurrencyToLogo,
    mapCurrencyToName,
    mapCurrencyToSymbol,
    sumAllWalletsInRub,
} from "@/utils/currency.ts";
import {useNavigate} from "react-router-dom";
import {ModalHeader} from "@telegram-apps/telegram-ui/dist/components/Overlays/Modal/components/ModalHeader/ModalHeader";
import {ModalClose} from "@telegram-apps/telegram-ui/dist/components/Overlays/Modal/components/ModalClose/ModalClose";
import {useSelector} from "react-redux";
import {RootState} from "@/store/store.ts";

export const Money: FC = () => {
    const userFromRedux = useSelector((state: RootState) => state.user);
    const [wallets, setWallets] = useState<Wallet[] | null>(null);
    const [rates, setRates] = useState<Record<string, number>>({});
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                const [{data: walletsData}, {data: usdRate}, {data: eurRate}, {data: cnyRate}] = await Promise.all([
                    WalletService.getUserWallets(),
                    RateService.getRate('USD'),
                    RateService.getRate('EUR'),
                    RateService.getRate('CNY')
                ]);
                setWallets(walletsData);
                setRates({
                    'RUB': 1,
                    'USD': usdRate.value,
                    'EUR': eurRate.value,
                    'CNY': cnyRate.value,
                });
                setError(null);
            } catch {
                setError('Ошибка загрузки данных');
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, []);

    const openQr = async () => {
        if (qrScanner.open.isAvailable()) {
            const res = await qrScanner.open({ text: 'Сканируйте QR-код для перевода' });
            const data = res;

            if (data && typeof data === "string") {
                try {
                    const url = new URL(data);

                    if (url.protocol === "yobank:") {
                        const method = url.searchParams.get("method");
                        const target = url.searchParams.get("target");
                        const amount = url.searchParams.get("amount");

                        const query = new URLSearchParams();

                        if (method) query.set("method", method);
                        if (target) query.set("target", target);
                        if (amount) query.set("amount", amount);

                        navigate(`/transfer-money?${query.toString()}`);
                    } else {
                        alert("Неверный формат QR-кода");
                    }
                } catch {
                    alert("Не удалось обработать QR-код");
                }
            } else {
                alert("QR-код не содержит данных");
            }
        }
    };


    return (
        <List>
            <Cell
                before={
                    <Avatar
                        alt="Telegram logo"
                        src={userFromRedux?.avatarUrl || `https://avatars.githubusercontent.com/u/${10 % 1000000}?v=4`}
                    />
                }
            />
            <Placeholder
                header="Баланс кошелька"
                description={
                    <Skeleton visible={loading}>
                        <LargeTitle weight="1" style={{color: 'var(--tgui--text_color)'}}>
                            {loading ? '' : `${sumAllWalletsInRub(wallets, rates).toFixed(2)} ₽`}
                        </LargeTitle>
                    </Skeleton>
                }
            />

            <InlineButtons mode="bezeled">
                <InlineButtonsItem text="Отправить" onClick={() => navigate("/transfer-money")}>
                    <Icon32SendCircle width={24} height={24}/>
                </InlineButtonsItem>
                <InlineButtonsItem text="Пополнить" onClick={() => navigate("/add-money")}>
                    <Icon16AddCircle width={24} height={24}/>
                </InlineButtonsItem>
                <InlineButtonsItem text="QR" onClick={openQr}>
                    <Icon24QR width={24} height={24}/>
                </InlineButtonsItem>
            </InlineButtons>

            {error && (
                <Section>
                    <Cell style={{color: 'var(--tgui--destructive_text_color)'}}>
                        {error}
                    </Cell>
                </Section>
            )}

            {wallets && wallets.length > 0 && (
                <Section header="Ваши счета">
                    {wallets.map(wallet => (
                        <Modal
                            key={wallet.id}
                            header={
                                <ModalHeader after={<ModalClose><Icon28Close style={{color: 'var(--tgui--plain_foreground)'}} /></ModalClose>}>
                                    {mapCurrencyToName(wallet.currency)}
                                </ModalHeader>
                            }
                            trigger={
                                <Cell
                                    subtitle={
                                        <Text>
                                            {rates[wallet.currency].toFixed(2)} {mapCurrencyToSymbol(wallet.currency)}
                                        </Text>
                                    }
                                    before={
                                        <Image
                                            src={mapCurrencyToLogo(wallet.currency)}
                                            style={{boxShadow: "none", backgroundColor: "transparent"}}
                                        />
                                    }
                                    after={
                                        <Info
                                            type="text"
                                            subtitle={`${convertToRub(wallet.balance, rates, wallet.currency)} ₽`}
                                        >
                                            {formatBalance(wallet.balance)} {mapCurrencyToSymbol(wallet.currency)}
                                        </Info>
                                    }
                                >
                                    <Title level="3">{mapCurrencyToName(wallet.currency)}</Title>
                                </Cell>
                            }
                        >
                            <Placeholder
                                header={`Баланс: ${formatBalance(wallet.balance)} ${mapCurrencyToSymbol(wallet.currency)}`}
                                description={
                                    <>
                                        {`В рублях: ${convertToRub(wallet.balance, rates, wallet.currency)} ₽`}
                                        <br />
                                        {`Номер карты: ${wallet.number}`}
                                    </>
                                }
                            >
                                <img
                                    alt="Telegram sticker"
                                    src={moneyGif}
                                    style={{
                                        display: 'block',
                                        height: '144px',
                                        width: '144px'
                                    }}
                                />
                            </Placeholder>
                        </Modal>
                    ))}
                </Section>
            )}
        </List>
    );
};
