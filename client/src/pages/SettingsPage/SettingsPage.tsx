import { FC, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
    Button,
    Cell,
    List,
    Section,
    Spinner,
    Switch,
    Text,
    Title,
} from "@telegram-apps/telegram-ui";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "@/store/store.ts";
import { Link } from "@/components/Link/Link.tsx";
import UserService from "@/api/services/userService.ts";
import {setUser} from "@/store/userSlice.ts";

export const SettingsPage: FC = () => {
    const navigate = useNavigate();
    const dispatch = useDispatch();

    const user = useSelector((state: RootState) => state.user);
    const loading = user === null;

    const [notificationsEnabled, setNotificationsEnabled] = useState<boolean>(false);
    const [updating, setUpdating] = useState<boolean>(false);

    useEffect(() => {
        if (user) {
            setNotificationsEnabled(user.notification ?? false);
        }
    }, [user]);

    const handleToggleNotifications = async () => {
        if (!user) return;

        const newValue = !notificationsEnabled;
        setNotificationsEnabled(newValue);
        setUpdating(true);

        try {
            await UserService.toggleNotification(user.id, newValue);
            dispatch(setUser({ ...user, notification: newValue }));
        } catch (error) {
            alert("Не удалось обновить уведомления");
            setNotificationsEnabled(!newValue); // откат
        } finally {
            setUpdating(false);
        }
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
                                    disabled={updating}
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
