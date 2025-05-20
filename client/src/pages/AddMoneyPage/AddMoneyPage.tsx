import {FC} from "react";
import {useNavigate} from "react-router-dom";
import {Page} from "@/components/Page";
import {Cell, List, Text, Title} from "@telegram-apps/telegram-ui";
import {bem} from "@/css/bem.ts";
import './AddMoneyPage.css'
import {Icon28PaymentCardAddOutline, Icon28QrCodeOutline} from "@vkontakte/icons";

const ADD_MONEY_METHODS = [
    {
        key: "card",
        icon: <Icon28PaymentCardAddOutline />,
        title: "Онлайн пополнение",
        subtitle: "Купить валюту онлайн",
        route: "/add-money/online"
    },
    {
        key: "qr",
        icon: <Icon28QrCodeOutline />,
        title: "Пополнение QR-кодом",
        subtitle: "QR-код для перевода денег вам",
        route: "/add-money/qr"
    }
];

const [, e] = bem('add-money-page');

export const AddMoneyPage: FC = () => {
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
                    Выберите способ пополнения
                </Title>
                {ADD_MONEY_METHODS.map(method => (
                    <Cell
                        key={method.key}
                        before={method.icon}
                        subtitle={
                            <Text style={{color: "#8f8f8f"}}>{method.subtitle}</Text>
                        }
                        onClick={() => navigate(method.route)}
                        className={e('cell')}
                    >
                        <span style={{fontWeight: 500}}>{method.title}</span>
                    </Cell>
                ))}
            </List>
        </Page>
    );
};
