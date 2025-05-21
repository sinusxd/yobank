import {FC, useState} from "react";
import { useNavigate } from "react-router-dom";
import {
    Button,
    Cell,
    List,
    Section,
    Spinner,
    Switch,
    Text,
    Title
} from "@telegram-apps/telegram-ui";
import { useSelector } from "react-redux";
import {RootState} from "@/store/store.ts";
import {Link} from "@/components/Link/Link.tsx";

export const SettingsPage: FC = () => {
    const navigate = useNavigate();

    const [notificationsEnabled, setNotificationsEnabled] = useState<boolean>(
        localStorage.getItem("notifications-enabled") === "true"
    );

    const user = useSelector((state: RootState) => state.user);
    const loading = user === null;

    const handleToggleNotifications = () => {
        const newValue = !notificationsEnabled;
        setNotificationsEnabled(newValue);
        sessionStorage.setItem("notifications-enabled", String(newValue));
    };

    const handleLogout = () => {
        sessionStorage.removeItem("access_token");
        sessionStorage.removeItem("refresh_token");
        navigate("/");
    };

    return (
        <List>
            <Title level="1" weight="2" style={{ textAlign: "center", marginTop: "24px", marginBottom: "16px" }}>
                Настройки
            </Title>

            {loading ? (
                <Spinner size="l" />
            ) : (
                <>
                    <Section header="Профиль">
                        <Cell>
                            <Text weight="2">{user?.username}</Text>
                        </Cell>
                    </Section>

                    <Section header="Оповещения">
                        <Cell
                            after={
                                <Switch
                                    checked={notificationsEnabled}
                                    onChange={handleToggleNotifications}
                                />
                            }
                        >
                            Уведомления
                        </Cell>
                    </Section>

                    <Section header="О приложении">
                        <Cell>
                            <Text>Версия: 1.0.0</Text>
                        </Cell>
                        <Cell>
                            <Text>
                                Разработчик: <Link to={'https://t.me/roflandown'}>@roflandown</Link>
                            </Text>
                        </Cell>
                    </Section>

                    <Section>
                        <Button
                            style={{ marginTop: "8px" }}
                            size="l"
                            stretched
                            onClick={handleLogout}
                        >
                            Выйти
                        </Button>
                    </Section>
                </>
            )}
        </List>
    );
};

export default SettingsPage;
