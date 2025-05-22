import {FC, useState} from "react";
import {Page} from "@/components/Page.tsx";
import {Tabbar} from "@telegram-apps/telegram-ui";
import {Icon28MoneyHistoryBackwardOutline, Icon28SettingsOutline, Icon28WalletOutline} from "@vkontakte/icons";
import {Money} from "@/components/Money/Money.tsx";
import SettingsPage from "@/pages/SettingsPage/SettingsPage.tsx";
import {TransactionHistory} from "@/components/TransactionHistory/TransactionHistory.tsx";

export const BankPage: FC = () => {
    const [activeTab, setActiveTab] = useState("wallet");

    const renderContent = () => {
        switch (activeTab) {
            case "wallet":
                return <Money/>;
            case "settings":
                return <SettingsPage/>;
            case "history":
                return <TransactionHistory/>;
            case "payments":
                return <div style={{padding: "16px"}}>Платежи (в разработке)</div>;
            default:
                return null;
        }
    };

    return (
        <Page back>
            <div style={{ flex: 1, overflowY: 'auto', paddingBottom: 80 }}>
                {renderContent()}
            </div>

            <Tabbar style={{paddingBottom: "20px"}}>
                {[
                    <Tabbar.Item
                        key="wallet"
                        text="Кошелёк"
                        selected={activeTab === "wallet"}
                        onClick={() => setActiveTab("wallet")}
                    >
                        <Icon28WalletOutline/>
                    </Tabbar.Item>,
                    <Tabbar.Item
                        key="history"
                        text="История"
                        selected={activeTab === "history"}
                        onClick={() => setActiveTab("history")}
                    >
                        <Icon28MoneyHistoryBackwardOutline/>
                    </Tabbar.Item>,
                    <Tabbar.Item
                        key="settings"
                        text="Настройки"
                        selected={activeTab === "settings"}
                        onClick={() => setActiveTab("settings")}
                    >
                        <Icon28SettingsOutline/>
                    </Tabbar.Item>,
                ]}
            </Tabbar>
        </Page>
    );
};
