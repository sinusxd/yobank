import {FC} from "react";
import {Avatar, Cell, InlineButtons, LargeTitle, List, Placeholder, Section} from "@telegram-apps/telegram-ui";
import {
    InlineButtonsItem
} from "@telegram-apps/telegram-ui/dist/components/Blocks/InlineButtons/components/InlineButtonsItem/InlineButtonsItem";
import {Icon16AddCircle, Icon32SendCircle} from "@vkontakte/icons";
import {Icon24QR} from "@telegram-apps/telegram-ui/dist/icons/24/qr";
import {useSignal,
    initDataRaw as _initDataRaw,
    initDataState as _initDataState,
    qrScanner
} from "@telegram-apps/sdk-react";

export const Money: FC = () => {
    const initDataState = useSignal(_initDataState);
    const user = initDataState?.user;

    const openQr = async () => {
        if (qrScanner.open.isAvailable()) {
            qrScanner.isOpened(); // false
            let promise = qrScanner.open({ text: 'Scan any QR' });
            qrScanner.isOpened(); // true
            await promise.then(res => console.log(res));
            qrScanner.isOpened(); // false
        }
    }

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
                    <LargeTitle
                        weight={'1'}
                        style={{
                            color: 'var(--tgui--text_color)'
                        }}
                    >
                        0.00 ₽
                    </LargeTitle>
                }>
            </Placeholder>

            <InlineButtons mode="bezeled">
                <InlineButtonsItem text="Отправить">
                    <Icon32SendCircle width={24} height={24} />
                </InlineButtonsItem>
                <InlineButtonsItem text="Пополнить">
                    <Icon16AddCircle width={24} height={24} />
                </InlineButtonsItem>
                <InlineButtonsItem text="QR" onClick={openQr}>
                    <Icon24QR width={24} height={24} />
                </InlineButtonsItem>

            </InlineButtons>
            <Section>
                <Cell>
                    aboba
                </Cell>
            </Section>

        </List>

    )
}