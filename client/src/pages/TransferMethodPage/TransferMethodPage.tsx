import { FC } from "react";
import { useNavigate } from "react-router-dom";
import { Page } from "@/components/Page";
import { Cell, List, Text, Title } from "@telegram-apps/telegram-ui";
import { bem } from "@/css/bem.ts";
import {
    Icon28UserOutline,
    Icon28MailOutline,
    Icon28MoneyTransferOutline
} from "@vkontakte/icons";
import './TransferMethodPage.css';


const TRANSFER_METHODS = [
    {
        key: "card",
        icon: <Icon28MoneyTransferOutline />,
        title: "По номеру карты",
        subtitle: "Отправьте деньги на банковскую карту",
        route: "/transfer/card"
    },
    {
        key: "username",
        icon: <Icon28UserOutline />,
        title: "По имени пользователя",
        subtitle: "Перевод внутри приложения",
        route: "/transfer/username"
    },
    {
        key: "email",
        icon: <Icon28MailOutline />,
        title: "По email",
        subtitle: "Отправка по электронной почте",
        route: "/transfer/email"
    }
];

const [, e] = bem('transfer-method-page');

export const TransferMethodPage: FC = () => {
    const navigate = useNavigate();

    return (
        <Page back>
            <List>
                <Title
                    style={{
                        textAlign: "center",
                        marginTop: "100px",
                        marginBottom: "50px"
                    }}
                    level="1"
                    weight="1"
                >
                    Куда перевести деньги?
                </Title>
                {TRANSFER_METHODS.map(method => (
                    <Cell
                        key={method.key}
                        before={method.icon}
                        subtitle={
                            <Text style={{ color: "#8f8f8f" }}>{method.subtitle}</Text>
                        }
                        onClick={() => navigate(method.route)}
                        className={e('cell')}
                    >
                        <span style={{ fontWeight: 500 }}>{method.title}</span>
                    </Cell>
                ))}
            </List>
        </Page>
    );
};
