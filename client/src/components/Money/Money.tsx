import {FC, useEffect, useState} from "react";
import {
    Avatar,
    Cell,
    Image,
    Info,
    InlineButtons,
    LargeTitle,
    List,
    Placeholder,
    Section,
    Skeleton,
    Text,
    Title
} from "@telegram-apps/telegram-ui";
import {
    InlineButtonsItem
} from "@telegram-apps/telegram-ui/dist/components/Blocks/InlineButtons/components/InlineButtonsItem/InlineButtonsItem";
import {Icon16AddCircle, Icon32SendCircle} from "@vkontakte/icons";
import {Icon24QR} from "@telegram-apps/telegram-ui/dist/icons/24/qr";
import {initDataState as _initDataState, qrScanner, useSignal} from "@telegram-apps/sdk-react";
import WalletService, {Wallet} from "@/api/services/walletService.ts";
import RateService from "@/api/services/rateService.ts";
import {
    convertToRub,
    formatBalance,
    mapCurrencyToLogo,
    mapCurrencyToName,
    mapCurrencyToSymbol,
    sumAllWalletsInRub,
} from "@/utils/currency.ts";
import {useNavigate} from "react-router-dom";

export const Money: FC = () => {
    const initDataState = useSignal(_initDataState);
    const user = initDataState?.user;
    const [wallets, setWallets] = useState<Wallet[] | null>(null);
    const [rates, setRates] = useState<Record<string, number>>({});
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    // Получаем кошельки и курсы при загрузке компонента
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
            } catch (err) {
                setError('Ошибка загрузки данных');
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, []);

    const openQr = async () => {
        if (qrScanner.open.isAvailable()) {
            let promise = qrScanner.open({text: 'Scan any QR'});
            await promise.then(res => console.log(res));
        }
    };

    return (
        <List>
            <Cell
                before={
                    <Avatar
                        alt={'Telegram logo'}
                        src={user?.photo_url}
                    />
                }
            />
            <Placeholder
                header={'Баланс кошелька'}
                description={
                    <Skeleton visible={loading} withoutAnimation={false}>
                        <LargeTitle weight={'1'} style={{color: 'var(--tgui--text_color)'}}>
                            {loading ? '' : `${sumAllWalletsInRub(wallets, rates).toFixed(2)} ₽`}
                        </LargeTitle>
                    </Skeleton>
                }>
            </Placeholder>

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
                        <Cell
                            key={wallet.id}
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
                                    type={'text'}
                                    subtitle={`${convertToRub(wallet.balance, rates, wallet.currency)} ${mapCurrencyToSymbol('RUB')}`}
                                >
                                    {`${formatBalance(wallet.balance)} ${mapCurrencyToSymbol(wallet.currency)}`}
                                </Info>

                            }
                        >
                            <Title level="3">
                                {mapCurrencyToName(wallet.currency)}
                            </Title>
                        </Cell>
                    ))}
                </Section>
            )}

        </List>
    );
}
