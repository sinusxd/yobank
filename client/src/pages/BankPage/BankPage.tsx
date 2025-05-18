import { FC, useState } from "react";
import { Page } from "@/components/Page.tsx";
import { Tabbar } from "@telegram-apps/telegram-ui";
import {
    Icon28MoneyHistoryBackwardOutline,
    Icon28PaymentCardAddOutline,
    Icon28SettingsOutline,
    Icon28WalletOutline
} from "@vkontakte/icons";
import {Money} from "@/components/Money/Money.tsx";

export const BankPage: FC = () => {
    const [activeTab, setActiveTab] = useState("wallet");

    return (
        <Page back={false}>
            {/* Здесь можно рендерить разный контент по activeTab */}
            <Money/>
            <Tabbar style={{ paddingBottom: "25px" }}>
                {[
                    <Tabbar.Item
                        key="wallet"
                        text="Кошелёк"
                        selected={activeTab === "wallet"}
                        onClick={() => setActiveTab("wallet")}
                    >
                        <Icon28WalletOutline />
                    </Tabbar.Item>,
                    <Tabbar.Item
                        key="history"
                        text="История"
                        selected={activeTab === "history"}
                        onClick={() => setActiveTab("history")}
                    >
                        <Icon28MoneyHistoryBackwardOutline />
                    </Tabbar.Item>,
                    <Tabbar.Item
                        key="payments"
                        text="Платежи"
                        selected={activeTab === "payments"}
                        onClick={() => setActiveTab("payments")}
                    >
                        <Icon28PaymentCardAddOutline />
                    </Tabbar.Item>,
                    <Tabbar.Item
                        key="settings"
                        text="Настройки"
                        selected={activeTab === "settings"}
                        onClick={() => setActiveTab("settings")}
                    >
                        <Icon28SettingsOutline />
                    </Tabbar.Item>,
                ]}
            </Tabbar>
        </Page>
    );
};
